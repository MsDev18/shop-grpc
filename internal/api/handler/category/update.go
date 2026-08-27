package category

import (
	dto "shop/internal/dto/category"
	"shop/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h Handler) Update(ctx *gin.Context) {
	const op = "category-handler.Update"
	// get slug param
	slug := ctx.Param("slug")
	// get request data
	var titlePtr *string
	if title, exists := ctx.GetPostForm("title"); exists {
		titlePtr = &title
	}

	var updatedSlugPtr *string
	if updatedSlug, exists := ctx.GetPostForm("slug"); exists {
		updatedSlugPtr = &updatedSlug
	}

	image, _ := ctx.FormFile("image")

	// map to dto.Update
	req := dto.UpdateRequest{
		Title: titlePtr,
		Slug:  updatedSlugPtr,
		Image: image,
	}

	// validator
	if validationErr := h.validator.Update(ctx.Request.Context(), req); validationErr != nil {
		response.New(ctx).Error(validationErr)
		return
	}

	// service
	res ,serviceErr := h.service.Update(ctx.Request.Context(), slug, req)
	if serviceErr != nil {
		response.New(ctx).Error(serviceErr)
		return
	}
	response.New(ctx).OK("update category successfully", res)
}
