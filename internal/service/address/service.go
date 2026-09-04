package address

import (
	"context"
	provincedto "shop/internal/dto/province"
	"shop/internal/entity"
)

type ProvinceService interface {
	GetOne(ctx context.Context, id uint) (provincedto.GetOneResponse, error)
	GetAll(ctx context.Context) ([]provincedto.GetOneResponse, error)
}
type Service struct {
	repository      Repository
	provinceService ProvinceService
}

type Repository interface {
	Create(ctx context.Context, address entity.Address) (entity.Address, error)
	GetAll(ctx context.Context, userID uint) ([]entity.Address, error)
	GetOne(ctx context.Context, userID uint, addressID uint) (entity.Address, error)
	Delete(ctx context.Context, userID uint, addressID uint) error
	Update(ctx context.Context, userID uint, addressID uint, title *string, provinceID *uint, city, address, postalCode *string) error
}

func New(repository Repository, provinceService ProvinceService) Service {
	return Service{
		repository:      repository,
		provinceService: provinceService,
	}
}
