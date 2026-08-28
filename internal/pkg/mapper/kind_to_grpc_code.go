package mapper

import (
	"shop/internal/pkg/richerror"

	"google.golang.org/grpc/codes"
)

func KindToGrpcCode(kind richerror.Kind) codes.Code {
	switch kind {
	case richerror.KindUnexpectedErr:
		return codes.Internal
	case richerror.KindConflictErr:
		return codes.AlreadyExists
	case richerror.KindNotFoundErr:
		return codes.NotFound
	case richerror.KindForbiddenErr:
		return codes.PermissionDenied
	case richerror.KindUnauthorizeErr:
		return codes.Unauthenticated
	case richerror.KindBadRequestErr:
		return codes.InvalidArgument
	default:
		return codes.Unknown
	}
}
