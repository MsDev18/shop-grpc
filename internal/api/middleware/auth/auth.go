package auth

import (
	"context"
	"shop/internal/entity"
)

const (
	USER_ID_KEY    = "user-id"
	SESSION_ID_KEY = "session-id"
	ROLE_KEY       = "role"
)

type Repository interface {
	GetSessionByID(ctx context.Context, sessionID uint) (entity.Session, error)
	GetUserByID(ctx context.Context, userID uint) (entity.User, error)
}

type Middleware struct {
	accessTokenSecret string
	repository        Repository
}

func New(repository Repository, accessTokenSecret string) Middleware {
	return Middleware{
		accessTokenSecret: accessTokenSecret,
		repository:        repository,
	}
}
