package auth

import (
	"context"
	dto "shop/internal/dto/auth"
	authpb "shop/internal/pb/auth"
	"shop/internal/pkg/mapper"
)

func (s Server) SendOtp(ctx context.Context, req *authpb.SendOtpRequest) (*authpb.SendOtpResponse, error) {
	const op = "auth-grpc.SendOtp"

	dtoReq := dto.SendOtpRequest{
		PhoneNumber: req.GetPhoneNumber(),
	}

	if validationErr := s.validator.SendOtp(ctx, dtoReq); validationErr != nil {
		return nil, mapper.ErrorToGrpc(validationErr)
	}

	if _, serviceErr := s.service.SendOtp(ctx, dtoReq); serviceErr != nil {
		return nil, mapper.ErrorToGrpc(serviceErr)
	}

	return &authpb.SendOtpResponse{}, nil
}
