package product

import (
	"context"
	dto "shop/internal/dto/product"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func (s Service) Create(ctx context.Context, req dto.CreateRequest) (dto.CreateResponse, error) {
	const op = "product-service-Create"
	// check category
	_, err := s.categoryService.GetOneByID(ctx, req.CategoryID)
	if err != nil {
		return dto.CreateResponse{}, err
	}
	// check exists slug
	exists, err := s.repository.IsExistsSlug(ctx, req.Slug)
	if err != nil {
		return dto.CreateResponse{}, err
	}
	if exists {
		return dto.CreateResponse{}, richerror.New().
			SetOp(op).
			SetMsg("conflict slug , this slug already exists").
			SetKind(richerror.KindConflictErr)
	}

	// process main image
	mainImage, err := s.imageProcessor.Process(ctx, req.MainImage)
	if err != nil {
		return dto.CreateResponse{}, err
	}

	// process gallery images
	imagePaths := make([]string, 0, len(req.Images))
	for _, image := range req.Images {
		path, err := s.imageProcessor.Process(ctx, image)
		if err != nil {
			return dto.CreateResponse{}, err
		}
		imagePaths = append(imagePaths, path)
	}

	// resolve optional stock (nil -> 0)
	var stock uint
	if req.Stock != nil {
		stock = *req.Stock
	}

	// build entity
	p := entity.Product{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Price:       req.Price,
		Stock:       stock,
		MainImage:   mainImage,
		CategoryID:  req.CategoryID,
	}

	// insert product + product_image atomically
	created, err := s.repository.Create(ctx, p, imagePaths)
	if err != nil {
		return dto.CreateResponse{}, err
	}

	// map to response
	response := dto.CreateResponse{
		ID:          created.ID,
		Name:        created.Name,
		Slug:        created.Slug,
		Description: created.Description,
		Price:       created.Price,
		Stock:       created.Stock,
		CategoryID:  created.CategoryID,
		MainImage:   created.MainImage,
		Images:      imagePaths,
	}
	return response, nil
}
