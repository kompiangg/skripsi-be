package order

import (
	"context"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/model"
	"skripsi-be/type/params"
	"skripsi-be/type/result"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/volatiletech/null/v9"
)

func (s service) FindOrder(ctx context.Context, param params.FindOrderService) (allOrders []result.Order, err error) {
	err = param.Validate()
	if err != nil {
		return nil, errors.Wrap(err)
	}

	if param.StartDate.After(s.config.Date.Now()) || param.EndDate.After(s.config.Date.Now()) {
		return nil, errors.Wrap(errors.ErrDataParamMustNotAFterCurrentTime)
	}

	if param.StartDate.After(param.EndDate) {
		return nil, errors.Wrap(errors.ErrDataParamStartDateMustNotAfterEndDate)
	}

	param.StartDate = time.Date(param.StartDate.Year(), param.StartDate.Month(), param.StartDate.Day(), 0, 0, 0, 0, param.StartDate.Location())
	param.EndDate = time.Date(param.EndDate.Year(), param.EndDate.Month(), param.EndDate.Day(), 23, 59, 59, 0, param.EndDate.Location())

	shardQuery, err := s.getShardWhereQuery(param.StartDate, param.EndDate)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	if s.config.IsUsingSharding {
		errorChan := make(chan error, len(shardQuery))
		defer close(errorChan)

		shardOrder := make([][]result.Order, len(shardQuery))
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		for idx, shardParam := range shardQuery {
			go func(idx int, shardParam result.ShardTimeSeriesWhereQuery) {
				var order []model.OrderWithOrderDetails

				if shardParam.ShardIndex == len(s.config.Shards) {
					order, err = s.orderRepo.FindAllOrderAndDetailsOnLongTermDB(ctx, params.LongTermWhereQuery{
						StartDate: null.TimeFrom(shardParam.StartDate),
						EndDate:   null.TimeFrom(shardParam.EndDate),
					})
					if err != nil {
						errorChan <- err
						return
					}
				} else {
					order, err = s.orderRepo.FindAllOrderAndDetailsOnShardDB(ctx, params.ShardTimeSeriesWhereQuery{
						ShardIndex: shardParam.ShardIndex,
						StartDate:  null.TimeFrom(shardParam.StartDate),
						EndDate:    null.TimeFrom(shardParam.EndDate),
					})
					if err != nil {
						errorChan <- err
						return
					}
				}

				orders := make([]result.Order, 0)
				orderIDMap := make(map[ulid.ULID]int)

				for idx := range order {
					orderID, err := ulid.Parse(order[idx].OrderID)
					if err != nil {
						errorChan <- err
						return
					}

					if _, ok := orderIDMap[orderID]; !ok {
						orders = append(orders, result.Order{
							ID:              orderID,
							CashierID:       order[idx].CashierID,
							StoreID:         order[idx].StoreID,
							PaymentID:       order[idx].PaymentID,
							CustomerID:      order[idx].CustomerID,
							TotalQuantity:   order[idx].TotalQuantity,
							TotalUnit:       order[idx].TotalUnit,
							TotalPrice:      order[idx].TotalPrice,
							TotalPriceInUSD: order[idx].TotalPriceInUSD,
							Currency:        order[idx].Currency,
							UsdRate:         order[idx].UsdRate,
							CreatedAt:       order[idx].CreatedAt.In(param.StartDate.Location()),
							OrderDetails:    make([]result.OrderDetail, 0),
						})

						orders[len(orders)-1].OrderDetails = append(orders[len(orders)-1].OrderDetails, result.OrderDetail{
							ID:       order[idx].OrderDetailID,
							OrderID:  orderID,
							ItemID:   order[idx].ItemID,
							Quantity: order[idx].Quantity,
							Unit:     order[idx].Unit,
							Price:    order[idx].Price,
						})

						orderIDMap[orderID] = len(orders) - 1
					} else {
						orders[orderIDMap[orderID]].OrderDetails = append(orders[orderIDMap[orderID]].OrderDetails, result.OrderDetail{
							ID:       order[idx].OrderDetailID,
							OrderID:  orderID,
							ItemID:   order[idx].ItemID,
							Quantity: order[idx].Quantity,
							Unit:     order[idx].Unit,
							Price:    order[idx].Price,
						})
					}
				}

				shardOrder[idx] = orders
				errorChan <- nil
			}(idx, shardParam)
		}

		notErrCount := 0
		for notErrCount != len(shardQuery) {
			err := <-errorChan
			if err != nil {
				cancel()
				return nil, errors.Wrap(err)
			}

			notErrCount++
		}

		for _, orders := range shardOrder {
			allOrders = append(allOrders, orders...)
		}
	} else {
		order, err := s.orderRepo.FindAllOrderAndDetailsOnLongTermDB(ctx, params.LongTermWhereQuery{
			StartDate: null.TimeFrom(param.StartDate),
			EndDate:   null.TimeFrom(param.EndDate),
		})
		if err != nil {
			return nil, errors.Wrap(err)
		}

		orders := make([]result.Order, 0)
		orderIDMap := make(map[ulid.ULID]int)

		for idx := range order {
			orderID, err := ulid.Parse(order[idx].OrderID)
			if err != nil {
				return nil, errors.Wrap(err)
			}

			if _, ok := orderIDMap[orderID]; !ok {
				orders = append(orders, result.Order{
					ID:              orderID,
					CashierID:       order[idx].CashierID,
					StoreID:         order[idx].StoreID,
					PaymentID:       order[idx].PaymentID,
					CustomerID:      order[idx].CustomerID,
					TotalQuantity:   order[idx].TotalQuantity,
					TotalUnit:       order[idx].TotalUnit,
					TotalPrice:      order[idx].TotalPrice,
					TotalPriceInUSD: order[idx].TotalPriceInUSD,
					Currency:        order[idx].Currency,
					UsdRate:         order[idx].UsdRate,
					CreatedAt:       order[idx].CreatedAt.In(param.StartDate.Location()),
					OrderDetails:    make([]result.OrderDetail, 0),
				})

				orders[len(orders)-1].OrderDetails = append(orders[len(orders)-1].OrderDetails, result.OrderDetail{
					ID:       order[idx].OrderDetailID,
					OrderID:  orderID,
					ItemID:   order[idx].ItemID,
					Quantity: order[idx].Quantity,
					Unit:     order[idx].Unit,
					Price:    order[idx].Price,
				})

				orderIDMap[orderID] = len(orders) - 1
			} else {
				orders[orderIDMap[orderID]].OrderDetails = append(orders[orderIDMap[orderID]].OrderDetails, result.OrderDetail{
					ID:       order[idx].OrderDetailID,
					OrderID:  orderID,
					ItemID:   order[idx].ItemID,
					Quantity: order[idx].Quantity,
					Unit:     order[idx].Unit,
					Price:    order[idx].Price,
				})
			}
		}

		allOrders = orders
	}

	return allOrders, nil
}
