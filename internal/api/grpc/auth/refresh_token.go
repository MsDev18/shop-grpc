package auth

import (
	"context"
	authpb "shop/internal/pb/auth"
	"shop/internal/pkg/mapper"
)

func (s Server) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {

	res, err := s.service.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, mapper.ErrorToGrpc(err)
	}

	return &authpb.RefreshTokenResponse{
		AccessToken: res.AccessToken,
	}, nil
}
