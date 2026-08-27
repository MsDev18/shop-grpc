package category

import "mime/multipart"

type CreateRequest struct {
	ParentID *uint
	Title string
	Slug string
	Image *multipart.FileHeader
}

type CreateResponse struct {
	ID uint `json:"id"`
	Title string `json:"title"`
	Slug string `json:"slug"`
	ParentID *uint `json:"parent-id"`
	Image *string `json:"image"`
}