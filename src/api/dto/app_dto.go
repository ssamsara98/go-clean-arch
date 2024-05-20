package dto

import "time"

type RegisterUserDto struct {
	Email    string `form:"email" binding:"required"`
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Name     string `form:"name" binding:"required"`
}

type LoginUserDto struct {
	UserSession string `form:"userSession" binding:"required"`
	Password    string `form:"password" binding:"required"`
}

type UpdateProfileDto struct {
	Name      string     `form:"name"`
	Birthdate *time.Time `form:"birthdate" time_format:"xxxx-xx-xx"`
}

type RenewAccessTokenReqDto struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
