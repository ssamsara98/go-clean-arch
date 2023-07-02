package dto

import "time"

type RegisterUserDto struct {
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

type LoginUserDto struct {
	UserSession string `json:"userSession" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type UpdateProfile struct {
	Name      string     `json:"name"`
	Birthdate *time.Time `json:"birthdate" time_format:"2006-01-02"`
}
