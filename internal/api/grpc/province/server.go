package province

import (
	pb "shop/internal/pb/province"
	service "shop/internal/service/province"
)

type Server struct {
	pb.UnimplementedProvinceServiceServer
	service service.Service
}

func New(service service.Service) Server {
	return Server{
		service: service,
	}
}