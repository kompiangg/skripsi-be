package handler

import (
	"skripsi-be/cmd/middleware"
	"skripsi-be/cmd/shardingloadservice/handler/healthz"
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
