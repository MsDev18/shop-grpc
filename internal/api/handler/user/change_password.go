package user

import (
	authmiddleware "shop/internal/api/middleware/auth"
	userdto "shop/internal/dto/user"
	"shop/internal/pkg/response"
	"shop/internal/pkg/richerror"

	"github.com/gin-gonic/gin"
)

func (h Handler) ChangePassword(ctx *gin.Context) {
	const op = "user-handler.ChangePassword"

	value , isExist := ctx.Get(authmiddleware.USER_ID_KEY)
	if !isExist {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("user id not found in context").
				SetKind(richerror.KindUnauthorizeErr),
		)
		return
	}

	userID , ok:= value.(uint)
	if !ok {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("failed to assert user id").
				SetKind(richerror.KindUnexpectedErr),
		)
		return
	}

	var req userdto.ChangePasswordRequest
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

	// validator
	if validationErr := h.validator.ChangePassword(ctx.Request.Context(), req); validationErr != nil {
		response.New(ctx).Error(validationErr)
		return
	}
	// service
	if serivceErr := h.service.ChangePassword(ctx.Request.Context() , userID , req) ; serivceErr != nil {
		response.New(ctx).Error(serivceErr)
		return
	}
	// response
	response.New(ctx).OK("change password successfully" , nil)
}
