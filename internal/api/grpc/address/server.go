package address

import (
	pb "shop/internal/pb/address"
	service "shop/internal/service/address"
	validator "shop/internal/validator/address"
)

type Server struct {
	pb.UnimplementedAddressServiceServer
	service   service.Service
	validator validator.Validator
}

func New(service service.Service, validator validator.Validator) Server {
	return Server{
		service:   service,
		validator: validator,
	}
}
