package category

import (
	"io"
	dto "shop/internal/dto/category"
	"shop/internal/pkg/response"
	"shop/internal/pkg/richerror"

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

	var image *dto.ImageFile
	if fileHeader, err := ctx.FormFile("image"); fileHeader != nil && err == nil {
		file, err := fileHeader.Open()
		if err != nil {
			response.New(ctx).Error(
				richerror.New().
					SetOp(op).
					SetMsg("cna't open uploaded file").
					SetKind(richerror.KindBadRequestErr).
					SetErr(err),
			)
		}
		defer file.Close()

		content, readErr := io.ReadAll(file)
		if readErr != nil {
			response.New(ctx).Error(
				richerror.New().SetOp(op).SetMsg("can't read uploaded image").
					SetKind(richerror.KindUnexpectedErr).SetErr(readErr),
			)
			return
		}

		image = &dto.ImageFile{Filename: fileHeader.Filename, Content: content}
	}

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
	res, serviceErr := h.service.Update(ctx.Request.Context(), slug, req)
	if serviceErr != nil {
		response.New(ctx).Error(serviceErr)
		return
	}
	response.New(ctx).OK("update category successfully", res)
}
