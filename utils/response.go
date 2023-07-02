package utils

import (
	"go-clean-arch/constants"

	"github.com/labstack/echo/v4"
)

// JSON : json response function
func JSON(c echo.Context, statusCode int, data any) error {
	return c.JSON(statusCode, echo.Map{
		"statusCode": statusCode,
		"result":     data,
	})
}

// ErrorJSON : json error response function
func ErrorJSON(c echo.Context, statusCode int, data error) error {
	return c.JSON(statusCode, echo.Map{
		"statusCode": statusCode,
		"message":    data.Error(),
		"error":      data,
	})
}

// SuccessJSON : json error response function
func SuccessJSON(c echo.Context, statusCode int, data any) error {
	return c.JSON(statusCode, echo.Map{
		"statusCode": statusCode,
		"message":    data,
	})
}

// JSONWithPagination : json response function
func JSONWithPagination(c echo.Context, statusCode int, response echo.Map) error {
	limit, _ := c.Get(constants.Limit).(int64)
	page, _ := c.Get(constants.Page).(int64)

	return c.JSON(
		statusCode,
		echo.Map{
			"result": response["result"],
			"pagination": echo.Map{
				"hasNext": (response["count"].(int64) - limit*page) > 0,
				"count":   response["count"],
			},
		})
}
