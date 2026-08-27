package province

import (
	"shop/internal/pkg/response"
	"shop/internal/pkg/richerror"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetOne(ctx *gin.Context) {
	const op = "province-handler.GetOne"

	// get id from param
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("invalid province id").
				SetKind(richerror.KindBadRequestErr).
				SetErr(err),
		)
		return
	}
	// call service
	res, err := h.service.GetOne(ctx.Request.Context(), uint(id))
	if err != nil {
		response.New(ctx).Error(err)
		return
	}

	response.New(ctx).OK("OK", res)
}
