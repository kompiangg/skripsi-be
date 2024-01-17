package handler

import (
	"skripsi-be/cmd/longtermloadservice/handler/healthz"
	"skripsi-be/cmd/middleware"
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
