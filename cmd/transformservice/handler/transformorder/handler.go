package transformorder

import (
	"skripsi-be/cmd/middleware"
	"skripsi-be/config"
	"skripsi-be/service/order"

	"github.com/labstack/echo/v4"
)

type eventHandler struct {
	orderService order.Service
}

func New(
	orderService order.Service,
) eventHandler {
	return eventHandler{
		orderService: orderService,
	}
}

type httpHandler struct {
	orderService order.Service
	config       config.Config
}

func InitHTTPHandler(
	e *echo.Echo,
	middleware middleware.Middleware,
	config config.Config,
	orderService order.Service,
) {
	h := httpHandler{
		orderService: orderService,
		config:       config,
	}

	e.POST("/v1/transform", h.TransformOrder)
}
