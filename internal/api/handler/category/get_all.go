package category

import (
	"shop/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetAll (ctx *gin.Context) {
	const op = "category-handler.GetAll"
	// call service 
	categories ,err := h.service.GetAll(ctx.Request.Context())
	if err != nil {
		response.New(ctx).Error(err)
		return
	}
	// return response
	response.New(ctx).OK("ok" , categories)
}