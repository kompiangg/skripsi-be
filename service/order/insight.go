package order

import (
	"context"
	"fmt"
	"math"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/model"
	"skripsi-be/type/result"
	"sort"
	"sync"
	"time"
)

func (s service) FindInsightBasedOnInterval(ctx context.Context, interval string, offset int) (result.GetAggregateOrderService, error) {
	res := result.GetAggregateOrderService{
		TopProducts: make([]result.GetAggregateOrderTopProductService, 0),
		Chart:       make([]result.GetAggregateOrderChartService, 0),
	}

	decrementBy, exists := getDecrementBy(interval)
	if !exists {
		return res, errors.Wrap(errors.ErrBadRequest)
	}

	loc := time.FixedZone(fmt.Sprintf("UTC+%d", offset), offset*3600)

	const divideBy = 7
	var dayRanges [divideBy]struct {
		StartDate time.Time
		EndDate   time.Time
	}

	dayRange := math.Ceil(float64(decrementBy) / divideBy)

	for i := 0; i < divideBy; i++ {
		startDate := s.config.Date.Now().In(loc).AddDate(0, 0, -(int(dayRange)*i + int(dayRange) - 1))
		endDate := s.config.Date.Now().In(loc).AddDate(0, 0, -int(dayRange)*i)

		dayRanges[i].StartDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
		dayRanges[i].EndDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), endDate.Hour(), endDate.Minute(), endDate.Second(), endDate.Nanosecond(), loc)
	}

	globalStartRange := dayRanges[divideBy-1].StartDate
	globalEndRange := dayRanges[0].EndDate

	if s.config.IsUsingSharding {
		mutex := sync.Mutex{}

		whereQuery, err := s.getShardWhereQuery(globalStartRange, globalEndRange)
		if err != nil {
			return res, errors.Wrap(err)
		}

		errorChan := make(chan error, len(whereQuery))
		defer close(errorChan)

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		mapShardTopSellingProducts := make(map[string]model.GetAggregateTopSellingProductResultRepo)

		for i := range whereQuery {
			go func(shardIndex int, startDate, endDate time.Time) {
				var topSellingProducts []model.GetAggregateTopSellingProductResultRepo = []model.GetAggregateTopSellingProductResultRepo{}

				if shardIndex == len(s.config.Shards) {
					topSellingProducts, err = s.orderRepo.GetAggregateTopSellingProductOnLongTermDB(ctx, startDate, endDate)
					if err != nil {
						errorChan <- errors.Wrap(err)
						return
					}
				} else {
					topSellingProducts, err = s.orderRepo.GetAggregateTopSellingProductOnShardDB(ctx, shardIndex, startDate, endDate)
					if err != nil {
						errorChan <- errors.Wrap(err)
						return
					}
				}

				mutex.Lock()
				for idx := range topSellingProducts {
					mapShardTopSellingProducts[topSellingProducts[idx].ItemID] = model.GetAggregateTopSellingProductResultRepo{
						ItemID:                topSellingProducts[idx].ItemID,
						ItemSoldTotalQuantity: mapShardTopSellingProducts[topSellingProducts[idx].ItemID].ItemSoldTotalQuantity + topSellingProducts[idx].ItemSoldTotalQuantity,
					}
				}
				mutex.Unlock()

				errorChan <- nil
			}(whereQuery[i].ShardIndex, whereQuery[i].StartDate, whereQuery[i].EndDate)
		}

		notErrCount := 0
		for notErrCount != len(whereQuery) {
			err := <-errorChan
			if err != nil {
				return res, errors.Wrap(err)
			}

			notErrCount++
		}

		topProducts := make([]result.GetAggregateOrderTopProductService, 0, len(mapShardTopSellingProducts))
		for _, v := range mapShardTopSellingProducts {
			topProducts = append(topProducts, result.GetAggregateOrderTopProductService{
				ItemID:                v.ItemID,
				ItemSoldTotalQuantity: v.ItemSoldTotalQuantity,
			})
		}

		// sort top products by descending
		sort.Slice(topProducts, func(i, j int) bool {
			return topProducts[i].ItemSoldTotalQuantity > topProducts[j].ItemSoldTotalQuantity
		})

		if len(topProducts) > 3 {
			res.TopProducts = topProducts[:3]
		} else {
			res.TopProducts = topProducts
		}

		res.Chart = make([]result.GetAggregateOrderChartService, len(dayRanges))

		for i := range dayRanges {
			whereQuery, err := s.getShardWhereQuery(dayRanges[i].StartDate, dayRanges[i].EndDate)
			if err != nil {
				return res, errors.Wrap(err)
			}

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			errorChan := make(chan error, len(whereQuery))
			defer close(errorChan)

			chart := result.GetAggregateOrderChartService{
				Date: dayRanges[i].EndDate,
			}

			for i := range whereQuery {
				go func(shardIndex int, startDate, endDate time.Time) {
					var aggregateOrderMember []model.GetAggregateOrderResultRepo

					if shardIndex == len(s.config.Shards) {
						aggregateOrderMember, err = s.orderRepo.GetAggregateOrderOnLongTermDB(ctx, startDate, endDate)
						if errors.Is(err, context.Canceled) {
							errorChan <- nil
							return
						} else if errors.Is(err, errors.ErrNotFound) {
							err = nil
						} else if err != nil {
							errorChan <- errors.Wrap(err)
							return
						}
					} else {
						aggregateOrderMember, err = s.orderRepo.GetAggregateOrderOnShardDB(ctx, shardIndex, startDate, endDate)
						if errors.Is(err, context.Canceled) {
							errorChan <- nil
							return
						} else if errors.Is(err, errors.ErrNotFound) {
							err = nil
						} else if err != nil {
							errorChan <- errors.Wrap(err)
							return
						}
					}

					mutex.Lock()
					customerOrderMap := make(map[string]model.GetAggregateOrderResultRepo)
					for _, v := range aggregateOrderMember {
						customerOrderMap[v.CustomerID] = model.GetAggregateOrderResultRepo{
							CustomerID:             v.CustomerID,
							OrderQuantity:          customerOrderMap[v.CustomerID].OrderQuantity + v.OrderQuantity,
							NotMemberOrderQuantity: customerOrderMap[v.CustomerID].NotMemberOrderQuantity + v.NotMemberOrderQuantity,
							ItemSoldTotalQuantity:  customerOrderMap[v.CustomerID].ItemSoldTotalQuantity + v.ItemSoldTotalQuantity,
							ItemSoldTotalPrice:     customerOrderMap[v.CustomerID].ItemSoldTotalPrice.Add(v.ItemSoldTotalPrice),
						}
					}

					for _, v := range customerOrderMap {
						if v.CustomerID == "null" {
							res.TotalNotCustomerOrderQuantity += v.NotMemberOrderQuantity
						} else {
							res.TotalCustomerOrderQuantity += v.OrderQuantity
						}

						res.ItemSoldTotalQuantity += v.ItemSoldTotalQuantity
						chart.TotalOrderQuantity += v.ItemSoldTotalQuantity

						res.ItemSoldTotalPrice = res.ItemSoldTotalPrice.Add(v.ItemSoldTotalPrice)
						chart.TotalOrderPrice = chart.TotalOrderPrice.Add(v.ItemSoldTotalPrice)
					}

					mutex.Unlock()

					errorChan <- nil
				}(whereQuery[i].ShardIndex, whereQuery[i].StartDate, whereQuery[i].EndDate)
			}

			errCount := 0
			var routineErr error

			for errCount != len(whereQuery) {
				routineErr = errors.Join(routineErr, <-errorChan)
				if routineErr != nil {
					cancel()
				}

				errCount++
			}

			if routineErr != nil {
				return res, errors.Wrap(routineErr)
			}

			res.Chart[i] = chart
			fmt.Println("startDate", dayRanges[i].StartDate)
			fmt.Println("endDate", dayRanges[i].EndDate)
		}
	} else {
		topSellingProducts, err := s.orderRepo.GetAggregateTopSellingProductOnLongTermDB(ctx, globalStartRange, globalEndRange)
		if err != nil {
			return res, errors.Wrap(err)
		}

		topSellingProductsLength := 3
		if len(topSellingProducts) < 3 {
			topSellingProductsLength = len(topSellingProducts)
		}

		for idx := 0; idx < topSellingProductsLength; idx++ {
			res.TopProducts = append(res.TopProducts, result.GetAggregateOrderTopProductService{
				ItemID:                topSellingProducts[idx].ItemID,
				ItemSoldTotalQuantity: topSellingProducts[idx].ItemSoldTotalQuantity,
			})
		}

		for i := range dayRanges {
			chart := result.GetAggregateOrderChartService{
				Date: dayRanges[i].EndDate,
			}

			aggregateOrderMember, err := s.orderRepo.GetAggregateOrderOnLongTermDB(ctx, dayRanges[i].StartDate, dayRanges[i].EndDate)
			if err != nil {
				return res, errors.Wrap(err)
			}

			customerOrderMap := make(map[string]model.GetAggregateOrderResultRepo)
			for _, v := range aggregateOrderMember {
				customerOrderMap[v.CustomerID] = model.GetAggregateOrderResultRepo{
					CustomerID:             v.CustomerID,
					OrderQuantity:          customerOrderMap[v.CustomerID].OrderQuantity + v.OrderQuantity,
					NotMemberOrderQuantity: customerOrderMap[v.CustomerID].NotMemberOrderQuantity + v.NotMemberOrderQuantity,
					ItemSoldTotalQuantity:  customerOrderMap[v.CustomerID].ItemSoldTotalQuantity + v.ItemSoldTotalQuantity,
					ItemSoldTotalPrice:     customerOrderMap[v.CustomerID].ItemSoldTotalPrice.Add(v.ItemSoldTotalPrice),
				}
			}

			for _, v := range customerOrderMap {
				if v.CustomerID == "null" {
					res.TotalNotCustomerOrderQuantity += v.NotMemberOrderQuantity
				} else {
					res.TotalCustomerOrderQuantity += v.OrderQuantity
				}

				res.ItemSoldTotalQuantity += v.ItemSoldTotalQuantity
				chart.TotalOrderQuantity += v.ItemSoldTotalQuantity

				res.ItemSoldTotalPrice = res.ItemSoldTotalPrice.Add(v.ItemSoldTotalPrice)
				chart.TotalOrderPrice = chart.TotalOrderPrice.Add(v.ItemSoldTotalPrice)
			}

			res.Chart = append(res.Chart, chart)
		}
	}

	return res, nil
}

func getDecrementBy(interval string) (int, bool) {
	switch interval {
	case "1w":
		return 7, true
	case "1m":
		return 30, true
	case "3m":
		return 90, true
	default:
		return 0, false
	}
}
