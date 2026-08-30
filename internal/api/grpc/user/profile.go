package user

import (
	"context"
	"shop/internal/api/grpc/auth"
	pb "shop/internal/pb/user"
	"shop/internal/pkg/mapper"
	"shop/internal/pkg/richerror"
)

func (s Server) Profile(ctx context.Context, req *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	const op = "user-grpc.Profile"

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, mapper.ErrorToGrpc(
			richerror.New().
				SetOp(op).
				SetMsg("user id not found in context").
				SetKind(richerror.KindUnexpectedErr),
		)
	}

	res, serviceErr := s.service.Profile(ctx, userID)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return &pb.ProfileResponse{
		Name:        res.Name,
		Avatar:      res.Avatar,
		PhoneNumber: res.PhoneNumber,
	} , nil
}
