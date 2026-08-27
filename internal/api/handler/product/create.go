package product

import (
	dto "shop/internal/dto/product"
	"shop/internal/pkg/response"
	"shop/internal/pkg/richerror"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h Handler) Create(ctx *gin.Context) {
	const op = "product-handler.Create"

	name := ctx.PostForm("name")
	slug := ctx.PostForm("slug")
	description := ctx.PostForm("description")

	priceStr := ctx.PostForm("price")
	price, err := strconv.ParseUint(priceStr, 10, 64)
	if err != nil {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("price must be integer").
				SetKind(richerror.KindBadRequestErr).
				SetErr(err),
		)
		return
	}

	var stock *uint
	stockStr, exists := ctx.GetPostForm("stock")
	if exists {
		stockUint64, err := strconv.ParseUint(stockStr, 10, 64)
		if err != nil {
			response.New(ctx).Error(
				richerror.New().
					SetOp(op).
					SetMsg("stock must be integer").
					SetKind(richerror.KindBadRequestErr).
					SetErr(err),
			)
			return
		}
		stockUint := uint(stockUint64)
		stock = &stockUint
	}

	categoryIDStr := ctx.PostForm("category-id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("category-id must be integer").
				SetKind(richerror.KindBadRequestErr).
				SetErr(err),
		)
		return
	}

	mainImage, _ := ctx.FormFile("main-image")

	form, err := ctx.MultipartForm()
	if err != nil {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("invalid form data").
				SetKind(richerror.KindBadRequestErr).
				SetErr(err),
		)
		return
	}

	images := form.File["images"]

	// create dto.CreateResponse
	req := dto.CreateRequest{
		Name:        name,
		Slug:        slug,
		Description: description,
		Price:       uint(price),
		Stock:       stock,
		CategoryID:  uint(categoryID),
		MainImage:   mainImage,
		Images:      images,
	}
	// validation
	if validationErr := h.validator.Create(ctx.Request.Context(), req); validationErr != nil {
		response.New(ctx).Error(validationErr)
		return
	}

	// service
	res, serviceErr := h.service.Create(ctx.Request.Context(), req)
	if serviceErr != nil {
		response.New(ctx).Error(serviceErr)
		return
	}

	// return response
	response.New(ctx).Created("create product successfully", res)
}
