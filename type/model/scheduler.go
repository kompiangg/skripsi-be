package model

import (
	"time"

	"github.com/volatiletech/null/v9"
)

type Scheduler struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	RunCount  int       `db:"run_count"`
	LastRunAt null.Time `db:"last_run_at"`
	CreatedAt time.Time `db:"created_at"`
}
