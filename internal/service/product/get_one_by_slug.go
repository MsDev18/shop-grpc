package product

import (
	"context"
	dto "shop/internal/dto/product"
)

func (s Service) GetOneBySlug(ctx context.Context, slug string) (dto.CreateResponse, error) {
	const op = "product-service.GetOneBySlug"

	p, pImage, err := s.repository.GetOneBySlug(ctx, slug)
	if err != nil {
		return dto.CreateResponse{}, err
	}
	// map pImage to []string
	var productImage = make([]string, len(pImage))
	for i, v := range pImage {
		productImage[i] = v.Image
	}
	// map to dto.CreateResponse
	response := dto.CreateResponse{
		ID:          p.ID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CategoryID:  p.CategoryID,
		MainImage:   p.MainImage,
		Images:      productImage,
	}

	return response, nil
}
