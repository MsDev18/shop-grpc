package category

import (
	"shop/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetOne(ctx *gin.Context) {
	const op = "category-handler.GetOne"
	// get slug 
	slug := ctx.Param("slug")
	// call service 
	category , err := h.service.GetOne(ctx.Request.Context(), slug)
	if err != nil {
		response.New(ctx).Error(err)
		return
	}
	// return response 
	response.New(ctx).OK("ok", category)
}