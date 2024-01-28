package order

import (
	"context"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/constant"
	"skripsi-be/type/params"
	"skripsi-be/type/result"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v9"
)

func (s service) FindOrder(ctx context.Context, param params.FindOrderService) (allOrders []result.Order, err error) {
	isUsingSharding := true

	shardingNichBos, ok := ctx.Value("isUsingSharding").(bool)
	if ok {
		if !shardingNichBos {
			isUsingSharding = false
		}
	}

	param.StartDate = time.Date(param.StartDate.Year(), param.StartDate.Month(), param.StartDate.Day(), 0, 0, 0, 0, param.StartDate.Location())
	param.EndDate = time.Date(param.EndDate.Year(), param.EndDate.Month(), param.EndDate.Day(), 23, 59, 59, 0, param.EndDate.Location())

	shardQuery, err := s.getShardWhereQuery(param.StartDate, param.EndDate)
	if errors.Is(err, constant.ErrOutOfShardRange) {
		isUsingSharding = false
		err = nil
	} else if err != nil {
		return nil, errors.Wrap(err)
	}

	if s.config.IsUsingSharding && isUsingSharding {
		var waitGroup sync.WaitGroup

		errorChan := make(chan error, 1)
		shardOrder := make([][]result.Order, len(shardQuery))
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		for _, shardParam := range shardQuery {
			waitGroup.Add(1)
			go func(shardParam result.ShardTimeSeriesWhereQuery) {
				defer waitGroup.Done()

				order, err := s.orderRepo.FindAllOnShardDB(ctx, params.ShardTimeSeriesWhereQuery{
					ShardIndex: shardParam.ShardIndex,
					StartDate:  null.TimeFrom(shardParam.StartDate),
					EndDate:    null.TimeFrom(shardParam.EndDate),
				})
				if err != nil {
					select {
					case errorChan <- err: // send the error to the error channel
					default: // if the error channel is already full, don't block
					}
					cancel() // cancel the context
					return
				}

				for idx := range order {
					order[idx].OrderDetails, err = s.orderRepo.FindOrderDetailsOnShardDB(ctx, params.FindOrderDetailsOnShardRepo{
						ShardIndex: shardParam.ShardIndex,
						OrderID:    uuid.NullUUID{UUID: order[idx].ID, Valid: true},
					})
					if err != nil {
						select {
						case errorChan <- err: // send the error to the error channel
						default: // if the error channel is already full, don't block
						}
						cancel() // cancel the context
						return
					}
				}

				orders := make([]result.Order, len(order))
				for idx := range order {
					orders[idx].FromModel(order[idx])
				}

				shardOrder[shardParam.ShardIndex] = orders
			}(shardParam)
		}

		waitGroup.Wait()
		close(errorChan)

		for _, orders := range shardOrder {
			allOrders = append(allOrders, orders...)
		}
	} else {
		order, err := s.orderRepo.FindAllOnLongTermDB(ctx, params.LongTermWhereQuery{
			StartDate: null.TimeFrom(param.StartDate),
			EndDate:   null.TimeFrom(param.EndDate),
		})
		if err != nil {
			return nil, errors.Wrap(err)
		}

		for idx := range order {
			order[idx].OrderDetails, err = s.orderRepo.FindOrderDetailsOnLongTermDB(ctx, params.FindOrderDetailsOnLongTermRepo{
				OrderID: uuid.NullUUID{UUID: order[idx].ID, Valid: true},
			})
			if err != nil {
				return nil, errors.Wrap(err)
			}
		}

		allOrders = make([]result.Order, len(order))
		for idx := range order {
			allOrders[idx].FromModel(order[idx])
		}
	}

	return allOrders, nil
}
