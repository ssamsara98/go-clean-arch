package controllers

import (
	"go-clean-arch/constants"
	"go-clean-arch/lib"
	"go-clean-arch/models"
	"go-clean-arch/src/dto"
	"go-clean-arch/src/services"
	"go-clean-arch/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AppController struct {
	logger     lib.Logger
	appService *services.AppService
}

func NewAppController(
	logger lib.Logger,
	appService *services.AppService,
) *AppController {
	return &AppController{
		logger,
		appService,
	}
}

func (app AppController) Home(c echo.Context) error {
	message := app.appService.Home()
	return utils.SuccessJSON(c, http.StatusOK, message)
}

func (app AppController) Register(c echo.Context) error {
	body := new(dto.RegisterUserDto)
	err := (&echo.DefaultBinder{}).BindBody(c, body)
	if err != nil {
		return utils.ErrorJSON(c, http.StatusBadRequest, err)
	}

	err = app.appService.FindEmailUsername(body)
	if err != nil {
		return utils.ErrorJSON(c, http.StatusConflict, err)
	}

	user, err := app.appService.Register(body)
	if err != nil {
		return utils.ErrorJSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, user)
}

func (app AppController) Login(c echo.Context) error {
	body := new(dto.LoginUserDto)
	err := c.Bind(body)
	if err != nil {
		return utils.ErrorJSON(c, http.StatusBadRequest, err)
	}

	token, err := app.appService.Login(body)
	if err != nil {
		return utils.ErrorJSON(c, http.StatusUnauthorized, err)
	}

	return c.JSON(http.StatusCreated, token)
}

func (app AppController) Me(c echo.Context) error {
	user, _ := c.Get(constants.User).(*models.User)

	return c.JSON(http.StatusOK, user)
}

func (app AppController) UpdateProfile(c echo.Context) error {
	body := new(dto.UpdateProfile)
	err := (&echo.DefaultBinder{}).BindBody(c, body)
	if err != nil {
		return utils.ErrorJSON(c, http.StatusBadRequest, err)
	}

	user, _ := c.Get(constants.User).(*models.User)
	err = app.appService.UpdateProfile(user.ID, body)
	if err != nil {
		return utils.ErrorJSON(c, http.StatusInternalServerError, err)
	}

	return utils.SuccessJSON(c, http.StatusOK, "success")
}
