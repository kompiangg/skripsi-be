package session

import (
	"skripsi-be/service/auth"

	"github.com/labstack/echo/v4"
)

type handler struct {
	authService auth.Service
}

func Init(e *echo.Echo, service auth.Service) {
	h := handler{
		authService: service,
	}

	e.POST("/v1/session", h.NewSession)
}
