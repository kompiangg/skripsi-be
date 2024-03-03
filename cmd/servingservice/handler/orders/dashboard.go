package orders

import (
	"net/http"
	"skripsi-be/lib/httpx"
	"skripsi-be/pkg/errors"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h handler) GetDashboard(c echo.Context) error {
	interval := c.QueryParam("interval")
	if interval == "" {
		return httpx.WriteErrorResponse(c, errors.ErrBadRequest, nil)
	}

	offsetStr := c.QueryParam("offset")
	if offsetStr == "" {
		return httpx.WriteErrorResponse(c, errors.ErrBadRequest, nil)
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return httpx.WriteErrorResponse(c, errors.Wrap(err), nil)
	}

	result, err := h.orderService.FindInsightBasedOnInterval(c.Request().Context(), interval, offset)
	if err != nil {
		return httpx.WriteErrorResponse(c, errors.Wrap(err), nil)
	}

	return httpx.WriteResponse(c, http.StatusOK, result)
}
