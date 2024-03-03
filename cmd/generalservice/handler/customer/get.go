package customer

import (
	"net/http"
	"skripsi-be/lib/httpx"

	"github.com/labstack/echo/v4"
)

func (h handler) FindLikeOneOfAllColumn(c echo.Context) error {
	req := c.QueryParam("q")

	customers, err := h.customerService.FindLikeOneOfAllColumn(c.Request().Context(), req)
	if err != nil {
		return httpx.WriteErrorResponse(c, err, nil)
	}

	return httpx.WriteResponse(c, http.StatusOK, customers)
}
