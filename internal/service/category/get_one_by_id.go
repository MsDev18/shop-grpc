package category

import (
	"context"
	dto "shop/internal/dto/category"
)

func (s Service) GetOneByID(ctx context.Context, id uint) (dto.CategoryResponse, error) {
	const op = "category-service.GetOneByID"

	c, err := s.repository.GetOneByID(ctx, id)
	if err != nil {
		return dto.CategoryResponse{}, err
	}


	// map to dto.CreateResponse
	category := dto.CategoryResponse{
		ID:       c.ID,
		Title:    c.Title,
		Slug:     c.Slug,
		Image:    c.Image,
	}

	return category, nil
}