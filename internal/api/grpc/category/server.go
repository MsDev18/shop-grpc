package category

import (
	pb "shop/internal/pb/category"
	service "shop/internal/service/category"
	validator "shop/internal/validator/category"
)

type Server struct {
	pb.UnimplementedCategoryServiceServer
	service   service.Service
	validator validator.Validator
}

func New(service service.Service, validator validator.Validator) Server {
	return Server{
		service:   service,
		validator: validator,
	}
}
