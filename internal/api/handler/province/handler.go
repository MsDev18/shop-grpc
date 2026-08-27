package province

import (
	service "shop/internal/service/province"
)

type Handler struct {
	service service.Service
}

func New (service service.Service) Handler {
	return Handler{
		service: service,
	}
}