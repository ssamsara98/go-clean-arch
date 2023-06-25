package dto

type RegisterUserDto struct {
	Email    string `form:"email"`
	Username string `form:"username"`
	Password string `form:"password"`
	Name     string `form:"name"`
}

type LoginUserDto struct {
	UserSession string `form:"userSession"`
	Password    string `form:"password"`
}
