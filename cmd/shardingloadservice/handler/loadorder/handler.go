package loadorder

import "skripsi-be/service/order"

type handler struct {
	orderService order.Service
}

func New(
	orderService order.Service,
) handler {
	return handler{
		orderService: orderService,
	}
}
