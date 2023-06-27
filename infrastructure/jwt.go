package infrastructure

import (
	"errors"
	"go-clean-arch/lib"
	"go-clean-arch/models"

	"github.com/golang-jwt/jwt/v4"
)

type Token struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

// JWTAuthHelper service relating to authorization
type JWTAuthHelper struct {
	env    *lib.Env
	logger lib.Logger
	db     Database
}

func NewJWTAuthHelper(
	env *lib.Env,
	logger lib.Logger,
	db Database,
) *JWTAuthHelper {
	return &JWTAuthHelper{
		env,
		logger,
		db,
	}
}

// CreateToken creates jwt auth token
func (j JWTAuthHelper) CreateToken(user models.User) *Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"name":     user.Name,
		"email":    *user.Email,
		"username": *user.Username,
	})

	tokenString, err := token.SignedString([]byte(j.env.JWTSecret))

	if err != nil {
		j.logger.Infof("JWT validation failed: ", err)
	}

	return &Token{
		Type:  "Bearer",
		Token: tokenString,
	}
}

// Authorize authorizes the generated token
func (j JWTAuthHelper) Authorize(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.env.JWTSecret), nil
	})

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errors.New("token claims error")
		}

		var user models.User
		j.db.Where("id = ?", uint(claims["sub"].(float64))).First(&user)

		return &user, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, errors.New("token malformed")
		}
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, errors.New("token expired")
		}
	}

	return nil, errors.New("couldn't handle token")
}
