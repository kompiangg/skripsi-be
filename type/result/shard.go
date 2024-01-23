package result

import "time"

type ShardTimeSeriesWhereQuery struct {
	ShardIndex int
	StartDate  time.Time
	EndDate    time.Time
}
