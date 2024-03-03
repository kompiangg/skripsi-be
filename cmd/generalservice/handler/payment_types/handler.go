package payment_types

import (
	"skripsi-be/service/payment_types"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service payment_types.Service
}

func Init(
	e *echo.Echo,
	service payment_types.Service,
) {
	h := handler{
		service: service,
	}

	e.GET("/v1/payment-types", h.FindLikeOneOfAllColumn)
}
