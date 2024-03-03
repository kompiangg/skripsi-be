package healthz

import (
	"net/http"
	"skripsi-be/lib/httpx"

	"github.com/labstack/echo/v4"
)

type handler struct {
}

func Init(e *echo.Echo) {
	h := handler{}

	e.GET("/v1/healthz", h.Healthz)
}

func (h handler) Healthz(c echo.Context) error {
	return httpx.WriteResponse(c, http.StatusOK, "All system running well :)")
}
