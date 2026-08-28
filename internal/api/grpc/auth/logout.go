package auth

import (
	"context"
	authpb "shop/internal/pb/auth"
	"shop/internal/pkg/mapper"
	"shop/internal/pkg/richerror"
)

func (s Server) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	const op = "auth-grpc.Logout"

	sessionID, ok := SessionIDFromContext(ctx)
	if !ok {
		return nil, mapper.ErrorToGrpc(
			richerror.New().
				SetOp(op).
				SetMsg("not found session id in request context").
				SetKind(richerror.KindUnexpectedErr),
		)
	}

	if serviceErr := s.service.Logout(ctx, sessionID); serviceErr != nil {
		return nil, mapper.ErrorToGrpc(serviceErr)
	}

	return &authpb.LogoutResponse{}, nil
}
