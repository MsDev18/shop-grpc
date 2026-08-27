package localstorage

import (
	"context"
	"io"
	"os"
	"path"
	"path/filepath"
	"shop/internal/pkg/imageprocessor"
	"shop/internal/pkg/richerror"
)

func (l LocalStorage) Save(ctx context.Context, filename string, data io.Reader) (string, error) {
	const op = "localstorage.Save"

	// 1. make sure directory exists on disk
	diskDir := filepath.Join(imageprocessor.UPLOADS_ROOT, l.subDir)
	if err := os.MkdirAll(diskDir, 0755); err != nil {
		return "", richerror.New().
			SetOp(op).
			SetMsg("can't create upload directory").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	// 2. create the file and copy the incoming bytes into it
	diskPath := filepath.Join(diskDir, filename)
	file, err := os.Create(diskPath)
	if err != nil {
		return "", richerror.New().
			SetOp(op).
			SetMsg("can't create file on disk").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	defer file.Close()

	// 3. copy file bytes
	if _, err := io.Copy(file, data); err != nil {
		return "", richerror.New().
			SetOp(op).
			SetMsg("can't write file to disk").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	// 4. build the public URL - always forward slashes, regardless of OS
	urlPath := path.Join(imageprocessor.UPLOADS_URL_PREFIX , l.subDir ,filename)
	return urlPath, nil
}
