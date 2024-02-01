package repository

import (
	"context"
	"math"
	"skripsi-be/config"
	"skripsi-be/type/constant"
	"skripsi-be/type/result"
	"time"

	"github.com/jmoiron/sqlx"
)

func beginLongTermDBTx(longtermDB *sqlx.DB) func(ctx context.Context) (*sqlx.Tx, error) {
	return func(ctx context.Context) (*sqlx.Tx, error) {
		return longtermDB.BeginTxx(context.Background(), nil)
	}
}

func beginShardDBTx(shardingDatabase []*sqlx.DB) func(ctx context.Context, dbIndex int) (*sqlx.Tx, error) {
	return func(ctx context.Context, dbIndex int) (*sqlx.Tx, error) {
		return shardingDatabase[dbIndex].BeginTxx(ctx, nil)
	}
}

func getShardIndexByDateTime(shards config.Shards, customDate config.Date) func(date time.Time) (int, error) {
	now := customDate.Now()

	return func(date time.Time) (int, error) {
		date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		diff := now.Sub(date)
		diffInDay := int(diff.Hours() / 24)

		for i, v := range shards {
			if diffInDay-v.DataRetention < 0 {
				return i, nil
			}
		}

		return 0, constant.ErrOutOfShardRange
	}
}

// Start date should smaller than end date
// Example:
// Start date: 2021-01-01
// End date: 2021-01-02
func getShardWhereQuery(shards config.Shards, customDate config.Date) func(startDate time.Time, endDate time.Time) ([]result.ShardTimeSeriesWhereQuery, error) {
	subtractedShard := make([]config.Shard, len(shards))
	for i, v := range shards {
		subtractedShard[i] = config.Shard{
			DataRetention: v.DataRetention - 1,
			URIConnection: v.URIConnection,
		}
	}

	subtractedShard = append(subtractedShard, config.Shard{
		DataRetention: math.MaxInt64,
	})

	return func(startDate time.Time, endDate time.Time) ([]result.ShardTimeSeriesWhereQuery, error) {
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())
		now := time.Date(customDate.Now().Year(), customDate.Now().Month(), customDate.Now().Day(), 0, 0, 0, 0, customDate.Now().Location())

		diffStart := now.Sub(startDate)
		diffStartInDay := int(diffStart.Hours() / 24)

		diffEnd := now.Sub(endDate)
		diffEndInDay := int(diffEnd.Hours() / 24)

		startIdx, endIdx := -1, -1 // -1 means not found

		if diffEndInDay < subtractedShard[0].DataRetention {
			return nil, constant.ErrOutOfShardRange
		}

		if diffEndInDay > subtractedShard[len(subtractedShard)-2].DataRetention {
			return nil, constant.ErrOutOfShardRange
		}

		for i, v := range subtractedShard {
			if diffStartInDay-v.DataRetention <= 0 && startIdx == -1 {
				startIdx = i
			}

			if diffEndInDay-v.DataRetention <= 0 && endIdx == -1 {
				endIdx = i
			}

			if startIdx != -1 && endIdx != -1 {
				break
			}
		}

		choosedDB := subtractedShard[endIdx : startIdx+1]
		whereQueries := []result.ShardTimeSeriesWhereQuery{}
		startDateQuery := time.Time{}
		endDateQuery := time.Time{}

		if len(choosedDB) == 1 {
			startDateQuery = now.Add(-time.Duration(int(time.Hour) * 24 * diffStartInDay))
			endDateQuery = now.Add(-time.Duration(int(time.Hour) * 24 * diffEndInDay))

			whereQueries = append(whereQueries, result.ShardTimeSeriesWhereQuery{
				StartDate:  time.Date(startDateQuery.Year(), startDateQuery.Month(), startDateQuery.Day(), 0, 0, 0, 0, startDateQuery.Location()),
				EndDate:    time.Date(endDateQuery.Year(), endDateQuery.Month(), endDateQuery.Day(), 23, 59, 59, 0, endDateQuery.Location()),
				ShardIndex: endIdx,
			})

			return whereQueries, nil
		}

		for i, v := range choosedDB {
			if i == 0 {
				startDateQuery = now.Add(-time.Duration(int(time.Hour) * 24 * v.DataRetention))
				endDateQuery = now.Add(-time.Duration(int(time.Hour) * 24 * diffEndInDay))
			} else if i == len(choosedDB)-1 {
				startDateQuery = now.Add(-time.Duration(int(time.Hour) * 24 * diffStartInDay))
				endDateQuery = now.Add(-time.Duration(int(time.Hour) * 24 * (choosedDB[i-1].DataRetention + 1)))
			} else {
				startDateQuery = now.Add(-time.Duration(int(time.Hour) * 24 * v.DataRetention))
				endDateQuery = now.Add(-time.Duration(int(time.Hour) * 24 * (choosedDB[i-1].DataRetention + 1)))
			}

			whereQueries = append(whereQueries, result.ShardTimeSeriesWhereQuery{
				StartDate: time.Date(startDateQuery.Year(), startDateQuery.Month(), startDateQuery.Day(), 0, 0, 0, 0, startDateQuery.Location()),
				EndDate:   time.Date(endDateQuery.Year(), endDateQuery.Month(), endDateQuery.Day(), 23, 59, 59, 0, endDateQuery.Location()),
			})
		}

		for i := range choosedDB {
			whereQueries[i].ShardIndex = endIdx + i
		}

		return whereQueries, nil
	}
}
