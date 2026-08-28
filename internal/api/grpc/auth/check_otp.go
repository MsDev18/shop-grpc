package auth

import (
	"context"
	dto "shop/internal/dto/auth"
	authpb "shop/internal/pb/auth"
	"shop/internal/pkg/mapper"
)

func (s Server) CheckOtp(ctx context.Context, req *authpb.CheckOtpRequest) (*authpb.CheckOtpResponse, error) {
	const op = "auth-server.CheckOtp"

	dtoReq := dto.CheckOtpRequest{
		PhoneNumber: req.GetPhoneNumber(),
		Code:        req.GetCode(),
	}

	if validationErr := s.validator.CheckOtp(ctx, dtoReq); validationErr != nil {
		return nil, mapper.ErrorToGrpc(validationErr)
	}

	res, serviceErr := s.service.CheckOtp(ctx, dtoReq)
	if serviceErr != nil {
		return nil, mapper.ErrorToGrpc(serviceErr)
	}

	return &authpb.CheckOtpResponse{
		AccessToken:  res.Tokens.AccessToken,
		RefreshToken: res.Tokens.RefreshToken,
	}, nil
}
