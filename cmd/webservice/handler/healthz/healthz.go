package healthz

import (
	"net/http"
	"skripsi-be/lib/httpx"

	"github.com/labstack/echo/v4"
)

func (h handler) Healthz(c echo.Context) error {
	return httpx.WriteResponse(c, http.StatusOK, "All system running well :)")
}
