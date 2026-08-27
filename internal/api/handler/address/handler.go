package address

import (
	service "shop/internal/service/address"
	validator "shop/internal/validator/address"
)

type Handler struct {
	service   service.Service
	validator validator.Validator
}

func New (service service.Service,validator validator.Validator) Handler {
	return Handler{
		service: service,
		validator: validator,
	}
}