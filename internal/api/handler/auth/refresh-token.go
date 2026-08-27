package auth

import (
	"net/http"
	"shop/internal/pkg/response"
	"shop/internal/pkg/richerror"

	"github.com/gin-gonic/gin"
)

func (h Handler) RefreshToken(ctx *gin.Context) {
	const op = "auth-handler.RefreshToken"

	refreshToken, err := ctx.Cookie("refresh-token")
	if err != nil {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("refresh token not found").
				SetKind(richerror.KindUnauthorizeErr).
				SetErr(err),
		)
		return
	}

	res, err := h.service.RefreshToken(ctx.Request.Context(), refreshToken)
	if err != nil {
		response.New(ctx).Error(err)
		return
	}

	response.New(ctx).Send(http.StatusOK, "access token refreshed" , res , nil)
}
