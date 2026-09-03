package product

import (
	dto "shop/internal/dto/product"
	pb "shop/internal/pb/product"
)

func toProtoProductResponse(res dto.CreateResponse) *pb.ProductResponse {
	return &pb.ProductResponse{
		Id:          uint64(res.ID),
		Name:        res.Name,
		Slug:        res.Slug,
		Description: res.Description,
		Price:       uint64(res.Price),
		Stock:       uint64(res.Stock),
		CategoryId:  uint64(res.CategoryID),
		MainImage:   res.MainImage,
		Images:      res.Images,
	}
}
