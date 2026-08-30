package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	userdto "shop/internal/dto/user"
	"shop/internal/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) UpdateProfile(ctx context.Context, req userdto.UpdateProfileRequest) error {
	const op = "user-validator.UpdateProfile"

	err := validation.ValidateStructWithContext(
		ctx,
		&req,
		validation.Field(&req.Name, validation.Length(3, 50)),
		validation.Field(&req.Avatar, validation.By(v.validationAvatarImage)),
	)

	if err != nil {
		var validationErr validation.Errors
		if errors.As(err, &validationErr) {
			meta := make(map[string]any)
			for key, value := range validationErr {
				meta[key] = value.Error()
			}
			return richerror.New().
				SetOp(op).
				SetMsg("input validation error").
				SetKind(richerror.KindBadRequestErr).
				SetErr(err).
				SetMeta(meta)
		}
		return richerror.New().
			SetOp(op).
			SetMsg("unexpected err").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	return nil
}

func (v Validator) validationAvatarImage(value interface{}) error {
	const op = "user-validator.validationAvatarImage"
	// 1. check value type and type assertion
	avatar, ok := value.(*userdto.AvatarFile)
	if !ok {
		return richerror.New().
			SetOp(op).
			SetMsg("type assertion error").
			SetKind(richerror.KindBadRequestErr)
	}

	// check fileHeader is nil, if nil -> no file uploaded, return nil
	if avatar == nil {
		return nil
	}

	// 2. check file size
	if len(avatar.Content) > 1*1024*1024 {
		return richerror.New().
			SetOp(op).
			SetMsg(fmt.Sprintf("file max size is %d MB", 1)).
			SetKind(richerror.KindBadRequestErr)
	}


	n := len(avatar.Content)
	if n > 512 {
		n = 512
	}

	contentType := http.DetectContentType(avatar.Content[:n])
	if contentType != "image/png" && contentType != "image/jpeg" {
		return richerror.New().
			SetOp(op).
			SetMsg("only jpeg and png allowed to upload").
			SetKind(richerror.KindBadRequestErr)
	}

	return nil
}
