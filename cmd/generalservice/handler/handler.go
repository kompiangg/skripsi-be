package handler

import (
	"skripsi-be/cmd/generalservice/handler/cashier"
	"skripsi-be/cmd/generalservice/handler/customer"
	"skripsi-be/cmd/generalservice/handler/healthz"
	"skripsi-be/cmd/generalservice/handler/item"
	"skripsi-be/cmd/generalservice/handler/payment_types"
	"skripsi-be/cmd/middleware"
	"skripsi-be/service"

	"github.com/labstack/echo/v4"
)

func Init(
	echo *echo.Echo,
	service service.Service,
	middleware middleware.Middleware,
) {
	healthz.Init(echo)
	item.Init(echo, service.Item)
	payment_types.Init(echo, service.PaymentTypes)
	customer.Init(echo, service.Customer)
	cashier.Init(echo, middleware, service.Cashier)
}
