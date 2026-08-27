package entity

import "time"

type ProductImage struct {
	ID        uint
	Image     string
	ProductID uint
	DeletedAt *time.Time
	CreatedAt time.Time
}
