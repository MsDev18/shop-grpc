package category

import (
	"context"
	dto "shop/internal/dto/category"
	"shop/internal/entity"
)

func (s Service) GetOne(ctx context.Context, slug string) (dto.CategoryResponse, error) {
	const op = "category-service.GetOne"
	// call repository
	c, err := s.repository.GetOneBySlug(ctx, slug)
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	childrenByParent := make(map[uint][]entity.Category)

	if c.ParentID == nil {
		childrens, err := s.repository.GetChildrenByParentID(ctx, c.ID)
		if err != nil {
			return dto.CategoryResponse{}, err
		}
		childrenByParent[c.ID] = childrens
	}

	response := s.mapToCategoryResponse([]entity.Category{c}, childrenByParent)
	return response[0], nil
}
