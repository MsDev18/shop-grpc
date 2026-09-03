package product

import (
	"context"
	pb "shop/internal/pb/product"
	"shop/internal/pkg/mapper"
)

func (s Server) GetOneBySlug(ctx context.Context, req *pb.GetOneBySlugRequest) (*pb.ProductResponse, error) {
	const op = "product-grpc.GetOneBySlug"

	res, err := s.service.GetOneBySlug(ctx, req.GetSlug())
	if err != nil {
		return nil, mapper.ErrorToGrpc(err)
	}

	return toProtoProductResponse(res), nil
}
