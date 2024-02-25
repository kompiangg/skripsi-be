package params

import (
	"github.com/google/uuid"
	"github.com/volatiletech/null/v9"
)

type ShardTimeSeriesWhereQuery struct {
	ShardIndex int
	StartDate  null.Time
	EndDate    null.Time
	CashierID  uuid.NullUUID
}

type FindOrderDetailsOnShardRepo struct {
	ShardIndex int
	OrderID    uuid.NullUUID
}

type LongTermWhereQuery struct {
	StartDate null.Time
	EndDate   null.Time
	CashierID uuid.NullUUID
}

type FindOrderDetailsOnLongTermRepo struct {
	OrderID uuid.NullUUID
}
