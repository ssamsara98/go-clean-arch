package utils

import (
	"go-clean-arch/src/constants"
	"math"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IPaginationMeta struct {
	Page       *int64 `json:"page"`
	Limit      *int64 `json:"limit"`
	Count      *int64 `json:"count"`
	ItemCount  *int64 `json:"itemCount"`
	TotalPages *int64 `json:"totalPages"`
}

type Pagination[T any] struct {
	Meta  *IPaginationMeta `json:"meta"`
	Items *[]T             `json:"items"`
}

func GetPaginationQuery(c *gin.Context) (*int64, *int64) {
	limit, _ := c.MustGet(constants.Limit).(int64)
	page, _ := c.MustGet(constants.Page).(int64)
	return &limit, &page
}

func Paginate(limit, page *int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := ((*page) - 1) * (*limit)
		return db.Offset(int(offset)).Limit(int(*limit))
	}
}

func CreatePagination[T any](items *[]T, count *int64, limit *int64, page *int64) *Pagination[T] {
	itemCount := int64(len(*items))
	totalPages := int64(math.Ceil(float64(*count) / float64(*limit)))

	meta := &IPaginationMeta{
		Page:       page,
		Limit:      limit,
		Count:      count,
		ItemCount:  &itemCount,
		TotalPages: &totalPages,
	}

	result := &Pagination[T]{
		Meta:  meta,
		Items: items,
	}
	return result
}
