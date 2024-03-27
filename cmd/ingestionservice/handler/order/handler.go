package order

import (
	"skripsi-be/cmd/middleware"
	"skripsi-be/service/order"

	"github.com/labstack/echo/v4"
)

type handler struct {
	e       *echo.Echo
	service order.Service
}

func Init(e *echo.Echo, service order.Service, middleware middleware.Middleware) {
	h := handler{
		e:       e,
		service: service,
	}

	e.POST("/v1/order/without-kappa", h.CreateNewOrderWithoutKappa, middleware.JWTRestricted())
	e.POST("/v1/order", h.CreateNewOrder, middleware.JWTRestricted())
}
