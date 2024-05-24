package models

import (
	"go-clean-arch/src/lib"
)

// User model
type Post struct {
	lib.ModelBase
	AuthorId    uint   `json:"authorId"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	IsPublished bool   `json:"isPublished"`
}
