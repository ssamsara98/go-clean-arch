package models

import (
	"go-clean-arch/src/lib"
)

// User model
type Post struct {
	lib.ModelBase
	AuthorID    *uint  `json:"authorId"`
	Author      *User  `json:"author" gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	IsPublished bool   `json:"isPublished"`
}
