package handler

import (
	"skripsi-be/cmd/middleware"
	"skripsi-be/cmd/servingservice/handler/healthz"
	"skripsi-be/cmd/servingservice/handler/orders"
	"skripsi-be/config"
	"skripsi-be/service"

	"github.com/labstack/echo/v4"
)

func Init(
	echo *echo.Echo,
	service service.Service,
	middleware middleware.Middleware,
	config config.Config,
) {
	healthz.Init(echo)
	orders.Init(echo, middleware, config, service.Order)
}
