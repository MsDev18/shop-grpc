package category

type CategoryResponse struct {
	ID       uint               `json:"id"`
	Title    string             `json:"title"`
	Slug     string             `json:"slug"`
	Image    *string            `json:"image"`
	Children []CategoryResponse `json:"children,omitempty"`
}
