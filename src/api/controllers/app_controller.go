package controllers

import (
	"go-clean-arch/src/api/dto"
	"go-clean-arch/src/api/services"
	"go-clean-arch/src/constants"
	"go-clean-arch/src/helpers"
	"go-clean-arch/src/lib"
	"go-clean-arch/src/models"
	"go-clean-arch/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (app AppController) Home(c *gin.Context) {
	message := app.appService.Home()
	utils.SuccessJSON(c, message)
}

func (app AppController) Register(c *gin.Context) {
	body, err := utils.BindBody[dto.RegisterUserDto](c)
	if err != nil {
		return
	}

	err = app.appService.FindEmailUsername(body)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusConflict)
		return
	}

	trxHandle, _ := c.MustGet(constants.DBTransaction).(*gorm.DB)

	result, err := app.appService.WithTrx(trxHandle).Register(body)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}

	utils.SuccessJSON(c, result, http.StatusCreated)
}

func (app AppController) Login(c *gin.Context) {
	body, err := utils.BindBody[dto.LoginUserDto](c)
	if err != nil {
		return
	}

	token, err := app.appService.Login(body)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusUnauthorized)
		return
	}

	utils.SuccessJSON(c, token, http.StatusCreated)
}

func (app AppController) Me(c *gin.Context) {
	user, _ := c.MustGet(constants.User).(*models.User)
	utils.SuccessJSON(c, user)
}

func (app AppController) UpdateProfile(c *gin.Context) {
	body, err := utils.BindBody[dto.UpdateProfileDto](c)
	if err != nil {
		return
	}

	user, _ := c.MustGet(constants.User).(*models.User)
	err = app.appService.UpdateProfile(user.ID, body)
	if err != nil {
		utils.ErrorJSON(c, err)
		return
	}

	utils.SuccessJSON(c, "success")
}

func (app AppController) TokenCheck(c *gin.Context) {
	claims, _ := c.MustGet(constants.User).(*helpers.Claims)
	utils.SuccessJSON(c, claims)
}

func (app AppController) TokenRefresh(c *gin.Context) {
	user, _ := c.MustGet(constants.User).(*models.User)
	tokens, err := app.appService.TokenRefresh(user)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusUnauthorized)
		return
	}
	utils.SuccessJSON(c, tokens)
}
