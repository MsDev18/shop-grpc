package user

import (
	pb "shop/internal/pb/user"
	service "shop/internal/service/user"
	validator "shop/internal/validator/user"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	service   service.Service
	validator validator.Validator
}

func New(service service.Service, validator validator.Validator) Server {
	return Server{
		service:   service,
		validator: validator,
	}
}
