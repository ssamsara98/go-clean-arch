package utils

import (
	"go-clean-arch/src/constants"

	"github.com/gin-gonic/gin"
)

func BindBody[T any](c *gin.Context) (*T, error) {
	var body T
	err := c.Bind(&body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func BindUri[T any](c *gin.Context) (*T, error) {
	var uri T
	err := c.ShouldBindUri(&uri)
	if err != nil {
		return nil, err
	}
	return &uri, nil
}

func GetUser[T any](c *gin.Context) (*T, bool) {
	user, boolean := c.MustGet(constants.User).(*T)
	return user, boolean
}
