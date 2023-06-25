package middlewares

import (
	"errors"
	"go-clean-arch/constants"
	"go-clean-arch/infrastructure"
	"go-clean-arch/lib"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type JWTAuthMiddleware struct {
	logger        lib.Logger
	jwtAuthHelper *infrastructure.JWTAuthHelper
}

func NewJWTAuthMiddleware(
	logger lib.Logger,
	jwtHelper *infrastructure.JWTAuthHelper,
) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		logger,
		jwtHelper,
	}
}

func (m JWTAuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")

		if len(t) == 2 {
			authToken := t[1]
			user, err := m.jwtAuthHelper.Authorize(authToken)
			if user != nil {
				c.Set(constants.User, user)
				c.Next()
				return
			}
			lib.ErrorJSON(c, http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		lib.ErrorJSON(c, http.StatusUnauthorized, errors.New("you are not authorized"))
		c.Abort()
	}
}
