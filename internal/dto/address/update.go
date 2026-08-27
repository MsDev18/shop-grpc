package address

type UpdateRequest struct {
	Title      *string `json:"title"`
	ProvinceID *uint   `json:"province-id"`
	City       *string `json:"city"`
	Address    *string `json:"address"`
	PostalCode *string `json:"postal-code" `
}
