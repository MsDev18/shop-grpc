package user

type ProfileResponse struct {
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	PhoneNumber string `json:"phone-number"`
}
