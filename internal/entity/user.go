package entity

import "time"

type Role string

const (
	AdminRole Role = "ADMIN"
	UserRole  Role = "USER"
)

type User struct {
	ID          uint
	Role        Role
	Name        string
	Avatar      string
	PhoneNumber string
	Password    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
