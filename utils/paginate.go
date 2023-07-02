package utils

import (
	"go-clean-arch/constants"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Paginate(c echo.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		limit, _ := c.Get(constants.Limit).(int64)
		page, _ := c.Get(constants.Page).(int64)

		offset := (page - 1) * limit

		return db.Offset(int(offset)).Limit(int(limit))
	}
}
