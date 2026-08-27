package product

type GetAllRequest struct {
	Page         int
	Limit        int
	CategorySlug *string
}

type ProductListItem struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	Price      uint   `json:"price"`
	Stock      uint   `json:"stock"`
	CategoryID uint   `json:"category-id"`
	MainImage  string `json:"main-image"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total-pages"`
}

type GetAllResponse struct {
	Products []ProductListItem `json:"products"`
	Meta     PaginationMeta    `json:"meta"`
}
