package auth

type CheckOtpRequest struct {
	PhoneNumber string `json:"phone-number" binding:"required"`
	Code        string `json:"code" binding:"required"`
}

type CheckOtpResponse struct{
	Tokens Tokens `json:"tokens"`
}


type Tokens struct {
	AccessToken  string `json:"access-token"`
	RefreshToken string `json:"refresh-token"`
}