package auth

import (
	authpb "shop/internal/pb/auth"
	service "shop/internal/service/auth"
	validator "shop/internal/validator/auth"
)

type Server struct {
	authpb.UnimplementedAuthServiceServer
	validator validator.Validator
	service   service.Service
}

func New(service service.Service, validator validator.Validator) Server {
	return Server{
		validator: validator,
		service:   service,
	}
}
