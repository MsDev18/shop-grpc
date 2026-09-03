package category

import (
	"context"
	pb "shop/internal/pb/category"
	"shop/internal/pkg/mapper"
)

func (s Server) Delete(ctx context.Context ,req *pb.DeleteRequest) (*pb.DeleteResponse, error ) {
	const op = "category-grpc.Delete"

	if err := s.service.Delete(ctx, req.GetSlug()); err != nil {
		return nil, mapper.ErrorToGrpc(err)
	}
	
	return &pb.DeleteResponse{}, nil
}