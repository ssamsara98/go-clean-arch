package middlewares

import (
	"errors"
	"go-clean-arch/constants"
	"go-clean-arch/infrastructure"
	"go-clean-arch/lib"
	"go-clean-arch/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
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

func (m JWTAuthMiddleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			t := strings.Split(authHeader, " ")

			if len(t) == 2 {
				authToken := t[1]
				user, err := m.jwtAuthHelper.Authorize(authToken)
				if user != nil {
					c.Set(constants.User, user)
					return next(c)
				}
				return utils.ErrorJSON(c, http.StatusInternalServerError, err)
			}

			return utils.ErrorJSON(c, http.StatusUnauthorized, errors.New("you are not authorized"))
		}
	}
}
