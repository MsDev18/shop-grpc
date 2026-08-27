package category

import (
	"context"
	"shop/internal/entity"
	"shop/internal/pkg/imageprocessor"
)

type Service struct {
	repository     Repository
	imageProcessor imageprocessor.Processor
}

type Repository interface {
	IsExistsSlug(ctx context.Context, slug string) (bool, error)
	Create(ctx context.Context, category entity.Category) (entity.Category, error)
	GetOneByID(ctx context.Context, id uint) (entity.Category, error)
	GetAll(ctx context.Context) ([]entity.Category, error)
	GetOneBySlug(ctx context.Context, slug string) (entity.Category, error)
	GetChildrenByParentID(ctx context.Context, parentID uint) ([]entity.Category, error)
	UpdateByID(ctx context.Context, id uint, title, slug, image *string) error
	DeleteByID(ctx context.Context, id uint) error
}

func New(repository Repository, imageprocessor imageprocessor.Processor) Service {
	return Service{
		repository:     repository,
		imageProcessor: imageprocessor,
	}
}
