package infrastructure

import (
	"errors"
	"fmt"
	"go-clean-arch/lib"
	"go-clean-arch/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTAuthHelper service relating to authorization
type JWTAuthHelper struct {
	env    *lib.Env
	logger *lib.Logger
	// db     *Database
}

type Claims struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

func NewJWTAuthHelper(
	env *lib.Env,
	logger *lib.Logger,
	// db *Database,
) *JWTAuthHelper {
	return &JWTAuthHelper{
		env,
		logger,
		// db,
	}
}

// CreateToken creates jwt auth token
func (j *JWTAuthHelper) CreateToken(user *models.User, duration time.Duration, tokenType string) (string, error) {
	iat := time.Now()
	exp := iat.Add(duration)
	claims := &Claims{
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			IssuedAt:  jwt.NewNumericDate(iat),
			ExpiresAt: jwt.NewNumericDate(exp),
			Subject:   fmt.Sprint(user.ID),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.env.JWTAccessSecret))
	if err != nil {
		j.logger.Error("jwt validation failed: ", err)
		return "", fmt.Errorf("jwt validation failed: %s", err)
	}

	return tokenString, nil
}

// Authorize authorizes the generated token
func (j *JWTAuthHelper) VerifyToken(tokenString string) (*Claims, error) {
	claims := new(Claims)
	var keyfunc jwt.Keyfunc = func(t *jwt.Token) (interface{}, error) { return []byte(j.env.JWTAccessSecret), nil }
	token, err := jwt.ParseWithClaims(tokenString, claims, keyfunc)

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
