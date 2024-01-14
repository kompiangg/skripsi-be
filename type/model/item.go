package model

import (
	"time"
)

type Item struct {
	ID            string    `db:"id"`
	Name          string    `db:"name"`
	Desc          string    `db:"desc"`
	Price         float64   `db:"price"`
	OriginCountry string    `db:"origin_country"`
	Supplier      string    `db:"supplier"`
	Unit          string    `db:"unit"`
	CreatedAt     time.Time `db:"created_at"`
}
