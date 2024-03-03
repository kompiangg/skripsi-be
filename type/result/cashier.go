package result

import (
	"time"

	"github.com/google/uuid"
)

type Cashier struct {
	ID        uuid.UUID `json:"id"`
	StoreID   string    `json:"store_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
