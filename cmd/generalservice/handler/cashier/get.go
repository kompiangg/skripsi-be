package cashier

import (
	"net/http"
	"skripsi-be/lib/httpx"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/constant"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h handler) FindCashierByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpx.WriteErrorResponse(c, errors.ErrBadRequest, "invalid id")
	}

	jwtClaims, err := httpx.GetJWTClaimsFromContext(c, constant.AuthContextKey)
	if err != nil {
		return httpx.WriteErrorResponse(c, errors.ErrInternalServer, "cannot get jwt claims from context")
	}

	if jwtClaims.Role == constant.RoleEnum.Cashier {
		id, err = uuid.Parse(jwtClaims.Subject)
		if err != nil {
			return httpx.WriteErrorResponse(c, errors.ErrInternalServer, "cannot parse jwt sub")
		}
	}

	res, err := h.cashierService.FindCashierByID(c.Request().Context(), id)
	if err != nil {
		return httpx.WriteErrorResponse(c, err, nil)
	}

	return httpx.WriteResponse(c, http.StatusOK, res)
}
