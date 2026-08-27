package imageprocessor

import (
	"bytes"
	"context"
	"image/png"
	"mime/multipart"
	"shop/internal/pkg/richerror"

	"github.com/disintegration/imaging"
)

func (p Processor) Process(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	const op = "imageprocessor.Process"

	// 1. open the upload file
	file, err := fileHeader.Open()
	if err != nil {
		return "", richerror.New().
			SetOp(op).
			SetMsg("can't open file").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	defer file.Close()

	// 2. decode - AutoOrientation fixes photos that look sideways, because
	// phone cameras store rotation as EXIF metadata instead of rotating
	// the actual pixels
	img, err := imaging.Decode(file, imaging.AutoOrientation(true))
	if err != nil {
		return "", richerror.New().
			SetOp(op).
			SetMsg("uploaded file is not valid image").
			SetKind(richerror.KindBadRequestErr).
			SetErr(err)
	}

	// 3. resize - Fit preserves aspect ratio and only shrinks, it never
	// upscales an already-small image
	resized := imaging.Fit(img, p.config.MaxDimension, p.config.MaxDimension, imaging.Lanczos)

	// 4. encode into an in-memory buffer instead of writing straight to
	// disk - this is exactly what makes this function storage-agnostic.
	// imaging.Encode only needs an io.Writer, and bytes.Buffer satisfies
	// that in memory, so Process never touches the filesystem itself.
	var buf bytes.Buffer
	if err := imaging.Encode(&buf, resized, imaging.PNG, imaging.PNGCompressionLevel(png.BestCompression)); err != nil {
		return "", richerror.New().
			SetOp(op).
			SetMsg("can't encode processed image ").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	// 5. generate a safe, unique filename - never trust the client-supplied name
	filename, err := generateFileName()
	if err != nil {
		return "", richerror.New().
			SetOp(op).
			SetMsg("can't generate file name").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	// 6. hand the encoded bytes to whatever storage backend is configured -
	// Process has no idea whether this ends up on local disk, S3, or
	// anything else.
	url, err := p.storage.Save(ctx, filename, &buf)
	if err != nil {
		return "", err
	}

	return url, nil
}
