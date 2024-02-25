package orders

import (
	"net/http"
	"skripsi-be/lib/httpx"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/constant"
	"skripsi-be/type/params"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h handler) GetBriefInformation(c echo.Context) error {
	var req params.FindOrderService
	err := c.Bind(&req)
	if err != nil {
		return httpx.WriteErrorResponse(c, errors.Wrap(err), nil)
	}

	jwtClaims, err := httpx.GetJWTClaimsFromContext(c, constant.AuthContextKey)
	if err != nil {
		return httpx.WriteErrorResponse(c, errors.Wrap(err), nil)
	}

	if jwtClaims.Role == constant.RoleEnum.Cashier {
		cashierID, err := uuid.Parse(jwtClaims.Subject)
		if err != nil {
			return httpx.WriteErrorResponse(c, errors.Wrap(err), nil)
		}

		req.CashierID = uuid.NullUUID{
			UUID:  cashierID,
			Valid: true,
		}
	}

	allOrders, err := h.orderService.FindBriefInformationOrder(c.Request().Context(), req)
	if err != nil {
		return httpx.WriteErrorResponse(c, errors.Wrap(err), nil)
	}

	return httpx.WriteResponse(c, http.StatusOK, allOrders)
}
