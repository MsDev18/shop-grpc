package entity

import "time"

type Address struct {
	ID         uint
	UserID     uint
	Title      string
	ProvinceID uint
	City       string
	Address    string
	PostalCode string
	DeletedAt  *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
