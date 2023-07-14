package infrastructure

import (
	"errors"
	"fmt"
	"go-clean-arch/lib"
	"go-clean-arch/models"
	"time"

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
	// db     Database
}

type Claims struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewJWTAuthHelper(
	env *lib.Env,
	logger lib.Logger,
	// db Database,
) *JWTAuthHelper {
	return &JWTAuthHelper{
		env,
		logger,
		// db,
	}
}

// CreateToken creates jwt auth token
func (j JWTAuthHelper) CreateToken(user *models.User) (*Token, error) {
	iat := time.Now()
	exp := iat.Add(30 * 24 * time.Hour)
	claims := &Claims{
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			IssuedAt:  jwt.NewNumericDate(iat),
			ExpiresAt: jwt.NewNumericDate(exp),
			Subject:   fmt.Sprint(user.ID),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.env.JWTSecret))
	if err != nil {
		j.logger.Error("jwt validation failed: ", err)
		return nil, fmt.Errorf("jwt validation failed: %s", err)
	}

	return &Token{
		Type:  "Bearer",
		Token: tokenString,
	}, nil
}

// Authorize authorizes the generated token
func (j JWTAuthHelper) VerifyToken(tokenString string) (*Claims, error) {
	claims := new(Claims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.env.JWTSecret), nil
	})

	if token.Valid {
		return claims, nil
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
