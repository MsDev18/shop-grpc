package province

import (
	"context"
	pb "shop/internal/pb/province"
	"shop/internal/pkg/mapper"
)

func (s Server) GetAll(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	const op = "province-grpc.GetAll"

	provinces, err := s.service.GetAll(ctx)
	if err != nil {
		return nil, mapper.ErrorToGrpc(err)
	}

	resp := &pb.GetAllResponse{
		Provinces: make([]*pb.ProvinceResponse, 0., len(provinces)),
	}

	for _, p := range provinces {
		resp.Provinces = append(resp.Provinces, &pb.ProvinceResponse{
			Id:   uint64(p.ID),
			Name: p.Name,
		})
	}

	return resp , nil
}
