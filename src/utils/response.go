package utils

import (
	"github.com/gin-gonic/gin"
)

// ErrorJSON : json error response function
func ErrorJSON(c *gin.Context, statusCode int, err error) {
	resp := gin.H{
		"status":     "error",
		"statusCode": statusCode,
		"message":    err.Error(),
		"error":      err,
	}
	c.JSON(statusCode, resp)
}

// SuccessJSON : json error response function
func SuccessJSON(c *gin.Context, statusCode int, data any) {
	resp := gin.H{
		"status":     "success",
		"statusCode": statusCode,
		"result":     data,
	}
	c.JSON(statusCode, resp)
}
