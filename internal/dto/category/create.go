package category

type CreateRequest struct {
	ParentID *uint
	Title string
	Slug string
	Image *ImageFile
}

type CreateResponse struct {
	ID uint `json:"id"`
	Title string `json:"title"`
	Slug string `json:"slug"`
	ParentID *uint `json:"parent-id"`
	Image *string `json:"image"`
}