package services

import (
	"errors"
	"go-clean-arch/infrastructure"
	"go-clean-arch/lib"
	"go-clean-arch/models"
	"go-clean-arch/src/dto"
	"go-clean-arch/utils"

	"gorm.io/gorm"
)

type AppService struct {
	logger        lib.Logger
	db            infrastructure.Database
	jwtAuthHelper *infrastructure.JWTAuthHelper
}

func NewAppService(
	logger lib.Logger,
	db infrastructure.Database,
	jwtAuthHelper *infrastructure.JWTAuthHelper,
) *AppService {
	return &AppService{
		logger,
		db,
		jwtAuthHelper,
	}
}

func (app AppService) Home() string {
	return "Hello, World!"
}

func (app AppService) Register(body *dto.RegisterUserDto) (*models.User, error) {
	hashedPassword := utils.HashPassword([]byte(body.Password))

	user := models.User{
		Email:    &body.Email,
		Username: &body.Username,
		Password: string(hashedPassword),
		Name:     body.Name,
	}

	err := app.db.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (app AppService) Login(body *dto.LoginUserDto) (*infrastructure.Token, error) {
	var user models.User
	res := app.db.Where("email = ? OR username = ?", body.UserSession, body.UserSession).First(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("email/username or password is invalid")
	}

	err := utils.CompareHash([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return nil, errors.New("email/username or password is invalid")
	}

	// create token
	token := app.jwtAuthHelper.CreateToken(user)

	return token, nil
}
