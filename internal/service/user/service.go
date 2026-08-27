package user

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
	GetUserByID(ctx context.Context, userID uint) (entity.User, error)
	UpdateProfile(ctx context.Context, userID uint, name *string, avatar *string) error
	UpdatePassword(ctx context.Context ,user entity.User) error
}

func New(repository Repository, imageProcessor imageprocessor.Processor) Service {
	return Service{
		repository:     repository,	
		imageProcessor: imageProcessor,
	}
}
