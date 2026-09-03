package address

import (
	authmiddleware "shop/internal/api/middleware/auth"
	"shop/internal/pkg/response"
	"shop/internal/pkg/richerror"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetAll(ctx *gin.Context) {
	const op = "address-handler.GetAll"

	value, exists := ctx.Get(authmiddleware.USER_ID_KEY)
	if !exists {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("not found user id in request context").
				SetKind(richerror.KindUnauthorizeErr),
		)
		return
	}

	userID , ok := value.(uint)
	if !ok {
		response.New(ctx).Error(
			richerror.New().
			SetOp(op).
			SetMsg("invalid user id type").
			SetKind(richerror.KindUnauthorizeErr),
		)
		return
	}

	// service 
	res, err := h.service.GetAll(ctx.Request.Context(), userID)
	if err != nil {
		response.New(ctx).Error(err)
		return
	}
	response.New(ctx).OK("ok", res)
}
