package entity

import "time"

type Session struct {
	ID uint
	UserID uint
	ExpiresAt time.Time
	// declare pointer of time.Time 
	// beacuse we need to check nil or not 
	RevokeAt *time.Time
}