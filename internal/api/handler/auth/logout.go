package auth

import (
	authmiddleware "shop/internal/api/middleware/auth"
	"shop/internal/pkg/response"
	"shop/internal/pkg/richerror"

	"github.com/gin-gonic/gin"
)

func (h Handler) Logout(ctx *gin.Context) {
	const op = "auth-handler.Logout"

	value, exists := ctx.Get(authmiddleware.SESSION_ID_KEY)
	if !exists {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("session id not found in context").
				SetKind(richerror.KindUnexpectedErr),
		)
		return
	}

	sessionID, ok := value.(uint)
	if !ok {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("invalid session id type in context").
				SetKind(richerror.KindUnexpectedErr),
		)
		return
	}

	if err := h.service.Logout(ctx.Request.Context() , sessionID) ; err != nil {
		response.New(ctx).Error(err)
		return
	}

	ctx.SetCookie("refresh-token" , "" , -1, "/", "" ,true, true)
	response.New(ctx).OK("logout successfully" , nil)
}
