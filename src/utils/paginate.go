package utils

import (
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/ssamsara98/go-clean-arch/src/constants"
	"gorm.io/gorm"
)

type IPaginationMeta struct {
	Page       *int64 `json:"page"`
	Limit      *int64 `json:"limit"`
	Count      *int64 `json:"count"`
	ItemCount  *int64 `json:"itemCount"`
	TotalPages *int64 `json:"totalPages"`
}
type IPaginationCursorMeta struct {
	Cursor    *int64 `json:"cursor"`
	Limit     *int64 `json:"limit"`
	ItemCount *int64 `json:"itemCount"`
	HasNext   *bool  `json:"hasNext"`
}

type Pagination[M any, T any] struct {
	Meta  *M   `json:"meta"`
	Items *[]T `json:"items"`
}

func GetPaginationQuery(c *fiber.Ctx) (*int64, *int64) {
	limit, _ := c.Locals(constants.Limit).(int64)
	page, _ := c.Locals(constants.Page).(int64)
	return &limit, &page
}
func GetPaginationCursorQuery(c *fiber.Ctx) (*int64, *int64) {
	limit, _ := c.Locals(constants.Limit).(int64)
	cursor, _ := c.Locals(constants.Cursor).(int64)
	return &limit, &cursor
}

func Paginate(limit, page *int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := ((*page) - 1) * (*limit)
		return db.Offset(int(offset)).Limit(int(*limit))
	}
}
func PaginateCursor(limit *int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(int(*limit))
	}
}

func CreatePagination[T any](items *[]T, count *int64, limit *int64, page *int64) *Pagination[IPaginationMeta, T] {
	itemCount := int64(len(*items))
	totalPages := int64(math.Ceil(float64(*count) / float64(*limit)))

	meta := &IPaginationMeta{
		Page:       page,
		Limit:      limit,
		Count:      count,
		ItemCount:  &itemCount,
		TotalPages: &totalPages,
	}

	result := &Pagination[IPaginationMeta, T]{
		Meta:  meta,
		Items: items,
	}
	return result
}
func CreatePaginationCursor[T any](items *[]T, limit *int64, cursor *int64) *Pagination[IPaginationCursorMeta, T] {
	itemCount := int64(len(*items))
	hasNext := itemCount >= *limit

	meta := &IPaginationCursorMeta{
		Cursor:    cursor,
		Limit:     limit,
		ItemCount: &itemCount,
		HasNext:   &hasNext,
	}

	result := &Pagination[IPaginationCursorMeta, T]{
		Meta:  meta,
		Items: items,
	}
	return result
}
