package result

import (
	"skripsi-be/type/model"
	"time"
)

type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Contact   string    `json:"contact"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Customer) FromModel(m model.Customer) {
	c.ID = m.ID
	c.Name = m.Name
	c.Contact = m.Contact
	c.CreatedAt = m.CreatedAt
}
