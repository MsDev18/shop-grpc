package auth

import (
	"context"
	authdto "shop/internal/dto/auth"
	"shop/internal/entity"
	"shop/internal/pkg/claims"
	"shop/internal/pkg/richerror"
	"time"
)

func (s Service) CheckOtp(ctx context.Context, req authdto.CheckOtpRequest) (authdto.CheckOtpResponse, error) {
	const op = "auth-service.CheckOtp"
	// 1. get user by phone number
	// 2. check user exist
	user, getUserErr := s.repository.GetUserByPhoneNumber(ctx, req.PhoneNumber)
	if getUserErr != nil {
		return authdto.CheckOtpResponse{}, getUserErr
	}
	// 3. if exist user get otp by user id
	otp, getOtpErr := s.repository.GetOtpByUserID(ctx, user.ID)
	if getOtpErr != nil {
		return authdto.CheckOtpResponse{}, getOtpErr
	}

	// 4. if otp.ExpiresAt < time.now return err
	if !otp.ExpiresAt.After(time.Now()) {
		return authdto.CheckOtpResponse{}, richerror.New().
			SetOp(op).
			SetMsg("otp code expired").
			SetKind(richerror.KindUnauthorizeErr)
	}
	// 5. check otp.code == req.code
	if otp.Code != req.Code {
		return authdto.CheckOtpResponse{}, richerror.New().
			SetOp(op).
			SetMsg("invalid otp code").
			SetKind(richerror.KindUnauthorizeErr)
	}

	// 6. create session 
	sessionData := entity.Session{
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(s.config.RefreshTokenDuration),
	}
	session, err := s.repository.CreateSession(ctx , sessionData)
	if err != nil {
		return authdto.CheckOtpResponse{}, err
	}
	// 7. generate jwt token
	accessToken, err := claims.CreateAccessToken(user.ID,session.ID, s.config.AccessTokenSecret, s.config.AccessTokenDuration)
	if err != nil {
		return authdto.CheckOtpResponse{}, err
	}
	refreshToken, err := claims.CreateRefreshToken(user.ID,session.ID, s.config.RefreshTokenSecret, s.config.RefreshTokenDuration)
	if err != nil {
		return authdto.CheckOtpResponse{}, err
	}

	// 78. return response
	return authdto.CheckOtpResponse{
		Tokens: authdto.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}
