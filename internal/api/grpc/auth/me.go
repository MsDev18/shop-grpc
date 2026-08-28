package auth

import (
	"context"
	authpb "shop/internal/pb/auth"
	"shop/internal/pkg/mapper"
	"shop/internal/pkg/richerror"
)

func (s Server) Me(ctx context.Context, req *authpb.MeRequest) (*authpb.MeResponse, error) {
	const op = "auth-grpc.Me"

	userID, ok := UserIDFromContext(ctx)
	if !ok {
		return nil, mapper.ErrorToGrpc(
			richerror.New().
				SetOp(op).
				SetMsg("not found user-id in request context").
				SetKind(richerror.KindUnauthorizeErr),
		)
	}

	return &authpb.MeResponse{
		UserId: uint64(userID),
	}, nil
}
