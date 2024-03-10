package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Currency struct {
	ID        string          `db:"id"`
	Base      string          `db:"base"`
	Quote     string          `db:"quote"`
	Rate      decimal.Decimal `db:"rate"`
	CreatedAt time.Time       `db:"created_at"`
}
