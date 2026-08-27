package auth

import (
	authmiddleware "shop/internal/api/middleware/auth"
	"shop/internal/pkg/response"
	"shop/internal/pkg/richerror"

	"github.com/gin-gonic/gin"
)

func (h Handler) Me (ctx *gin.Context) {
	const op = "auth-handler.Me"

	value , exists := ctx.Get(authmiddleware.USER_ID_KEY)
	if !exists {
		response.New(ctx).Error(richerror.New().
			SetOp(op).
			SetMsg("user id not found in context").
			SetKind(richerror.KindUnexpectedErr),
		)
		return
	}
	
	userID, ok := value.(uint)
	if !ok {
		response.New(ctx).Error(richerror.New().
			SetOp(op).
			SetMsg("invalid user id type in context").
			SetKind(richerror.KindUnexpectedErr),
		)
		return
	}

	
	response.New(ctx).OK("/auth/me route successfully" , map[string]any{
		"user-id" : userID,
	})
}