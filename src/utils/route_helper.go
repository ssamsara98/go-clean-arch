package utils

import (
	"go-clean-arch/src/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindBody[T any](c *gin.Context) (*T, error) {
	body := new(T)
	if err := c.Bind(body); err != nil {
		ErrorJSON(c, err, http.StatusBadRequest)
		return nil, err
	}
	return body, nil
}

func BindUri[T any](c *gin.Context) (*T, error) {
	uri := new(T)
	if err := c.BindUri(uri); err != nil {
		ErrorJSON(c, err, http.StatusBadRequest)
		return nil, err
	}
	return uri, nil
}

func GetUser[T any](c *gin.Context) (*T, bool) {
	user, boolean := c.MustGet(constants.User).(*T)
	return user, boolean
}
