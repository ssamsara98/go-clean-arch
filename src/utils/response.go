package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorJSON : json error response function
func ErrorJSON(c *gin.Context, err error, opts ...any) {
	statusCode := http.StatusInternalServerError
	if len(opts) > 0 {
		statusCode, _ = opts[0].(int)
	}
	resp := gin.H{
		"status":     "error",
		"statusCode": statusCode,
		"message":    err.Error(),
		"error":      err,
	}
	c.JSON(statusCode, resp)
}

// SuccessJSON : json error response function
func SuccessJSON(c *gin.Context, data any, opts ...any) {
	statusCode := http.StatusOK
	if len(opts) > 0 {
		statusCode, _ = opts[0].(int)
	}
	resp := gin.H{
		"status":     "success",
		"statusCode": statusCode,
		"result":     data,
	}
	c.JSON(statusCode, resp)
}
