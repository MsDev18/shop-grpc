package product

import (
	"context"
	"errors"
	dto "shop/internal/dto/product"
	"shop/internal/pkg/richerror"
)

const (
	DEFAULT_PAGE  = 1
	DEFAULT_LIMIT = 20
	MAX_LIMIT     = 50
)

func (s Service) GetAll(ctx context.Context, req dto.GetAllRequest) (dto.GetAllResponse, error) {
	const op = "product-service.GetAll"

	page := req.Page
	if page < 1 {
		page = DEFAULT_PAGE
	}

	limit := req.Limit
	if limit < 1 {
		limit = DEFAULT_LIMIT
	}
	if limit > MAX_LIMIT {
		limit = MAX_LIMIT
	}

	offset := (page - 1) * limit

	var categoryID *uint
	if req.CategorySlug != nil {
		category, err := s.categoryService.GetOne(ctx, *req.CategorySlug)
		if err != nil {
			var richErr *richerror.RichError
			if errors.As(err, &richErr) && richErr.GetKind() == richerror.KindNotFoundErr {
				// دسته‌بندی‌ای با این slug وجود نداره -> یعنی هیچ محصولی نمی‌تونه match بشه
				return dto.GetAllResponse{
					Products: []dto.ProductListItem{},
					Meta: dto.PaginationMeta{
						Page:       page,
						Limit:      limit,
						Total:      0,
						TotalPages: 0,
					},
				}, nil
			}
			return dto.GetAllResponse{}, err
		}
		categoryID = &category.ID
	}

	products, total, err := s.repository.GetAll(ctx, limit, offset, categoryID)
	if err != nil {
		return dto.GetAllResponse{}, err
	}

	items := make([]dto.ProductListItem, 0, len(products))
	for _, p := range products {
		items = append(items, dto.ProductListItem{
			ID:         p.ID,
			Name:       p.Name,
			Slug:       p.Slug,
			Price:      p.Price,
			Stock:      p.Stock,
			CategoryID: p.CategoryID,
			MainImage:  p.MainImage,
		})
	}

	totalPages := (total + limit - 1) / limit

	response := dto.GetAllResponse{
		Products: items,
		Meta: dto.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return response, nil
}
