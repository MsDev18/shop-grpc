package category

import (
	"context"
	pb "shop/internal/pb/category"
	"shop/internal/pkg/mapper"
	dto "shop/internal/dto/category"
)

func (s Server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CategoryResponse, error) {
	const op = "category-grpc.Create"

	var parentID *uint
	if req.ParentId != nil {
		v := uint(req.GetParentId())
		parentID = &v
	}

	var image *dto.ImageFile
	if len(req.GetImage()) > 0 {
		image = &dto.ImageFile{Content: req.GetImage()}
	}

	dtoReq := dto.CreateRequest{
		ParentID: parentID,
		Title:    req.GetTitle(),
		Slug:     req.GetSlug(),
		Image:    image,
	}

	if validationErr := s.validator.Create(ctx,dtoReq) ; validationErr != nil {
		return nil , mapper.ErrorToGrpc(validationErr)
	}

	res, serviceErr := s.service.Create(ctx , dtoReq)
	if serviceErr != nil {
		return nil , mapper.ErrorToGrpc(serviceErr)
	}

	var respParentID *uint64
	if res.ParentID != nil {
		v := uint64(*res.ParentID)
		respParentID = &v
	}

	return &pb.CategoryResponse{
		Id:       uint64(res.ID),
		Title:    res.Title,
		Slug:     res.Slug,
		ParentId: respParentID,
		Image:    res.Image,
	}, nil
}