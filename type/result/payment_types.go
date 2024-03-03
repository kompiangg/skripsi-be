package result

import (
	"skripsi-be/type/model"
	"time"
)

type PaymentTypes struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Bank      string    `json:"bank"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *PaymentTypes) FromModel(m model.PaymentType) {
	p.ID = m.ID
	p.Type = m.Type
	p.Bank = m.Bank
	p.CreatedAt = m.CreatedAt
}
