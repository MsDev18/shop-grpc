package auth

import (
	"context"
	"shop/internal/entity"
	"shop/internal/pkg/mapper"
	"shop/internal/pkg/richerror"

	"google.golang.org/grpc"
)

var requiredRoles = map[string][]entity.Role{
	"/category.CategoryService/Create": {entity.AdminRole},
	"/category.CategoryService/Update": {entity.AdminRole},
	"/category.CategoryService/Delete": {entity.AdminRole},
	"/product.ProductService/Create":   {entity.AdminRole},
}

type RoleInterceptor struct{}

func NewRoleInterceptor() RoleInterceptor {
	return RoleInterceptor{}
}

func (RoleInterceptor) checkRole(ctx context.Context, fullMethod string) error {
	const op = "role-interceptor"

	allowed, exists := requiredRoles[fullMethod]
	if !exists {
		return nil
	}

	role, ok := RoleFromContext(ctx)
	if !ok {
		return mapper.ErrorToGrpc(
			richerror.New().
				SetOp(op).
				SetMsg("role not found in request context").
				SetKind(richerror.KindForbiddenErr),
		)
	}

	for _, r := range allowed {
		if role == r {
			return nil
		}
	}

	return mapper.ErrorToGrpc(
		richerror.New().
			SetOp(op).
			SetMsg("you don't have permission to perform this action").
			SetKind(richerror.KindForbiddenErr),
	)
}

func (i RoleInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := i.checkRole(ctx, info.FullMethod); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func (i RoleInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := i.checkRole(ss.Context(), info.FullMethod); err != nil {
			return err
		}
		return handler(srv, ss)
	}
}