package middlewares

import (
	"errors"
	"go-clean-arch/constants"
	"go-clean-arch/infrastructure"
	"go-clean-arch/lib"
	"go-clean-arch/models"
	"go-clean-arch/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JWTAuthMiddleware struct {
	logger        lib.Logger
	jwtAuthHelper *infrastructure.JWTAuthHelper
	db            infrastructure.Database
}

func NewJWTAuthMiddleware(
	logger lib.Logger,
	jwtHelper *infrastructure.JWTAuthHelper,
	db infrastructure.Database,
) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		logger,
		jwtHelper,
		db,
	}
}

func (m JWTAuthMiddleware) Handle() gin.HandlerFunc {
	m.logger.Debug("Setting up jwt auth middleware")

	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, "Bearer ")

		if len(t) == 2 {
			authToken := t[1]
			claims, err := m.jwtAuthHelper.VerifyToken(authToken)
			if err != nil {
				m.logger.Error("claims error")
				utils.ErrorJSON(c, http.StatusUnauthorized, errors.New("you are not authorized"))
				return
			}

			id, err := utils.ConvertStringToInt(claims.Subject)
			if err != nil {
				m.logger.Error("convert id error")
				utils.ErrorJSON(c, http.StatusUnauthorized, errors.New("you are not authorized"))
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
			return
		}

		utils.ErrorJSON(c, http.StatusUnauthorized, errors.New("you are not authorized"))
		c.Abort()
	}
}
