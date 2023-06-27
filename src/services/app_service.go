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

func (app AppService) FindEmailUsername(body *dto.RegisterUserDto) error {
	var user models.User

	result := app.db.Where("email = ?", body.Email).Or("username = ?", body.Username).First(&user)
	if result.Error != nil {
		return nil
	}

	if *user.Email == body.Email {
		_ = result.AddError(errors.New("email already exist"))
	}
	if *user.Username == body.Username {
		_ = result.AddError(errors.New("username already exist"))
	}
	return result.Error
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

func (app AppService) UpdateProfile(id uint, body *dto.UpdateProfile) error {
	user := &models.User{
		ModelBase: lib.ModelBase{ID: id},
		Name:      body.Name,
		Birthdate: body.Birthdate,
	}

	err := app.db.Model(&user).Updates(user).Error
	if err != nil {
		return err
	}

	return nil
}
