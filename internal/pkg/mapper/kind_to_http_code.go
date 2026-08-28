package mapper

import (
	"net/http"
	"shop/internal/pkg/richerror"
)

func KindToHttpStatusCode (kind richerror.Kind) int {
	switch kind {
	case richerror.KindUnexpectedErr:
		return http.StatusInternalServerError
	case richerror.KindConflictErr:
		return http.StatusConflict
	case richerror.KindNotFoundErr:
		return  http.StatusNotFound
	case richerror.KindForbiddenErr:
		return http.StatusForbidden
	case richerror.KindUnauthorizeErr:
		return http.StatusUnauthorized
	case richerror.KindBadRequestErr:
		return http.StatusBadRequest
	default:
		return 0
	}
}