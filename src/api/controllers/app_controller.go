package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/ssamsara98/go-clean-arch/src/api/dto"
	"github.com/ssamsara98/go-clean-arch/src/api/services"
	"github.com/ssamsara98/go-clean-arch/src/constants"
	"github.com/ssamsara98/go-clean-arch/src/helpers"
	"github.com/ssamsara98/go-clean-arch/src/lib"
	"github.com/ssamsara98/go-clean-arch/src/models"
	"github.com/ssamsara98/go-clean-arch/src/utils"
	"gorm.io/gorm"
)

type AppController struct {
	logger     *lib.Logger
	appService *services.AppService
}

func NewAppController(
	logger *lib.Logger,
	appService *services.AppService,
) *AppController {
	return &AppController{
		logger,
		appService,
	}
}

func (app AppController) Home(c *fiber.Ctx) error {
	message := app.appService.Home()
	return utils.SuccessJSON(c, message)
}

func (app AppController) Register(c *fiber.Ctx) error {
	body, err := utils.BindBody[dto.RegisterUserDto](c)
	if err != nil {
		return err
	}

	err = app.appService.FindEmailUsername(body)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusConflict)
		return err
	}

	trxHandle, _ := c.Locals(constants.DBTransaction).(*gorm.DB)

	result, err := app.appService.Register(trxHandle, body)
	if err != nil {
		utils.ErrorJSON(c, err)
		return err
	}

	return utils.SuccessJSON(c, result, http.StatusCreated)
}

func (app AppController) Login(c *fiber.Ctx) error {
	body, err := utils.BindBody[dto.LoginUserDto](c)
	if err != nil {
		return err
	}

	token, err := app.appService.Login(body)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusUnauthorized)
		return err
	}

	return utils.SuccessJSON(c, token, http.StatusCreated)
}

func (app AppController) Me(c *fiber.Ctx) error {
	user, _ := utils.GetUser[models.User](c)
	return utils.SuccessJSON(c, user)
}

func (app AppController) UpdateProfile(c *fiber.Ctx) error {
	body, err := utils.BindBody[dto.UpdateProfileDto](c)
	if err != nil {
		return err
	}

	user, _ := utils.GetUser[models.User](c)
	err = app.appService.UpdateProfile(user.ID, body)
	if err != nil {
		utils.ErrorJSON(c, err)
		return err
	}

	return utils.SuccessJSON(c, "success")
}

func (app AppController) TokenCheck(c *fiber.Ctx) error {
	claims, _ := utils.GetUser[helpers.Claims](c)
	return utils.SuccessJSON(c, claims)
}

func (app AppController) TokenRefresh(c *fiber.Ctx) error {
	user, _ := utils.GetUser[models.User](c)
	tokens, err := app.appService.TokenRefresh(user)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusUnauthorized)
		return err
	}
	return utils.SuccessJSON(c, tokens)
}
