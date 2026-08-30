package auth

import (
	"context"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
)

func RequireRole(ctx context.Context, allowed ...entity.Role) error {
	const op = "auth-grpc.RequireRole"

	role, ok := RoleFromContext(ctx)
	if !ok {
		return richerror.New().
			SetOp(op).
			SetMsg("role not found in request context").
			SetKind(richerror.KindForbiddenErr)
	}

	for _, r := range allowed {
		if role == r {
			return nil
		}
	}

	return richerror.New().
		SetOp(op).
		SetMsg("you don't have permission to perform this action").
		SetKind(richerror.KindForbiddenErr)
}
