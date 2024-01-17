package handler

import (
	"skripsi-be/cmd/orderservice/handler/healthz"
	"skripsi-be/cmd/orderservice/middleware"
	"skripsi-be/service"

	"github.com/labstack/echo/v4"
)

func Init(
	echo *echo.Echo,
	service service.Service,
	middleware middleware.Middleware,
) {
	healthz.Init(echo)
}
