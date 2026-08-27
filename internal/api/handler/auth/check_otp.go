package auth

import (
	"net/http"
	authdto "shop/internal/dto/auth"
	"shop/internal/pkg/response"
	"shop/internal/pkg/richerror"

	"github.com/gin-gonic/gin"
)

const (
	REFRESH_TOKEN_TTL = 60 * 60 *24 *30
)

func (h Handler) CheckOtp(ctx *gin.Context) {
	const op = "auth-handler.CheckOtp"

	// bind data
	var req authdto.CheckOtpRequest
	if bindingErr := ctx.ShouldBindJSON(&req); bindingErr != nil {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("can't bind body of request data").
				SetKind(richerror.KindBadRequestErr).
				SetErr(bindingErr),
		)
		return
	}

	// validation data
	if validationErr := h.validator.CheckOtp(ctx.Request.Context() ,req) ; validationErr != nil {
		response.New(ctx).Error(validationErr)
		return
	}

	// call service
	res, serviceErr := h.service.CheckOtp(ctx.Request.Context(), req)
	if serviceErr != nil {
		response.New(ctx).Error(serviceErr)
		return
	}

	// set refresh token in cookie 
	ctx.SetCookie(
		"refresh-token" , 
		res.Tokens.RefreshToken,
		REFRESH_TOKEN_TTL ,
		"/",
		"",
		true,
		true,
	)

	// return response 
	response.New(ctx).Send(http.StatusOK , "login successfully" , map[string]any{
		"access-token" : res.Tokens.AccessToken ,
	} , nil)
}
