package item

import (
	"net/http"
	"skripsi-be/lib/httpx"

	"github.com/labstack/echo/v4"
)

func (h handler) FindLikeNameOrID(c echo.Context) error {
	nameOrID := c.QueryParam("nameOrID")

	items, err := h.itemService.FindLikeNameOrID(c.Request().Context(), nameOrID)
	if err != nil {
		return httpx.WriteErrorResponse(c, err, nil)
	}

	return httpx.WriteResponse(c, http.StatusOK, items)
}
