package healthz

import (
	"github.com/labstack/echo/v4"
)

type handler struct {
}

func Init(e *echo.Echo) {
	h := handler{}

	e.GET("/v1/healthz", h.Healthz)
}
