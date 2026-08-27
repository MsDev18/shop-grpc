package entity

import "time"

type Category struct {
	ID        uint
	ParentID  *uint
	Title     string
	Slug      string
	Image     *string
	DeletedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
