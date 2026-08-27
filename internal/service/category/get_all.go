package category

import (
	"context"
	dto "shop/internal/dto/category"
	"shop/internal/entity"
)

func (s Service) GetAll(ctx context.Context) ([]dto.CategoryResponse, error) {
	const op = "category-service.GetAll"
	// call GetAll method from category repository
	categoies, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return s.buildTree(categoies), nil
}

func (s Service) buildTree(categories []entity.Category) ([]dto.CategoryResponse) {
	var childrentByParent = make(map[uint][]entity.Category)
	var root []entity.Category

	for _, c := range categories {
		if c.ParentID == nil {
			root = append(root, c)
			continue
		}
		childrentByParent[*c.ParentID]= append(childrentByParent[*c.ParentID], c)
	}

	// call method for map to []dto.CategoryResponse
	return s.mapToCategoryResponse(root, childrentByParent)

}


func (s Service) mapToCategoryResponse (categories []entity.Category, childrenByParent map[uint][]entity.Category) []dto.CategoryResponse {
	result := make([]dto.CategoryResponse , 0 , len(categories))
	for _, c := range categories {
		result =  append(result, dto.CategoryResponse{
			ID:       c.ID,
			Title:    c.Title,
			Slug:     c.Slug,
			Image:    c.Image,
			Children: s.mapToCategoryResponse(childrenByParent[c.ID] , childrenByParent),
		})
	}
	return result
}