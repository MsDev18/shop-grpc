package province

import (
	"shop/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetAll(ctx *gin.Context) {
	const op = "province-handler.GetAll"

	// call service
	res, err := h.service.GetAll(ctx.Request.Context())
	if err != nil {
		response.New(ctx).Error(err)
		return 
	}
	// return response
	response.New(ctx).OK("OK", res)
}
