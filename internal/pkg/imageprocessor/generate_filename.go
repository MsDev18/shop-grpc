package imageprocessor

import (
	"shop/internal/pkg/richerror"

	"github.com/google/uuid"
)

func generateFileName() (string, error) {
	const op = "imageprocessor.generateFileName"
	id, err := uuid.NewRandom()
	if err != nil {
		return "", richerror.New().
			SetOp(op).
			SetMsg("can't generate file name").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	return id.String() + ".png" , nil
}
