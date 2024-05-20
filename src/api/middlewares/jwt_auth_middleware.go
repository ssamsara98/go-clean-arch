package middlewares

import (
	"errors"
	"go-clean-arch/src/constants"
	"go-clean-arch/src/infrastructure"
	"go-clean-arch/src/lib"
	"go-clean-arch/src/models"
	"go-clean-arch/src/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JWTAuthMiddleware struct {
	logger        *lib.Logger
	jwtAuthHelper *infrastructure.JWTAuthHelper
	db            *infrastructure.Database
}

func NewJWTAuthMiddleware(
	logger *lib.Logger,
	jwtHelper *infrastructure.JWTAuthHelper,
	db *infrastructure.Database,
) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		logger,
		jwtHelper,
		db,
	}
}

func (m *JWTAuthMiddleware) Handle() gin.HandlerFunc {
	m.logger.Debug("Setting up jwt auth middleware")

	return func(c *gin.Context) {
		authorizationHeader := c.Request.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, constants.TokenPrefix) {
			utils.ErrorJSON(c, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		tokenString := strings.Replace(authorizationHeader, constants.TokenPrefix+" ", "", -1)

		claims, err := m.jwtAuthHelper.VerifyToken(tokenString)
		if err != nil {
			m.logger.Error("claims error")
			utils.ErrorJSON(c, http.StatusUnauthorized, err)
			c.Abort()
			return
		}
		if claims.Type != constants.TokenAccess {
			utils.ErrorJSON(c, http.StatusUnauthorized, errors.New("wrong token type"))
			c.Abort()
			return
		}

		id, err := utils.ConvertStringToInt(claims.Subject)
		if err != nil {
			m.logger.Error("convert id error")
			utils.ErrorJSON(c, http.StatusUnauthorized, errors.New("you are not authorized"))
			c.Abort()
			return
		}

		user := new(models.User)
		res := m.db.Where("id = ?", id).First(user)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			utils.ErrorJSON(c, http.StatusUnauthorized, errors.New("user not found"))
			c.Abort()
			return
		}

		c.Set(constants.User, user)
		c.Next()
	}
}
