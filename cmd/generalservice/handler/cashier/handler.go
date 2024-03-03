package cashier

import (
	"skripsi-be/cmd/middleware"
	"skripsi-be/service/cashier"

	"github.com/labstack/echo/v4"
)

type handler struct {
	cashierService cashier.Service
}

func Init(e *echo.Echo, mw middleware.Middleware, cashierService cashier.Service) {
	h := handler{
		cashierService: cashierService,
	}

	e.GET("/v1/cashiers/:id", h.FindCashierByID, mw.JWTRestricted())
}
