package middlewares

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ssamsara98/go-clean-arch/src/constants"
	"github.com/ssamsara98/go-clean-arch/src/helpers"
	"github.com/ssamsara98/go-clean-arch/src/infrastructure"
	"github.com/ssamsara98/go-clean-arch/src/lib"
	"github.com/ssamsara98/go-clean-arch/src/models"
	"github.com/ssamsara98/go-clean-arch/src/utils"
	"gorm.io/gorm"
)

type JWTAuthMiddleware struct {
	logger  *lib.Logger
	JWTAuth *helpers.JWTAuth
	db      *infrastructure.Database
}

func NewJWTAuthMiddleware(
	logger *lib.Logger,
	jwtHelper *helpers.JWTAuth,
	db *infrastructure.Database,
) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		logger,
		jwtHelper,
		db,
	}
}

func (m JWTAuthMiddleware) Handle(tokenType string, needUser bool) fiber.Handler {
	m.logger.Debug("setting up jwt auth middleware")

	return func(c *fiber.Ctx) (err error) {
		authorizationHeader := c.Get(fiber.HeaderAuthorization)
		if authorizationHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "no token")
		} else if !strings.Contains(authorizationHeader, constants.TokenPrefix) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}

		tokenString := strings.Replace(authorizationHeader, constants.TokenPrefix+" ", "", -1)
		claims, err := m.JWTAuth.VerifyToken(tokenString, tokenType)
		if err != nil {
			m.logger.Error("claims error => ", err.Error())
			if errors.Is(err, jwt.ErrTokenExpired) {
				err = fiber.NewError(fiber.StatusForbidden, err.Error())
			} else {
				err = fiber.NewError(fiber.StatusUnauthorized, err.Error())
			}
			return err
		}
		if (claims.Type != constants.TokenAccess) && (claims.Type != constants.TokenRefresh) {
			return fiber.NewError(fiber.StatusUnauthorized, "wrong token type")
		}

		id, err := utils.ConvertStringToInt(claims.Subject)
		if err != nil {
			m.logger.Error("convert id error")
			return fiber.NewError(fiber.StatusUnauthorized, "you are not authorized")
		}

		if needUser {
			user := new(models.User)
			res := m.db.Where("id = ?", id).First(user)
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				return fiber.NewError(fiber.StatusUnauthorized, "user not found")
			}
			c.Locals(constants.User, user)
		} else {
			c.Locals(constants.User)
		}

		return c.Next()
	}
}
