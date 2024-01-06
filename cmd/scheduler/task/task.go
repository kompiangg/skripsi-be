package task

import "skripsi-be/service"

type task struct {
	service service.Service
}

func New(
	service service.Service,
) task {
	return task{
		service: service,
	}
}
