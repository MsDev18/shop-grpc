package auth

import (
	authdto "shop/internal/dto/auth"
	"shop/internal/pkg/response"
	"shop/internal/pkg/richerror"

	"github.com/gin-gonic/gin"
)

func (h Handler) SendOtp(ctx *gin.Context) {
	const op = "auth-handler.SendOtp"
	// bind request body to json
	var req authdto.SendOtpRequest
	// in this section
	// we need to convert error type to richerror
	if bindingErr := ctx.ShouldBindJSON(&req); bindingErr != nil {
		response.New(ctx).Error(richerror.New().
			SetOp(op).
			SetMsg("can't bind body of request data").
			SetKind(richerror.KindBadRequestErr).
			SetErr(bindingErr),
		)
		return
	}
	// validate request body
	if validationErr := h.validator.SendOtp(ctx.Request.Context(), req); validationErr != nil {
		response.New(ctx).Error(validationErr)
		return
	}
	// call service
	res, serviceErr := h.service.SendOtp(ctx.Request.Context(), req)
	if serviceErr != nil {
		response.New(ctx).Error(serviceErr)
		return
	}
	response.New(ctx).OK("OTP sent successfully", res)
}
