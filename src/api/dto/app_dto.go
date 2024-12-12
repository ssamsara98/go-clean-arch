package dto

import (
	"github.com/ssamsara98/go-clean-arch/src/utils"
)

type RegisterUserDto struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

type LoginUserDto struct {
	UserSession string `json:"userSession" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type UpdateProfileDto struct {
	Name      string            `json:"name"`
	Birthdate *utils.CustomDate `json:"birthdate" time_format:"2006-01-02"`
}

type RenewAccessTokenReqDto struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
