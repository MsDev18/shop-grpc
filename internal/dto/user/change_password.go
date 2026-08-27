package user

type ChangePasswordRequest struct {
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm-password" binding:"required"`
}
