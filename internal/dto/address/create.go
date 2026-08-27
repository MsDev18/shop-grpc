package address

type CreateRequest struct {
	Title      string `json:"title" binding:"required"`
	ProvinceID uint   `json:"province-id" binding:"required"`
	City       string `json:"city" binding:"required"`
	Address    string `json:"address" binding:"required"`
	PostalCode string `json:"postal-code" binding:"required"`
}

type CreateResponse struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Address    string `json:"address"`
	PostalCode string `json:"postal-code"`
}
