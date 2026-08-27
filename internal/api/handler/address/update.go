package address

import (
	authmiddleware "shop/internal/api/middleware/auth"
	dto "shop/internal/dto/address"
	"shop/internal/pkg/response"
	"shop/internal/pkg/richerror"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h Handler) Update(ctx *gin.Context) {
	const op = "address-handler.Update"
	// get address id
	addressIDStr := ctx.Param("id")
	addressID, err := strconv.ParseUint(addressIDStr, 10, 64)
	if err != nil {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("invalid address id").
				SetKind(richerror.KindBadRequestErr).
				SetErr(err),
		)
		return
	}

	// get user id
	value, exists := ctx.Get(authmiddleware.USER_ID_KEY)
	if !exists {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("not fund user id in request context").
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

	// bind request
	var req dto.UpdateRequest
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

	// call validator
	if validationErr := h.validator.Update(ctx.Request.Context(), req); validationErr != nil {
		response.New(ctx).Error(validationErr)
		return
	}

	// call service
	if serviceErr := h.service.Update(ctx.Request.Context(), userID, uint(addressID), req); serviceErr != nil {
		response.New(ctx).Error(serviceErr)
		return
	}

	// return response
	response.New(ctx).OK("updated successfully", nil)
}
