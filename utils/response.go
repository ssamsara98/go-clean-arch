package utils

import (
	"go-clean-arch/constants"

	"github.com/gin-gonic/gin"
)

// JSON : json response function
func JSON(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, gin.H{
		"statusCode": statusCode,
		"result":     data,
	})
}

// ErrorJSON : json error response function
func ErrorJSON(c *gin.Context, statusCode int, data error) {
	c.JSON(statusCode, gin.H{
		"statusCode": statusCode,
		"message":    data.Error(),
		"error":      data,
	})
}

// SuccessJSON : json error response function
func SuccessJSON(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, gin.H{
		"statusCode": statusCode,
		"message":    data,
	})
}

// JSONWithPagination : json response function
func JSONWithPagination(c *gin.Context, statusCode int, response gin.H) {
	limit, _ := c.MustGet(constants.Limit).(int64)
	page, _ := c.MustGet(constants.Page).(int64)

	c.JSON(
		statusCode,
		gin.H{
			"result": response["result"],
			"pagination": gin.H{
				"hasNext": (response["count"].(int64) - limit*page) > 0,
				"count":   response["count"],
			},
		})
}
