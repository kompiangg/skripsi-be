package handler

import (
	"skripsi-be/cmd/orderservice/middleware"
	"skripsi-be/cmd/servingservice/handler/healthz"
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
