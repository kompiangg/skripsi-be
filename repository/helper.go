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

func beginGeneralDBTx(generalDB *sqlx.DB) func(ctx context.Context) (*sqlx.Tx, error) {
	return func(ctx context.Context) (*sqlx.Tx, error) {
		return generalDB.BeginTxx(context.Background(), nil)
	}
}

func getShardIndexByDateTime(shards config.Shards, customDate config.Date) func(date time.Time) (int, error) {
	now := customDate.Now()

	return func(date time.Time) (int, error) {
		now = now.In(date.Location())

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
		whereQueries := []result.ShardTimeSeriesWhereQuery{}

		startDateOriginTimeZone := startDate.Location()
		endDateOriginTimeZone := endDate.Location()

		now := customDate.Now().In(startDate.Location())

		diffStart := now.Sub(startDate)
		diffStartInDay := int(diffStart.Hours() / 24)

		diffEnd := now.Sub(endDate)
		diffEndInDay := int(diffEnd.Hours() / 24)

		startIdx, endIdx := -1, -1 // -1 means not found

		if diffEndInDay < subtractedShard[0].DataRetention {
			return nil, constant.ErrOutOfShardRange
		}

		if diffEndInDay > subtractedShard[len(subtractedShard)-2].DataRetention {
			whereQueries = append(whereQueries, result.ShardTimeSeriesWhereQuery{
				ShardIndex: len(subtractedShard) - 1,
				StartDate:  startDate.In(startDateOriginTimeZone),
				EndDate:    endDate.In(endDateOriginTimeZone),
			})
			return whereQueries, nil
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
		startDateQuery := time.Time{}
		endDateQuery := time.Time{}

		if len(choosedDB) == 1 {
			startDateQuery = now.Add(-time.Duration(time.Nanosecond * diffStart))
			endDateQuery = now.Add(-time.Duration(time.Nanosecond * diffEnd))

			whereQueries = append(whereQueries, result.ShardTimeSeriesWhereQuery{
				StartDate:  startDateQuery.In(startDateOriginTimeZone),
				EndDate:    endDateQuery.In(endDateOriginTimeZone),
				ShardIndex: endIdx,
			})

			return whereQueries, nil
		}

		for i, v := range choosedDB {
			durationInDay := time.Duration(v.DataRetention+1) * 24 * time.Hour
			dayInNanoSecond := durationInDay.Nanoseconds()

			var oneIndexBeforeDayInNanoSecond int64
			if i != 0 {
				oneIndexBeforeDurationInDay := time.Duration(choosedDB[i-1].DataRetention+1) * 24 * time.Hour
				oneIndexBeforeDayInNanoSecond = oneIndexBeforeDurationInDay.Nanoseconds()
			}

			if i == 0 {
				startDateQuery = now.Add(-time.Duration(dayInNanoSecond+1) * time.Nanosecond)
				endDateQuery = now.Add(-time.Duration(time.Nanosecond * diffEnd))
			} else if i == len(choosedDB)-1 {
				startDateQuery = now.Add(-time.Duration(time.Nanosecond * (diffStart + 1)))
				endDateQuery = now.Add(-time.Duration(oneIndexBeforeDayInNanoSecond) * time.Nanosecond)
			} else {
				startDateQuery = now.Add(-time.Duration(dayInNanoSecond+1) * time.Nanosecond)
				endDateQuery = now.Add(-time.Duration(oneIndexBeforeDayInNanoSecond) * time.Nanosecond)
			}

			whereQueries = append(whereQueries, result.ShardTimeSeriesWhereQuery{
				StartDate: startDateQuery.In(startDateOriginTimeZone),
				EndDate:   endDateQuery.In(endDateOriginTimeZone),
			})
		}

		for i := range choosedDB {
			whereQueries[i].ShardIndex = endIdx + i
		}

		return whereQueries, nil
	}
}
