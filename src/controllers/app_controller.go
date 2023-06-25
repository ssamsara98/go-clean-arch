package controllers

import (
	"go-clean-arch/constants"
	"go-clean-arch/lib"
	"go-clean-arch/models"
	"go-clean-arch/src/dto"
	"go-clean-arch/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (app AppController) Home(c *gin.Context) {
	message := app.appService.Home()
	lib.SuccessJSON(c, http.StatusOK, message)
}

func (app AppController) Register(c *gin.Context) {
	var body dto.RegisterUserDto
	err := c.Bind(&body)
	if err != nil {
		lib.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	user, err := app.appService.Register(&body)
	if err != nil {
		lib.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (app AppController) Login(c *gin.Context) {
	var body dto.LoginUserDto
	err := c.Bind(&body)
	if err != nil {
		lib.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	token, err := app.appService.Login(&body)
	if err != nil {
		lib.ErrorJSON(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusCreated, token)
}

func (app AppController) Me(c *gin.Context) {
	user, _ := c.MustGet(constants.User).(*models.User)

	c.JSON(http.StatusOK, user)
}
