package user

import (
	authmiddleware "shop/internal/api/middleware/auth"
	"shop/internal/pkg/response"
	"shop/internal/pkg/richerror"

	"github.com/gin-gonic/gin"
)

func (h Handler) Profile(ctx *gin.Context) {
	const op = "user-handler.Profile"

	// get user-id from context of request
	value, exists := ctx.Get(authmiddleware.USER_ID_KEY)
	if !exists {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("can't find user-id in context od request").
				SetKind(richerror.KindUnauthorizeErr),
		)
		return
	}

	// parse to uint
	userID, ok := value.(uint)
	if !ok {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("invalid user-id").
				SetKind(richerror.KindUnauthorizeErr),
		)
		return 
	}
	// call service
	res, err := h.service.Profile(ctx.Request.Context() , userID)
	if err != nil {
		response.New(ctx).Error(err)
		return
	}
	// return response
	response.New(ctx).OK("" , res)
}
