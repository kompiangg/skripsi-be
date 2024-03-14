package order

import (
	"context"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/constant"
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

	if s.config.IsUsingSharding {
		// Why i did this?
		// param.StartDate = time.Date(param.StartDate.Year(), param.StartDate.Month(), param.StartDate.Day(), 0, 0, 0, 0, param.StartDate.Location())
		// param.EndDate = time.Date(param.EndDate.Year(), param.EndDate.Month(), param.EndDate.Day(), 23, 59, 59, 0, param.EndDate.Location())

		shardQuery, err := s.getShardWhereQuery(param.StartDate, param.EndDate)
		if err != nil {
			return nil, errors.Wrap(err)
		}

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
					if errors.Is(err, context.Canceled) {
						errorChan <- nil
						return
					} else if err != nil {
						errorChan <- err
						return
					}
				} else {
					order, err = s.orderRepo.FindAllOrderAndDetailsOnShardDB(ctx, params.ShardTimeSeriesWhereQuery{
						ShardIndex: shardParam.ShardIndex,
						StartDate:  null.TimeFrom(shardParam.StartDate),
						EndDate:    null.TimeFrom(shardParam.EndDate),
					})
					if errors.Is(err, context.Canceled) {
						errorChan <- nil
						return
					} else if err != nil {
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

		errCount := 0
		var routineErr error

		for errCount != len(shardQuery) {
			routineErr = errors.Join(routineErr, <-errorChan)
			if routineErr != nil {
				cancel()
			}

			errCount++
		}

		if routineErr != nil {
			return nil, errors.Wrap(routineErr)
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

func (s service) FindBriefInformationOrder(ctx context.Context, param params.FindOrderService) (orders []result.OrderBriefInformation, err error) {
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

	orders = make([]result.OrderBriefInformation, 0)

	if s.config.IsUsingSharding {
		shardQuery, err := s.getShardWhereQuery(param.StartDate, param.EndDate)
		if err != nil {
			return nil, errors.Wrap(err)
		}

		errorChan := make(chan error, len(shardQuery))
		defer close(errorChan)

		shardOrders := make([][]result.OrderBriefInformation, len(shardQuery))
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		for idx, shardParam := range shardQuery {
			go func(idx int, shardParam result.ShardTimeSeriesWhereQuery) {
				var order []model.Order

				if shardParam.ShardIndex == len(s.config.Shards) {
					order, err = s.orderRepo.FindAllOnLongTermDB(ctx, params.LongTermWhereQuery{
						StartDate: null.TimeFrom(shardParam.StartDate),
						EndDate:   null.TimeFrom(shardParam.EndDate),
						CashierID: param.CashierID,
					})
					if errors.Is(err, context.Canceled) {
						errorChan <- nil
						return
					} else if err != nil {
						errorChan <- err
						return
					}
				} else {
					order, err = s.orderRepo.FindAllOnShardDB(ctx, params.ShardTimeSeriesWhereQuery{
						ShardIndex: shardParam.ShardIndex,
						StartDate:  null.TimeFrom(shardParam.StartDate),
						EndDate:    null.TimeFrom(shardParam.EndDate),
						CashierID:  param.CashierID,
					})
					if errors.Is(err, context.Canceled) {
						errorChan <- nil
						return
					} else if err != nil {
						errorChan <- err
						return
					}
				}

				orders := make([]result.OrderBriefInformation, 0)

				for idx := range order {
					orderID, err := ulid.Parse(order[idx].ID)
					if err != nil {
						errorChan <- err
						return
					}

					orders = append(orders, result.OrderBriefInformation{
						ID:              orderID,
						CashierID:       order[idx].CashierID,
						StoreID:         order[idx].StoreID,
						PaymentID:       order[idx].PaymentID,
						CustomerID:      order[idx].CustomerID,
						TotalQuantity:   order[idx].TotalQuantity,
						TotalPrice:      order[idx].TotalPrice,
						TotalPriceInUSD: order[idx].TotalPriceInUSD,
						Currency:        order[idx].Currency,
						UsdRate:         order[idx].UsdRate,
						CreatedAt:       order[idx].CreatedAt.In(param.StartDate.Location()),
					})
				}

				shardOrders[idx] = orders
				errorChan <- nil
			}(idx, shardParam)
		}

		errCount := 0
		var routineErr error

		for errCount != len(shardQuery) {
			routineErr = errors.Join(routineErr, <-errorChan)
			if routineErr != nil {
				cancel()
			}

			errCount++
		}

		if routineErr != nil {
			return nil, errors.Wrap(routineErr)
		}

		for _, shardOrder := range shardOrders {
			orders = append(orders, shardOrder...)
		}
	} else {
		order, err := s.orderRepo.FindAllOnLongTermDB(ctx, params.LongTermWhereQuery{
			StartDate: null.TimeFrom(param.StartDate),
			EndDate:   null.TimeFrom(param.EndDate),
			CashierID: param.CashierID,
		})
		if err != nil {
			return nil, errors.Wrap(err)
		}

		for idx := range order {
			orderID, err := ulid.Parse(order[idx].ID)
			if err != nil {
				return nil, errors.Wrap(err)
			}

			orders = append(orders, result.OrderBriefInformation{
				ID:              orderID,
				CashierID:       order[idx].CashierID,
				StoreID:         order[idx].StoreID,
				PaymentID:       order[idx].PaymentID,
				CustomerID:      order[idx].CustomerID,
				TotalQuantity:   order[idx].TotalQuantity,
				TotalPrice:      order[idx].TotalPrice,
				TotalPriceInUSD: order[idx].TotalPriceInUSD,
				Currency:        order[idx].Currency,
				UsdRate:         order[idx].UsdRate,
				CreatedAt:       order[idx].CreatedAt.In(param.StartDate.Location()),
			})
		}
	}

	return orders, nil
}

func (s service) FindOrderDetails(ctx context.Context, param params.FindOrderDetailsService) (res result.Order, err error) {
	err = param.Validate()
	if err != nil {
		return res, errors.Wrap(err)
	}

	orderID, err := ulid.Parse(param.OrderID)
	if err != nil {
		return res, errors.Wrap(err)
	}

	isUsingSharding := true
	createdAtUTC := time.Unix(int64(orderID.Time())/1000, 0).UTC()
	shardIndex, err := s.getShardIndexByDateTime(createdAtUTC)
	if errors.Is(err, constant.ErrOutOfShardRange) {
		err = nil
		isUsingSharding = false
		shardIndex = -1
	} else if err != nil {
		return res, errors.Wrap(err)
	}

	if s.config.IsUsingSharding && isUsingSharding {
		order, err := s.orderRepo.FindOrderByIDOnShardDB(ctx, shardIndex, orderID.String())
		if errors.Is(err, errors.ErrNotFound) {
			return res, errors.Wrap(errors.ErrNotFound)
		} else if err != nil {
			return res, errors.Wrap(err)
		}

		orderDetails, err := s.orderRepo.FindOrderDetailsByOrderIDOnShardDB(ctx, shardIndex, orderID.String())
		if err != nil {
			return res, errors.Wrap(err)
		}

		res.OrderDetails = make([]result.OrderDetail, len(orderDetails))
		order.OrderDetails = orderDetails
		res.FromModel(order)
	} else {
		order, err := s.orderRepo.FindOrderByIDOnLongTermDB(ctx, orderID.String())
		if errors.Is(err, errors.ErrNotFound) {
			return res, errors.Wrap(errors.ErrNotFound)
		} else if err != nil {
			return res, errors.Wrap(err)
		}

		orderDetails, err := s.orderRepo.FindOrderDetailsByOrderIDOnLongTermDB(ctx, orderID.String())
		if err != nil {
			return res, errors.Wrap(err)
		}

		res.OrderDetails = make([]result.OrderDetail, len(orderDetails))
		order.OrderDetails = orderDetails
		res.FromModel(order)
	}

	return res, nil
}
