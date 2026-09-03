package address

import (
	authmiddleware "shop/internal/api/middleware/auth"
	"shop/internal/pkg/response"
	"shop/internal/pkg/richerror"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetOne(ctx *gin.Context) {
	const op = "address-handler.GetOne"

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

	userID, ok := value.(uint)
	if !ok {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("invalid user id type").
				SetKind(richerror.KindUnauthorizeErr),
		)
		return
	}

	// get address id from param
	addressIDStr := ctx.Param("id")
	addressID, err := strconv.ParseUint(addressIDStr, 10, 64)
	if err != nil {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("invalid address id type").
				SetKind(richerror.KindBadRequestErr).
				SetErr(err),
		)
		return
	}

	// call service
	res, err := h.service.GetOne(ctx.Request.Context(), userID, uint(addressID))
	if err != nil {
		response.New(ctx).Error(err)
		return
	}
	// return response
	response.New(ctx).OK("OK", res)
}
