package product

import (
	"context"
	dto "shop/internal/dto/product"
	pb "shop/internal/pb/product"
	"shop/internal/pkg/mapper"
)

func (s Server) GetAll(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	const op = "product-grpc.GetAll"

	var categorySlug *string
	if req.CategorySlug != nil {
		v := req.GetCategorySlug()
		categorySlug = &v
	}

	dtoReq := dto.GetAllRequest{
		Page:         int(req.GetPage()),
		Limit:        int(req.GetLimit()),
		CategorySlug: categorySlug,
	}

	res, err := s.service.GetAll(ctx, dtoReq)
	if err != nil {
		return nil, mapper.ErrorToGrpc(err)
	}

	items := make([]*pb.ProductListItem, 0, len(res.Products))
	for _, p := range res.Products {
		items = append(items, &pb.ProductListItem{
			Id:         uint64(p.ID),
			Name:       p.Name,
			Slug:       p.Slug,
			Price:      uint64(p.Price),
			Stock:      uint64(p.Stock),
			CategoryId: uint64(p.CategoryID),
			MainImage:  p.MainImage,
		})
	}

	return &pb.GetAllResponse{
		Products: items,
		Meta: &pb.PaginationMeta{
			Page:       int32(res.Meta.Page),
			Limit:      int32(res.Meta.Limit),
			Total:      int32(res.Meta.Total),
			TotalPages: int32(res.Meta.TotalPages),
		},
	}, nil
}
