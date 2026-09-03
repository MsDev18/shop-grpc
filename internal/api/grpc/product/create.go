package product

import (
	"io"
	dto "shop/internal/dto/product"
	pb "shop/internal/pb/product"
	"shop/internal/pkg/mapper"
	"shop/internal/pkg/richerror"
)

func (s Server) Create(stream pb.ProductService_CreateServer) error {
	const op = "product-grpc.Create"

	ctx := stream.Context()

	var metadata *pb.ProductMetadata
	var galleryImages []*dto.ImageFile

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			return mapper.ErrorToGrpc(
				richerror.New().
					SetOp(op).
					SetMsg("error reciving image stream").
					SetKind(richerror.KindUnauthorizeErr).
					SetErr(err),
			)
		}

		switch payload := req.Payload.(type) {
		case *pb.CreateRequest_Metadata:
			if metadata != nil {
				return mapper.ErrorToGrpc(
					richerror.New().
						SetOp(op).
						SetMsg("meta data send more than once").
						SetKind(richerror.KindBadRequestErr),
				)
			}
			metadata = payload.Metadata
		case *pb.CreateRequest_Image:
			galleryImages = append(galleryImages, &dto.ImageFile{Content: payload.Image.GetImage()})
		}
	}

	if metadata == nil {
		return mapper.ErrorToGrpc(
			richerror.New().
				SetOp(op).
				SetMsg("product metadata was not send").
				SetKind(richerror.KindBadRequestErr),
		)
	}

	var stock *uint
	if metadata.Stock != nil {
		v := uint(metadata.GetStock())
		stock = &v
	}

	dtoReq := dto.CreateRequest{
		Name:        metadata.GetName(),
		Slug:        metadata.GetSlug(),
		Description: metadata.GetDescription(),
		Price:       uint(metadata.GetPrice()),
		Stock:       stock,
		CategoryID:  uint(metadata.GetCategoryId()),
		MainImage:   &dto.ImageFile{Content: metadata.GetMainImage()},
		Images:      galleryImages,
	}

	if validationErr := s.validator.Create(ctx, dtoReq); validationErr != nil {
		return mapper.ErrorToGrpc(validationErr)
	}

	res, serviceErr := s.service.Create(ctx, dtoReq)
	if serviceErr != nil {
		return mapper.ErrorToGrpc(serviceErr)
	}

	return stream.SendAndClose(toProtoProductResponse(res))
}
