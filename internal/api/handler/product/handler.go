package product

import (
	service "shop/internal/service/product"
	validator "shop/internal/validator/product"
)

type Handler struct {
	service   service.Service
	validator validator.Validator
}

func New(service service.Service , validator validator.Validator) Handler {
	return Handler{
		service: service,
		validator: validator,
	}
}
