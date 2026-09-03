package category

import (
	"context"
	dto "shop/internal/dto/category"
	pb "shop/internal/pb/category"
	"shop/internal/pkg/mapper"
)

func (s Server) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.CategoryResponse, error) {
	const op = "category-grpc.Update"

	var image *dto.ImageFile
	if len(req.GetImage()) > 0 {
		image = &dto.ImageFile{Content: req.GetImage()}
	}

	dtoReq := dto.UpdateRequest{
		Title: req.Title,
		Slug:  req.NewSlug,
		Image: image,
	}

	if validationErr := s.validator.Update(ctx, dtoReq); validationErr != nil {
		return nil, mapper.ErrorToGrpc(validationErr)
	}

	res, err := s.service.Update(ctx, req.GetSlug(), dtoReq)
	if err != nil {
		return nil, mapper.ErrorToGrpc(err)
	}

	var parentID *uint64
	if res.ParentID != nil {
		v := uint64(*res.ParentID)
		parentID = &v
	}

	return &pb.CategoryResponse{
		Id:       uint64(res.ID),
		Title:    res.Title,
		Slug:     res.Slug,
		ParentId: parentID,
		Image:    res.Image,
	}, nil
}
