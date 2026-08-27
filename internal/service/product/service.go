package product

import (
	"context"
	"shop/internal/entity"
	"shop/internal/pkg/imageprocessor"
	"shop/internal/service/category"
)

type Service struct {
	repository      Repository
	imageProcessor  imageprocessor.Processor
	categoryService category.Service
}

type Repository interface {
	Create(ctx context.Context, product entity.Product, imagePaths []string) (entity.Product, error)
	GetOneBySlug(ctx context.Context, slug string) (entity.Product, []entity.ProductImage, error)
	GetProductImage(ctx context.Context, productID uint) ([]entity.ProductImage, error)
	IsExistsSlug(ctx context.Context, slug string) (bool, error)
	GetAll(ctx context.Context, limit int, offset int, categoryID *uint) ([]entity.Product, int, error)
}

func New(repository Repository, categoryService category.Service, imageProcessor imageprocessor.Processor) Service {
	return Service{
		repository:      repository,
		categoryService: categoryService,
		imageProcessor:  imageProcessor,
	}
}
