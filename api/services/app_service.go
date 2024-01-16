package services

import (
	"errors"
	"go-clean-arch/api/dto"
	"go-clean-arch/constants"
	"go-clean-arch/infrastructure"
	"go-clean-arch/lib"
	"go-clean-arch/models"
	"go-clean-arch/utils"

	"gorm.io/gorm"
)

type AppService struct {
	env           *lib.Env
	logger        lib.Logger
	db            infrastructure.Database
	jwtAuthHelper *infrastructure.JWTAuthHelper
}

func NewAppService(
	env *lib.Env,
	logger lib.Logger,
	db infrastructure.Database,
	jwtAuthHelper *infrastructure.JWTAuthHelper,
) *AppService {
	return &AppService{
		env,
		logger,
		db,
		jwtAuthHelper,
	}
}

// WithTrx delegates transaction to repository database
func (app AppService) WithTrx(trxHandle *gorm.DB) AppService {
	app.db = app.db.WithTrx(trxHandle)
	return app
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
	Type         string `json:"type"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
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

	// create token
	accessToken, err := app.jwtAuthHelper.CreateToken(user, app.env.AccessTokenDuration, constants.TokenAccess)
	if err != nil {
		return nil, err
	}
	refreshToken, err := app.jwtAuthHelper.CreateToken(user, app.env.RefreshTokenDuration, constants.TokenRefresh)
	if err != nil {
		return nil, err
	}

	tokens := &Tokens{
		Type:         constants.TokenPrefix,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, nil
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

func (app AppService) TokenCheck(accessToken string) (*infrastructure.Claims, error) {
	claims, err := app.jwtAuthHelper.VerifyToken(accessToken)
	if err != nil || claims == nil {
		app.logger.Error("claims error")
		return nil, err
	}
	return claims, nil
}

func (app AppService) TokenRenew(body *dto.RenewAccessTokenReqDto) (*Tokens, error) {
	claims, err := app.jwtAuthHelper.VerifyToken(body.RefreshToken)
	if err != nil || claims == nil {
		app.logger.Error("claims error")
		return nil, err
	}
	if claims.Type != constants.TokenRefresh {
		return nil, errors.New("wrong token type")
	}

	user := new(models.User)
	res := app.db.Where("id = ?", claims.Subject).First(user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("userId not found")
	}

	// create token
	accessToken, err := app.jwtAuthHelper.CreateToken(user, app.env.AccessTokenDuration, constants.TokenAccess)
	if err != nil {
		return nil, err
	}
	refreshToken, err := app.jwtAuthHelper.CreateToken(user, app.env.RefreshTokenDuration, constants.TokenRefresh)
	if err != nil {
		return nil, err
	}

	tokens := &Tokens{
		Type:         constants.TokenPrefix,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, nil
}
