package auth

import (
	"context"
	"shop/internal/entity"
	"time"
)

type Repository interface {
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	UpsertOtp(ctx context.Context, otp entity.Otp) (entity.Otp, error)
	GetOtpByUserID(ctx context.Context, userID uint) (entity.Otp, error)
	CreateSession(ctx context.Context, session entity.Session) (entity.Session, error)
	GetSessionByID(ctx context.Context, sessionID uint) (entity.Session, error)
	RevokeSession(ctx context.Context, sessionID uint) error
}

type Service struct {
	repository Repository
	config     Config
}
type Config struct {
	AccessTokenSecret    string        `koanf:"access_token_secret"`
	RefreshTokenSecret   string        `koanf:"refresh_token_secret"`
	AccessTokenDuration  time.Duration `koanf:"access_token_duration"`
	RefreshTokenDuration time.Duration `koanf:"refresh_token_duration"`
	OtpCodeDuration      time.Duration `koanf:"otp_code_duration"`
}

func New(repository Repository, config Config) Service {
	return Service{
		repository: repository,
		config:     config,
	}
}
