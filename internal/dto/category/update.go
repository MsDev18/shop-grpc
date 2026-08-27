package category

import "mime/multipart"

type UpdateRequest struct {
	Title *string
	Slug  *string
	Image *multipart.FileHeader
}
