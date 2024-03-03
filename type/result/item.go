package result

import (
	"skripsi-be/type/model"
	"time"

	"github.com/shopspring/decimal"
)

type Item struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Desc          string          `json:"desc"`
	Price         decimal.Decimal `json:"price"`
	OriginCountry string          `json:"origin_country"`
	Supplier      string          `json:"supplier"`
	Unit          string          `json:"unit"`
	CreatedAt     time.Time       `json:"created_at"`
}

func (i *Item) FromModel(m model.Item) {
	i.ID = m.ID
	i.Name = m.Name
	i.Desc = m.Desc
	i.Price = m.Price
	i.OriginCountry = m.OriginCountry
	i.Supplier = m.Supplier
	i.Unit = m.Unit
	i.CreatedAt = m.CreatedAt
}
