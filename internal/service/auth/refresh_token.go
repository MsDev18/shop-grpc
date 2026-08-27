package auth

import (
	"context"
	authdto "shop/internal/dto/auth"
	"shop/internal/pkg/claims"
	"shop/internal/pkg/richerror"
	"strconv"
	"time"
)

func (s Service) RefreshToken(ctx context.Context, refreshToken string) (authdto.RefreshTokenResponse, error) {
	const op = "auth-service.RefreshToken"

	refreshClaims, err := claims.ParseRefreshToken(refreshToken, s.config.RefreshTokenSecret)
	if err != nil {
		return authdto.RefreshTokenResponse{}, err
	}

	userID, err := strconv.ParseUint(refreshClaims.Subject, 10, 64)
	if err != nil {
		return authdto.RefreshTokenResponse{}, richerror.New().
			SetOp(op).
			SetMsg("invalid subject claim in refresh token").
			SetKind(richerror.KindUnauthorizeErr).
			SetErr(err)
	}

	session, err := s.repository.GetSessionByID(ctx, refreshClaims.SessionID)
	if err != nil {
		return authdto.RefreshTokenResponse{}, err
	}

	if session.RevokeAt != nil || !session.ExpiresAt.After(time.Now()) {
		return authdto.RefreshTokenResponse{}, richerror.New().
			SetOp(op).
			SetMsg("session revoke or expired").
			SetKind(richerror.KindUnauthorizeErr)
	}

	accessToken, err := claims.CreateAccessToken(uint(userID), session.ID, s.config.AccessTokenSecret, s.config.AccessTokenDuration)
	if err != nil {
		return authdto.RefreshTokenResponse{}, err
	}

	return authdto.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}
