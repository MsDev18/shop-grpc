package product

import "mime/multipart"

type CreateRequest struct {
	Name        string
	Slug        string
	Description string
	Price       uint
	Stock       *uint
	CategoryID  uint
	MainImage   *multipart.FileHeader
	Images      []*multipart.FileHeader
}

type CreateResponse struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Price       uint     `json:"price"`
	Stock       uint     `json:"stock"`
	CategoryID  uint     `json:"category-id"`
	MainImage   string   `json:"main-image"`
	Images      []string `json:"images"`
}
