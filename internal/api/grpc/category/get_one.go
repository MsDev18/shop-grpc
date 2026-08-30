package category

import (
	"context"
	pb "shop/internal/pb/category"
	"shop/internal/pkg/mapper"
)

func (s Server) GetOne(ctx context.Context, req *pb.GetOneRequest) (*pb.CategoryResponse, error) {
	c, err := s.service.GetOne(ctx, req.GetSlug())
	if err != nil {
		return nil, mapper.ErrorToGrpc(err)
	}
	return toProtoCategoryTree(c), nil
}