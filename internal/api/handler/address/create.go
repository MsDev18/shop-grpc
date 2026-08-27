package address

import (
	authmiddleware "shop/internal/api/middleware/auth"
	dto "shop/internal/dto/address"
	"shop/internal/pkg/response"
	"shop/internal/pkg/richerror"

	"github.com/gin-gonic/gin"
)

func (h Handler) Create(ctx *gin.Context) {
	const op = "address-handler.Create"
	// get user id for request context
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
	// type assertion any to uint
	userID, ok := value.(uint)
	if !ok {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("invalid user id format").
				SetKind(richerror.KindUnauthorizeErr),
		)
		return
	}
	// bind request
	var req dto.CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("can't bind request body").
				SetKind(richerror.KindBadRequestErr).
				SetErr(err),
		)
		return
	}
	// validation inputs
	if err := h.validator.Create(ctx.Request.Context(), req); err != nil {
		response.New(ctx).Error(err)
		return
	}
	// call service 
	res, err := h.service.Create(ctx.Request.Context(), userID, req)
	if err != nil {
		response.New(ctx).Error(err)
		return
	}
	response.New(ctx).Created("create address successfully", res)
}
