package auth

import (
	"shop/internal/entity"
	"shop/internal/pkg/response"
	"shop/internal/pkg/richerror"

	"github.com/gin-gonic/gin"
)

func (m Middleware) RoleRequired(requiredRole ...entity.Role) gin.HandlerFunc {
	const op = "auth-middleware.RoleRequired"
	return func(ctx *gin.Context) {
		value, exists := ctx.Get(ROLE_KEY)
		if !exists {
			response.New(ctx).Error(
				richerror.New().
					SetOp(op).
					SetMsg("not found role in request context").
					SetKind(richerror.KindForbiddenErr),
			)
			ctx.Abort()
			return
		}

		userRole, ok := value.(entity.Role)
		if !ok {
			response.New(ctx).Error(
				richerror.New().
					SetOp(op).
					SetMsg("invalid role type").
					SetKind(richerror.KindForbiddenErr),
			)
			ctx.Abort()
			return
		}

		var checkAccess bool = false
		for _, role := range requiredRole {
			if userRole == role {
				checkAccess = true
				break
			}
		}

		if !checkAccess {
			response.New(ctx).Error(
				richerror.New().
					SetOp(op).
					SetMsg("invalid role type").
					SetKind(richerror.KindForbiddenErr),
			)
			ctx.Abort()
			return 
		}

		ctx.Next()
	}
}
