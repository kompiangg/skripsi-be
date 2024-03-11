package model

import (
	"time"
)

type Store struct {
	ID          string    `db:"id"`
	Nation      string    `db:"nation"`
	Region      string    `db:"region"`
	District    string    `db:"district"`
	SubDistrict string    `db:"sub_district"`
	Currency    string    `db:"currency"`
	CreatedAt   time.Time `db:"created_at"`
}
