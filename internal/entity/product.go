package entity

import "time"

type Product struct {
	ID          uint
	Name        string
	Slug        string
	Description string
	Price       uint
	Stock       uint
	MainImage   string
	CategoryID  uint
	DeletedAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
