package address

import (
	"context"
	"shop/internal/entity"
	"shop/internal/service/province"
)

type Service struct {
	repository      Repository
	provinceService province.Service
}

type Repository interface {
	Create(ctx context.Context, address entity.Address) (entity.Address, error)
	GetAll(ctx context.Context, userID uint) ([]entity.Address, error)
	GetOne(ctx context.Context, userID uint, addressID uint) (entity.Address, error)
	Delete(ctx context.Context, userID uint, addressID uint) error
	Update(ctx context.Context, userID uint, addressID uint, title *string, provinceID *uint, city, address, postalCode *string) error
}

func New(repository Repository, provinceService province.Service) Service {
	return Service{
		repository:      repository,
		provinceService: provinceService,
	}
}
