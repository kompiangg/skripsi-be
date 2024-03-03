package item

import (
	"skripsi-be/service/item"

	"github.com/labstack/echo/v4"
)

type handler struct {
	itemService item.Service
}

func Init(
	e *echo.Echo,
	itemService item.Service,
) {
	h := handler{
		itemService: itemService,
	}

	e.GET("/v1/items", h.FindLikeNameOrID)
}
