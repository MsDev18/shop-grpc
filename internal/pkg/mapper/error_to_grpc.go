package mapper

import (
	"shop/internal/pkg/richerror"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorToGrpc(err error) error {
	if err == nil {
		return nil
	}

	richErr , ok := err.(*richerror.RichError)
	if !ok {
		return status.Error(codes.Unknown , err.Error())
	}

	return status.Error(KindToGrpcCode(richErr.GetKind()) , richErr.GetMessage())
}