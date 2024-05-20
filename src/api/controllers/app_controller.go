package controllers

import (
	"errors"
	"go-clean-arch/src/api/dto"
	"go-clean-arch/src/api/services"
	"go-clean-arch/src/constants"
	"go-clean-arch/src/lib"
	"go-clean-arch/src/models"
	"go-clean-arch/src/utils"
	"net/http"
	"strings"

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

func (app *AppController) Home(c *gin.Context) {
	message := app.appService.Home()
	utils.SuccessJSON(c, http.StatusOK, message)
}

func (app *AppController) Register(c *gin.Context) {
	var body dto.RegisterUserDto
	err := c.Bind(&body)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	err = app.appService.FindEmailUsername(&body)
	if err != nil {
		utils.ErrorJSON(c, http.StatusConflict, err)
		return
	}

	trxHandle, _ := c.MustGet(constants.DBTransaction).(*gorm.DB)

	user, err := app.appService.WithTrx(trxHandle).Register(&body)
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (app *AppController) Login(c *gin.Context) {
	var body dto.LoginUserDto
	err := c.Bind(&body)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	token, err := app.appService.Login(&body)
	if err != nil {
		utils.ErrorJSON(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusCreated, token)
}

func (app *AppController) Me(c *gin.Context) {
	user, _ := c.MustGet(constants.User).(*models.User)

	c.JSON(http.StatusOK, user)
}

func (app *AppController) UpdateProfile(c *gin.Context) {
	var body dto.UpdateProfileDto
	err := c.Bind(&body)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	user, _ := c.MustGet(constants.User).(*models.User)
	err = app.appService.UpdateProfile(user.ID, &body)
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessJSON(c, http.StatusOK, "success")
}

func (app *AppController) TokenCheck(c *gin.Context) {
	authorizationHeader := c.Request.Header.Get("Authorization")
	if !strings.Contains(authorizationHeader, constants.TokenPrefix) {
		utils.ErrorJSON(c, http.StatusUnauthorized, errors.New("invalid token"))
		return
	}

	tokenString := strings.Replace(authorizationHeader, constants.TokenPrefix+" ", "", -1)

	claims, err := app.appService.TokenCheck(tokenString)
	if err != nil || claims == nil {
		utils.ErrorJSON(c, http.StatusUnauthorized, err)
		return
	}

	utils.JSON(c, http.StatusOK, claims)
}

func (app *AppController) TokenRenew(c *gin.Context) {
	var body dto.RenewAccessTokenReqDto
	err := c.ShouldBindJSON(&body)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	tokens, err := app.appService.TokenRenew(&body)
	if err != nil {
		utils.ErrorJSON(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}
