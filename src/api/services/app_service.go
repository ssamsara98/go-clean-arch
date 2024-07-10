package services

import (
	"errors"
	"go-clean-arch/src/api/dto"
	"go-clean-arch/src/constants"
	"go-clean-arch/src/helpers"
	"go-clean-arch/src/infrastructure"
	"go-clean-arch/src/lib"
	"go-clean-arch/src/models"
	"go-clean-arch/src/utils"

	"gorm.io/gorm"
)

type AppService struct {
	env     *lib.Env
	logger  *lib.Logger
	db      *infrastructure.Database
	JWTAuth *helpers.JWTAuth
}

func NewAppService(
	env *lib.Env,
	logger *lib.Logger,
	db *infrastructure.Database,
	JWTAuth *helpers.JWTAuth,
) *AppService {
	return &AppService{
		env,
		logger,
		db,
		JWTAuth,
	}
}

// WithTrx delegates transaction to repository database
func (app AppService) WithTrx(trxHandle *gorm.DB) *AppService {
	app.db = app.db.WithTrx(trxHandle)
	return &app
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

	if user.Email == body.Email {
		_ = result.AddError(errors.New("email already exist"))
	}
	if user.Username == body.Username {
		_ = result.AddError(errors.New("username already exist"))
	}
	return result.Error
}

func (app AppService) Register(body *dto.RegisterUserDto) (*models.User, error) {
	hashedPassword := utils.HashPassword([]byte(body.Password))

	user := models.User{
		Email:    body.Email,
		Username: body.Username,
		Password: string(hashedPassword),
		Name:     body.Name,
	}

	err := app.db.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type Tokens struct {
	TokenType    string `json:"tokenType"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (app AppService) createToken(user *models.User) (*Tokens, error) {
	accessToken, err := app.JWTAuth.CreateToken(user, constants.TokenAccess)
	if err != nil {
		return nil, err
	}
	refreshToken, err := app.JWTAuth.CreateToken(user, constants.TokenRefresh)
	if err != nil {
		return nil, err
	}

	tokens := &Tokens{
		TokenType:    constants.TokenPrefix,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, nil
}

func (app AppService) Login(body *dto.LoginUserDto) (*Tokens, error) {
	user := new(models.User)
	res := app.db.Where("email = ? OR username = ?", body.UserSession, body.UserSession).First(user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("email/username or password is invalid")
	}

	err := utils.CompareHash(user.Password, body.Password)
	if err != nil {
		return nil, errors.New("email/username or password is invalid")
	}

	return app.createToken(user)
}

func (app AppService) UpdateProfile(id uint, body *dto.UpdateProfileDto) error {
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

func (app AppService) TokenRefresh(user *models.User) (*Tokens, error) {
	return app.createToken(user)
}
