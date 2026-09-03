package user

import (
	"io"
	authmiddleware "shop/internal/api/middleware/auth"
	userdto "shop/internal/dto/user"
	"shop/internal/pkg/response"
	"shop/internal/pkg/richerror"

	"github.com/gin-gonic/gin"
)

func (h Handler) UpdateProfile(ctx *gin.Context) {
	const op = "user-handler.UpdateProfile"

	value, exists := ctx.Get(authmiddleware.USER_ID_KEY)
	if !exists {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("can't find user-id in context od request").
				SetKind(richerror.KindUnauthorizeErr),
		)
		return
	}

	userID, ok := value.(uint)
	if !ok {
		response.New(ctx).Error(
			richerror.New().
				SetOp(op).
				SetMsg("invalid user-id").
				SetKind(richerror.KindUnauthorizeErr),
		)
		return
	}

	name , exists := ctx.GetPostForm("name")
	var namePtr *string
	if exists {
		namePtr = &name
	}
	// beacuse this field is optional, so we don't need to check error

	var avatar *userdto.AvatarFile
	if fileHeader, ferr := ctx.FormFile("avatar"); ferr == nil && fileHeader != nil {
		file, openErr := fileHeader.Open()
		if openErr != nil {
			response.New(ctx).Error(
				richerror.New().
					SetOp(op).
					SetMsg("can't open uploaded file").
					SetKind(richerror.KindBadRequestErr).
					SetErr(openErr),
			)
			return
		}
		defer file.Close()

		content, readErr := io.ReadAll(file)
		if readErr != nil {
			response.New(ctx).Error(
				richerror.New().
					SetOp(op).
					SetMsg("can't read uploaded file").
					SetKind(richerror.KindUnexpectedErr).
					SetErr(readErr),
			)
			return
		}
		avatar = &userdto.AvatarFile{Filename: fileHeader.Filename, Content: content}
	}

	req := userdto.UpdateProfileRequest{
		Name:   namePtr,
		Avatar: avatar,
	}

	// validation
	validationErr := h.validator.UpdateProfile(ctx.Request.Context(), req)
	if validationErr != nil {
		response.New(ctx).Error(validationErr)
		return
	}
	// service
	serviceErr := h.service.UpdateProfile(ctx.Request.Context(), userID, req)
	if serviceErr != nil {
		response.New(ctx).Error(serviceErr)
		return
	}
	// response
	// TODO - write response in body of request
	response.New(ctx).OK("edit your profile successfully.", nil)
}
