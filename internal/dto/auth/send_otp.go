package auth

type SendOtpRequest struct {
	PhoneNumber string `json:"phone-number" binding:"required"`
}

type SendOtpResponse struct {}