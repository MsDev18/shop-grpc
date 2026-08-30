package province

import (
	"context"
	pb "shop/internal/pb/province"
	"shop/internal/pkg/mapper"
)

func (s Server) GetOne(ctx context.Context, req *pb.GetOneRequest) (*pb.ProvinceResponse, error) {
	const op = "province-grpc.Getone"

	p, err := s.service.GetOne(ctx, uint(req.GetId()))
	if err != nil {
		return nil, mapper.ErrorToGrpc(err)
	}

	return &pb.ProvinceResponse{
		Id:   uint64(p.ID),
		Name: p.Name,
	}, err
}
