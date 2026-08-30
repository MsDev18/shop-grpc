package user

import (
	"context"
	"shop/internal/api/grpc/auth"
	dto "shop/internal/dto/user"
	pb "shop/internal/pb/user"
	"shop/internal/pkg/mapper"
	"shop/internal/pkg/richerror"
)

func (s Server) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	const op = "user-grpc.UpdateProfile"

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, mapper.ErrorToGrpc(
			richerror.New().
				SetOp(op).
				SetMsg("user id not found in context").
				SetKind(richerror.KindUnexpectedErr),
		)
	}

	var avatar *dto.AvatarFile
	if len(req.GetAvatar()) > 0 {
		avatar = &dto.AvatarFile{Content: req.GetAvatar()}
	}

	dtoReq := dto.UpdateProfileRequest{
		Name:   req.Name,
		Avatar: avatar,
	}

	if validationErr := s.validator.UpdateProfile(ctx, dtoReq); validationErr != nil {
		return nil, mapper.ErrorToGrpc(validationErr)
	}

	if serviceErr := s.service.UpdateProfile(ctx, userID, dtoReq); serviceErr != nil {
		return nil, mapper.ErrorToGrpc(serviceErr)
	}

	return &pb.UpdateProfileResponse{}, nil
}
