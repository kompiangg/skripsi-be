package model

import (
	"time"
)

type Store struct {
	ID          string    `db:"id"`
	Region      string    `db:"region"`
	District    string    `db:"district"`
	SubDistrict string    `db:"sub_district"`
	CreatedAt   time.Time `db:"created_at"`
}
