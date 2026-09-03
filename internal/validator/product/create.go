package product

import (
	"context"
	"errors"
	"net/http"
	dto "shop/internal/dto/product"
	"shop/internal/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) Create(ctx context.Context, req dto.CreateRequest) error {
	const op = "product-validator.Create"

	err := validation.ValidateStructWithContext(
		ctx,
		&req,
		validation.Field(&req.Name, validation.Required, validation.Length(NAME_MIN_LENGTH, NAME_MAX_LENGTH)),
		validation.Field(&req.Slug, validation.Required, validation.Length(SLUG_MIN_LENGTH, SLUG_MAX_LENGTH), validation.Match(SLUG_REGEX)),
		validation.Field(&req.Description, validation.Required),
		validation.Field(&req.Price, validation.Required),
		validation.Field(&req.Stock),
		validation.Field(&req.CategoryID, validation.Required),
		validation.Field(&req.MainImage, validation.Required, validation.By(v.validationImage)),
		validation.Field(&req.Images, validation.Each(validation.By(v.validationImage)), validation.Length(1, 0)),
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
				SetMsg("invalid inputs").
				SetKind(richerror.KindBadRequestErr).
				SetErr(err).
				SetMeta(meta)
		}
		return richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	return nil
}

func (v Validator) validationImage(value any) error {
	const op = "product-validator.validationImage"

	// var fileHeader *multipart.FileHeader
	// switch fh := value.(type) {
	// case *multipart.FileHeader:
	// 	fileHeader = fh
	// case multipart.FileHeader:
	// 	fileHeader = &fh
	// default:
	// 	return richerror.New().
	// 		SetOp(op).
	// 		SetMsg("type assertion error").
	// 		SetKind(richerror.KindBadRequestErr)
	// }

	image, ok := value.(*dto.ImageFile)
	if !ok || image == nil {
		return richerror.New().
			SetOp(op).
			SetMsg("invalid image").
			SetKind(richerror.KindBadRequestErr)
	}

	if len(image.Content) > 1024*300 {
		return richerror.New().
			SetOp(op).
			SetMsg("product image max size is 300KB").
			SetKind(richerror.KindBadRequestErr)
	}

	n := len(image.Content)
	if n > 512 {
		n = 512
	}

	contentType := http.DetectContentType(image.Content[:n])
	if contentType != "image/png" && contentType != "image/jpeg" {
		return richerror.New().
			SetOp(op).
			SetMsg("invalid image format").
			SetKind(richerror.KindBadRequestErr)
	}

	return nil
}
