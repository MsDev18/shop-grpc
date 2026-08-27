package category

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	categorydto "shop/internal/dto/category"
	"shop/internal/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) Create(ctx context.Context, req categorydto.CreateRequest) error {
	const op = "category-validator.Create"
	// bussinuse rule , only parent category have image
	if req.ParentID != nil && req.Image != nil {
		return richerror.New().
			SetOp(op).
			SetMsg("child categories can't have an image").
			SetKind(richerror.KindBadRequestErr)
	}

	err := validation.ValidateStructWithContext(
		ctx,
		&req,
		validation.Field(&req.Title, validation.Required, validation.Length(TITLE_MIN_LENGTH, TITLE_MAX_LENGTH)),
		validation.Field(&req.Slug, validation.Required, validation.Length(SLUG_MIN_LENGTH, SLUG_MAX_LENGTH), validation.Match(SLUG_REGEX)),
		validation.Field(&req.Image, validation.By(v.validateImage)),
	)

	if err != nil {
		var validationErr validation.Errors
		if errors.As(err, &validationErr) {
			fieldErr := make(map[string]any, 0)
			for key, value := range validationErr {
				fieldErr[key] = value.Error()
			}
			return richerror.New().
				SetOp(op).
				SetMsg("invalid input").
				SetKind(richerror.KindBadRequestErr).
				SetErr(err).
				SetMeta(fieldErr)
		}
		return richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	return nil
}

func (v Validator) validateImage(value any) error {
	const op = "category-validator.validateImage"

	fileHeader, ok := value.(*multipart.FileHeader)
	if !ok {
		return richerror.New().
			SetOp(op).
			SetMsg("invalid image type").
			SetKind(richerror.KindBadRequestErr)
	}

	if fileHeader == nil {
		return nil
	}

	if fileHeader.Size > 1024*1024*1 {
		return richerror.New().
			SetOp(op).
			SetMsg("file max size is 1 MB").
			SetKind(richerror.KindBadRequestErr)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return richerror.New().
			SetOp(op).
			SetMsg("can't open file").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return richerror.New().
			SetOp(op).
			SetMsg("can't read file").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	contentType := http.DetectContentType(buf[:n])
	if contentType != "image/png" && contentType != "image/jpeg" {
		return richerror.New().
			SetOp(op).
			SetMsg("invalid image content type").
			SetKind(richerror.KindBadRequestErr)
	}
	return nil
}
