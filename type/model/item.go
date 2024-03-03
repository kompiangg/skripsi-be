package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Item struct {
	ID            string          `db:"id"`
	Name          string          `db:"name"`
	Desc          string          `db:"desc"`
	Price         decimal.Decimal `db:"price"`
	OriginCountry string          `db:"origin_country"`
	Supplier      string          `db:"supplier"`
	Unit          string          `db:"unit"`
	CreatedAt     time.Time       `db:"created_at"`
}
