package order

import (
	"net/http"
	"skripsi-be/lib/httpx"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/constant"
	"skripsi-be/type/params"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h handler) CreateNewOrder(c echo.Context) error {
	var req []params.ServiceIngestionOrder
	err := c.Bind(&req)
	if err != nil {
		return httpx.WriteErrorResponse(c, errors.ErrBadRequest, "bad request")
	}

	jwtClaims, err := httpx.GetJWTClaimsFromContext(c, constant.AuthContextKey)
	if err != nil {
		return httpx.WriteErrorResponse(c, errors.ErrInternalServer, "cannot get jwt claims from context")
	}

	jwtSub, err := jwtClaims.GetSubject()
	if err != nil {
		return httpx.WriteErrorResponse(c, errors.ErrInternalServer, "cannot get jwt sub from context")
	}

	for idx := range req {
		req[idx].CashierID, err = uuid.Parse(jwtSub)
		if err != nil {
			return httpx.WriteErrorResponse(c, errors.ErrInternalServer, "cannot parse jwt sub")
		}
	}

	res, err := h.service.IngestOrder(c.Request().Context(), req)
	if err != nil {
		return httpx.WriteErrorResponse(c, err, nil)
	}

	return httpx.WriteResponse(c, http.StatusCreated, res)
}

func (h handler) CreateNewOrderWithoutKappa(c echo.Context) error {
	return httpx.WriteErrorResponse(c, errors.ErrInternalServer, "Internal Server Error")
}
