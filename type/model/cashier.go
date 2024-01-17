package model

import (
	"time"

	"github.com/google/uuid"
)

type Cashier struct {
	ID        uuid.UUID `db:"id"`
	AccountID uuid.UUID `db:"account_id"`
	StoreID   string    `db:"store_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func (c *Cashier) AssignUUID() bool {
	if c.ID != uuid.Nil {
		return false
	}

	c.ID = uuid.New()
	return true
}
