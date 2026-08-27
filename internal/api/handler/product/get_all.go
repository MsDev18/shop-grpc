package product

import (
	dto "shop/internal/dto/product"
	"shop/internal/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetAll(ctx *gin.Context) {
	const op = "product-handler.Getall"

	limit, _ := strconv.Atoi(ctx.Query("limit"))
	page, _ := strconv.Atoi(ctx.Query("page"))

	var categorySlug *string
	if raw, exists := ctx.GetQuery("category"); exists && raw != "" {
		categorySlug = &raw
	}

	req := dto.GetAllRequest{
		Page:  page,
		Limit: limit,
		CategorySlug: categorySlug,
	}

	res, err := h.service.GetAll(ctx.Request.Context(), req)
	if err != nil {
		response.New(ctx).Error(err)
		return
	}

	response.New(ctx).OK("OK", res)
}
