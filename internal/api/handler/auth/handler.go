package auth

import (
	authservice "shop/internal/service/auth"
	authvalidator "shop/internal/validator/auth"
)

type Handler struct {
	service    authservice.Service
	validator  authvalidator.Validator
}


func New ( service authservice.Service, validator authvalidator.Validator) Handler {
	return Handler{
		service:    service,
		validator:  validator,
	}
}