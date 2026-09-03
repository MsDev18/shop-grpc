package product

import (
	pb "shop/internal/pb/product"
	service "shop/internal/service/product"
	validator "shop/internal/validator/product"
)

type Server struct {
	pb.UnimplementedProductServiceServer
	service service.Service
	validator validator.Validator
}

func New(service service.Service , validator validator.Validator) Server {
	return Server{
		service:   service,
		validator: validator,
	}
}