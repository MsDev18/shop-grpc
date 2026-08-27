package category

import (
	"shop/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h Handler) Delete(ctx *gin.Context) {
	const op = "category-handler.Delete"

	slug := ctx.Param("slug")

	// we dont need to call validator
	// call service
	if err := h.service.Delete(ctx.Request.Context(), slug); err != nil {
		response.New(ctx).Error(err)
		return
	}
	response.New(ctx).OK("deleted successfully", nil)
}
