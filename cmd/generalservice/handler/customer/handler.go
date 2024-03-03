package customer

import (
	"skripsi-be/service/customer"

	"github.com/labstack/echo/v4"
)

type handler struct {
	customerService customer.Service
}

func Init(
	e *echo.Echo,
	customerService customer.Service,
) {
	h := handler{
		customerService: customerService,
	}

	e.GET("/v1/customers", h.FindLikeOneOfAllColumn)
}
