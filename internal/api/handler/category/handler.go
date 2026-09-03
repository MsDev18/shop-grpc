package category

import (
	service "shop/internal/service/category"
	validator "shop/internal/validator/category"

)

type Handler struct {
	service service.Service
	validator validator.Validator
}

func New (service service.Service, validator validator.Validator) Handler {
	return Handler{
		service: service,
		validator: validator,
	}
}