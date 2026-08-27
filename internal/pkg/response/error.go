package response

import (
	"errors"
	"net/http"
	"shop/internal/pkg/mapper"
	"shop/internal/pkg/richerror"
)

func (r *Responder) Error(err error) {
	var richErr *richerror.RichError
	if errors.As(err, &richErr) {
		statusCode := mapper.KindToHttpStatusCode(richErr.GetKind())
		r.Send(statusCode, richErr.GetMessage(), nil, richErr.GetMeta())
		return
	}
	r.Send(http.StatusInternalServerError, "internal server error", nil, nil)
}
