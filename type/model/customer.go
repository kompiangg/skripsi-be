package model

import (
	"time"
)

type Customer struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Contact   string    `db:"contact"`
	CreatedAt time.Time `db:"created_at"`
}
