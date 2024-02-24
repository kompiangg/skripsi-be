package orders

import (
	"skripsi-be/cmd/middleware"
	"skripsi-be/config"
	"skripsi-be/service/order"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service order.Service
	config  config.Config
}

func Init(
	e *echo.Echo,
	middleware middleware.Middleware,
	config config.Config,
	service order.Service,
) {
	h := handler{
		service: service,
		config:  config,
	}

	e.GET("/v1/orders-details", h.GetAllDetails, middleware.JWTRestricted())
	e.GET("/v1/orders", h.GetAllDetails, middleware.JWTRestricted())
	e.GET("/v1/orders/:id/:", h.GetAllDetails, middleware.JWTRestricted())
}
