package category

import (
	"context"
	pb "shop/internal/pb/category"
	"shop/internal/pkg/mapper"
	dto "shop/internal/dto/category"
)

func (s Server) GetAll (ctx context.Context , req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	const op = "category-grpc.GetAll"

	categories, err := s.service.GetAll(ctx)
	if err != nil {
		return nil , mapper.ErrorToGrpc(err)
	}

	resp := &pb.GetAllResponse{
		Categories: make([]*pb.CategoryResponse,0, len(categories)),
	}

	for _, c := range categories {
		resp.Categories = append(resp.Categories, toProtoCategoryTree(c))
	}

	return resp, nil
}

func toProtoCategoryTree(c dto.CategoryResponse) *pb.CategoryResponse {
	pc := &pb.CategoryResponse{
		Id:    uint64(c.ID),
		Title: c.Title,
		Slug:  c.Slug,
		Image: c.Image,
	}
	for _, child := range c.Children {
		pc.Children = append(pc.Children, toProtoCategoryTree(child))
	}
	return pc
}