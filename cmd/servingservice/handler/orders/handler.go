package orders

import (
	"skripsi-be/cmd/middleware"
	"skripsi-be/config"
	"skripsi-be/service/order"

	"github.com/labstack/echo/v4"
)

type handler struct {
	orderService order.Service
	config       config.Config
}

func Init(
	e *echo.Echo,
	middleware middleware.Middleware,
	config config.Config,
	orderService order.Service,
) {
	h := handler{
		orderService: orderService,
		config:       config,
	}

	e.GET("/v1/orders-details", h.GetAllDetails, middleware.JWTRestricted())
	e.GET("/v1/orders", h.GetBriefInformation, middleware.JWTRestricted())
	e.GET("/v1/orders/:id", h.GetOrderDetails, middleware.JWTRestricted())
}
