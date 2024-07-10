package middlewares

import (
	"errors"
	"go-clean-arch/src/constants"
	"go-clean-arch/src/helpers"
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

func (m JWTAuthMiddleware) Handle(tokenType string, needUser bool) gin.HandlerFunc {
	m.logger.Debug("setting up jwt auth middleware")

	return func(c *gin.Context) {
		authorizationHeader := c.Request.Header.Get("Authorization")
		if authorizationHeader == "" {
			utils.ErrorJSON(c, http.StatusUnauthorized, errors.New("no token"))
			c.Abort()
			return
		} else if !strings.Contains(authorizationHeader, constants.TokenPrefix) {
			utils.ErrorJSON(c, http.StatusUnauthorized, errors.New("invalid token"))
			c.Abort()
			return
		}

		tokenString := strings.Replace(authorizationHeader, constants.TokenPrefix+" ", "", -1)
		claims, err := m.JWTAuth.VerifyToken(tokenString, tokenType)
		if err != nil {
			m.logger.Error("claims error")
			utils.ErrorJSON(c, http.StatusUnauthorized, err)
			c.Abort()
			return
		}
		if (claims.Type != constants.TokenAccess) && (claims.Type != constants.TokenRefresh) {
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

		if needUser {
			user := new(models.User)
			res := m.db.Where("id = ?", id).First(user)
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				utils.ErrorJSON(c, http.StatusUnauthorized, errors.New("user not found"))
				c.Abort()
				return
			}
			c.Set(constants.User, user)
		} else {
			c.Set(constants.User, claims)
		}

		c.Next()
	}
}
