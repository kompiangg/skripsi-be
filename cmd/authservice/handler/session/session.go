package session

import (
	"net/http"
	"skripsi-be/lib/httpx"
	"skripsi-be/type/params"

	"github.com/labstack/echo/v4"
)

func (h handler) NewSession(c echo.Context) error {
	var req params.ServiceLogin
	err := c.Bind(&req)
	if err != nil {
		return httpx.WriteErrorResponse(c, err, nil)
	}

	res, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		return httpx.WriteErrorResponse(c, err, nil)
	}

	return httpx.WriteResponse(c, http.StatusOK, res)
}
