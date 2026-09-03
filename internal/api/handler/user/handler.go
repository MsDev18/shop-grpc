package user

import (
	userservice "shop/internal/service/user"
	uservalidator "shop/internal/validator/user"
)

type Handler struct {
	service userservice.Service
	validator uservalidator.Validator
}

func New ( service userservice.Service , validator uservalidator.Validator) Handler {
	return Handler{
		service: service,
		validator: validator,
	}
}