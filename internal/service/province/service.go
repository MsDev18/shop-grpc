package province

import (
	"context"
	"shop/internal/entity"
)

type Service struct {
	repository Repository
}

type Repository interface {
	GetAll(ctx context.Context) ([]entity.Province, error)
	GetOneByID(ctx context.Context, id uint) (entity.Province, error)
}

func New(repository Repository) Service {
	return Service{
		repository: repository,
	}
}
