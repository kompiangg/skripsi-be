package orders

import (
	"net/http"
	"skripsi-be/lib/httpx"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/params"

	"github.com/labstack/echo/v4"
)

func (h handler) GetAllDetails(c echo.Context) error {
	var req params.FindOrderService
	err := c.Bind(&req)
	if err != nil {
		return httpx.WriteErrorResponse(c, errors.Wrap(err), nil)
	}

	allOrders, err := h.service.FindOrder(c.Request().Context(), req)
	if err != nil {
		return httpx.WriteErrorResponse(c, errors.Wrap(err), nil)
	}

	return httpx.WriteResponse(c, http.StatusOK, allOrders)
}
