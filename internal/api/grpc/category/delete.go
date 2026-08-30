package category

import (
	"context"
	"shop/internal/api/grpc/auth"
	"shop/internal/entity"
	pb "shop/internal/pb/category"
	"shop/internal/pkg/mapper"
)

func (s Server) Delete(ctx context.Context ,req *pb.DeleteRequest) (*pb.DeleteResponse, error ) {
	const op = "category-grpc.Selete"

	if err := auth.RequireRole(ctx, entity.AdminRole); err != nil {
		return nil, mapper.ErrorToGrpc(err)
	}

	if err := s.service.Delete(ctx, req.GetSlug()); err != nil {
		return nil, mapper.ErrorToGrpc(err)
	}
	
	return &pb.DeleteResponse{}, nil
}