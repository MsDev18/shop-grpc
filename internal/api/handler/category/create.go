package category

import (
	categorydto "shop/internal/dto/category"
	"shop/internal/pkg/response"
	"shop/internal/pkg/richerror"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h Handler) Create(ctx *gin.Context) {
	const op = "category-handler.Create"
	// get data from form 
	title := ctx.PostForm("title")
	slug := ctx.PostForm("slug")
	var parentID *uint
	if raw , exists := ctx.GetPostForm("parent-id") ; exists && raw != "" {
		id , err := strconv.ParseUint(raw , 10, 64)
		if err != nil {
			response.New(ctx).Error(
				richerror.New().
				SetOp(op).
				SetMsg("invalid parent-id").
				SetKind(richerror.KindBadRequestErr).
				SetErr(err),
			)
			return 
		}
		v := uint(id)
		parentID = &v
	}
	image , _ := ctx.FormFile("image")

	// create CreateRequestDTO 
	req := categorydto.CreateRequest{
		ParentID: parentID,
		Title:    title,
		Slug:     slug,
		Image:    image,
	}
	// call validator
	err := h.validator.Create(ctx.Request.Context(),req)
	if err != nil {
		response.New(ctx).Error(err)
		return 
	}
	// call service 
	res , err := h.service.Create(ctx.Request.Context() , req)
	if err != nil {
		response.New(ctx).Error(err)
		return
	}
	// response
	response.New(ctx).Created("create category successfully" , res)
}