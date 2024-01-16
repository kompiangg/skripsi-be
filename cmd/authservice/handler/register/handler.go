package register

import (
	"skripsi-be/cmd/middleware"
	"skripsi-be/service/auth"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service auth.Service
}

func Init(
	echo *echo.Echo,
	service auth.Service,
	middleware middleware.Middleware,
) {
	h := handler{
		service: service,
	}

	echo.POST("/v1/register/admin", h.Admin, middleware.Admin())
	echo.POST("/v1/register/cashier", h.Cashier, middleware.Cashier())
}
