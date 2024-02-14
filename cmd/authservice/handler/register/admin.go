package register

import (
	"net/http"
	"skripsi-be/lib/httpx"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/constant"
	"skripsi-be/type/params"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h handler) Admin(c echo.Context) error {
	var req params.ServiceRegisterAdmin
	err := c.Bind(&req)
	if err != nil {
		return httpx.WriteErrorResponse(c, errors.New(err), nil)
	}

	jwtClaims, err := httpx.GetJWTClaimsFromContext(c, constant.AuthContextKey)
	if err != nil {
		return httpx.WriteErrorResponse(c, errors.New(err), nil)
	}

	requestBy, err := jwtClaims.GetSubject()
	if err != nil {
		return httpx.WriteErrorResponse(c, err, nil)
	}

	req.RequestBy, err = uuid.Parse(requestBy)
	if err != nil {
		err = errors.New(err)
		return httpx.WriteErrorResponse(c, err, nil)
	}

	err = h.service.RegisterAdmin(c.Request().Context(), req)
	if err != nil {
		return httpx.WriteErrorResponse(c, err, nil)
	}

	return httpx.WriteResponse(c, http.StatusOK, nil)
}
