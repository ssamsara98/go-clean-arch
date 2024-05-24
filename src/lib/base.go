package lib

import (
	"time"

	"gorm.io/gorm"
)

type ModelBase struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	ID        uint           `json:"id" gorm:"primarykey"`
}
