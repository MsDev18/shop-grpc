package product

import (
	"shop/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetOneBySlug(ctx *gin.Context) {
	const op = "product-handler.GetOneBySlug"

	slug := ctx.Param("slug")

	// we don't need validation
	// call service 
	res, err := h.service.GetOneBySlug(ctx.Request.Context(), slug)
	if err != nil {
		response.New(ctx).Error(err)
		return
	}

	response.New(ctx).OK("OK" , res)
}