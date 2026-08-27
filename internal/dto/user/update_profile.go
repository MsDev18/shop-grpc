package user

import "mime/multipart"

type UpdateProfileRequest struct {
	Name   *string
	Avatar *multipart.FileHeader
}
