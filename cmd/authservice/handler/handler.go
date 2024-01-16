package handler

import (
	"skripsi-be/cmd/authservice/handler/healthz"
	"skripsi-be/cmd/authservice/handler/register"
	"skripsi-be/cmd/authservice/handler/session"
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
	session.Init(echo, service.Auth)
	register.Init(echo, service.Auth, middleware)
}
