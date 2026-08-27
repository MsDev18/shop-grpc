package entity

import "time"

type Otp struct {
	ID uint
	UserID uint
	Code string
	ExpiresAt time.Time
}