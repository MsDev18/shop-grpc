package category

import (
	"context"
	dto "shop/internal/dto/category"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (s Service) Create(ctx context.Context, req dto.CreateRequest) (dto.CreateResponse, error) {
	const op = "category-service.Create"
	// 1. check parnet
	if req.ParentID != nil {
		parent, err := s.repository.GetOneByID(ctx, *req.ParentID)
		if err != nil {
			return dto.CreateResponse{}, err
		}
		if parent.ParentID != nil {
			return dto.CreateResponse{}, richerror.New().
				SetOp(op).
				SetMsg("category hierarchy can't exceed two levels").
				SetKind(richerror.KindBadRequestErr)
		}
	}
	// 2. check uniqueness slug
	isExists, err := s.repository.IsExistsSlug(ctx, req.Slug)
	if err != nil {
		return dto.CreateResponse{}, err
	}
	if isExists {
		return dto.CreateResponse{}, richerror.New().
			SetOp(op).
			SetMsg("this slug already exists").
			SetKind(richerror.KindConflictErr)
	}

	// 3. upload image
	var imageURL *string
	if req.Image != nil {
		image, err := s.imageProcessor.Process(ctx, req.Image)
		if err != nil {
			return dto.CreateResponse{}, err
		}
		imageURL = &image
	}

	// 4. generate new record
	categoryEntity := entity.Category{
		ParentID: req.ParentID,
		Title:    req.Title,
		Slug:     req.Slug,
		Image:    imageURL,
	}

	// 5. insert into repository
	category, err := s.repository.Create(ctx, categoryEntity)
	if err != nil {
		return dto.CreateResponse{}, err
	}

	// 6. map to dto.CreateResponse
	response := dto.CreateResponse{
		ID:       category.ID,
		Title:    category.Title,
		Slug:     category.Slug,
		ParentID: category.ParentID,
		Image:    category.Image,
	}
	// 7. return response
	return response, nil
}
