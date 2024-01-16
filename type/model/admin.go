package model

import (
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	ID        uuid.UUID `db:"id"`
	AccountID uuid.UUID `db:"account_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func (a *Admin) AssignUUID() bool {
	if a.ID != uuid.Nil {
		return false
	}

	a.ID = uuid.New()
	return true
}
