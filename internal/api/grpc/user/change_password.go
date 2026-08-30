package user

import (
	"context"
	"shop/internal/api/grpc/auth"
	dto "shop/internal/dto/user"
	pb "shop/internal/pb/user"
	"shop/internal/pkg/mapper"
	"shop/internal/pkg/richerror"
)

func (s Server) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	const op = "user-grpc.ChangePassword"

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, mapper.ErrorToGrpc(
			richerror.New().
				SetOp(op).
				SetMsg("not found user id in request context").
				SetKind(richerror.KindUnexpectedErr),
		)
	}

	dtoReq := dto.ChangePasswordRequest{
		Password:        req.GetPassword(),
		ConfirmPassword: req.GetConfirmPassword(),
	}

	if validationErr := s.validator.ChangePassword(ctx, dtoReq); validationErr != nil {
		return nil, mapper.ErrorToGrpc(validationErr)
	}

	if serviceErr := s.service.ChangePassword(ctx, userID, dtoReq); serviceErr != nil {
		return nil, mapper.ErrorToGrpc(serviceErr)
	}

	return &pb.ChangePasswordResponse{}, nil

}
