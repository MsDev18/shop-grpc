package category

import (
	"context"
	dto "shop/internal/dto/category"
	"shop/internal/pkg/richerror"
)

func (s Service) Update(ctx context.Context, slug string, req dto.UpdateRequest) (dto.CreateResponse, error) {
	const op = "category-service.Update"

	// get category by slug
	category, err := s.repository.GetOneBySlug(ctx, slug)
	if err != nil {
		return dto.CreateResponse{}, err
	}
	// check uniqueness req.slug
	if req.Slug != nil && *req.Slug != category.Slug {
		isExistsSlug, err := s.repository.IsExistsSlug(ctx, *req.Slug)
		if err != nil {
			return dto.CreateResponse{}, err
		}
		if isExistsSlug {
			return dto.CreateResponse{}, richerror.New().
				SetOp(op).
				SetMsg("this slug already exists").
				SetKind(richerror.KindConflictErr)
		}
	}
	// child category have not image
	if category.ParentID != nil && req.Image != nil {
		return dto.CreateResponse{}, richerror.New().
			SetOp(op).
			SetMsg("child category is not allowed to upload image").
			SetKind(richerror.KindBadRequestErr)
	}
	// upload image & call imageProcessor
	var imageURI *string
	if req.Image != nil {
		url, err := s.imageProcessor.Process(ctx, req.Image)
		if err != nil {
			return dto.CreateResponse{}, err
		}
		imageURI = &url
	}
	// call update method from repository
	err = s.repository.UpdateByID(ctx, category.ID, req.Title, req.Slug, imageURI)
	if err != nil {
		return dto.CreateResponse{}, err
	}

	updated, err := s.repository.GetOneByID(ctx, category.ID)
	if err != nil {
		return dto.CreateResponse{}, err
	}

	return dto.CreateResponse{
		ID:       updated.ID,
		Title:    updated.Title,
		Slug:     updated.Slug,
		ParentID: updated.ParentID,
		Image:    updated.Image,
	}, nil
}
