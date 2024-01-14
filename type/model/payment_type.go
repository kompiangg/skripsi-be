package model

import (
	"time"
)

type PaymentType struct {
	ID        string    `db:"id"`
	Type      string    `db:"type"`
	Bank      string    `db:"bank"`
	CreatedAt time.Time `db:"created_at"`
}
