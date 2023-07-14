package controllers

import (
	"errors"
	"go-clean-arch/constants"
	"go-clean-arch/infrastructure"
	"go-clean-arch/lib"
	"go-clean-arch/models"
	"go-clean-arch/src/dto"
	"go-clean-arch/src/services"
	"go-clean-arch/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppController struct {
	logger        lib.Logger
	appService    *services.AppService
	jwtAuthHelper *infrastructure.JWTAuthHelper
}

func NewAppController(
	logger lib.Logger,
	appService *services.AppService,
	jwtAuthHelper *infrastructure.JWTAuthHelper,
) *AppController {
	return &AppController{
		logger,
		appService,
		jwtAuthHelper,
	}
}

func (app AppController) Home(c *gin.Context) {
	message := app.appService.Home()
	utils.SuccessJSON(c, http.StatusOK, message)
}

func (app AppController) Register(c *gin.Context) {
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

func (app AppController) Login(c *gin.Context) {
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

func (app AppController) Me(c *gin.Context) {
	user, _ := c.MustGet(constants.User).(*models.User)

	c.JSON(http.StatusOK, user)
}

func (app AppController) UpdateProfile(c *gin.Context) {
	var body dto.UpdateProfile
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

func (app AppController) TokenCheck(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	t := strings.Split(authHeader, "Bearer ")

	if len(t) == 2 {
		authToken := t[1]
		claims, err := app.jwtAuthHelper.VerifyToken(authToken)
		if err != nil || claims == nil {
			app.logger.Error("claims error")
			utils.ErrorJSON(c, http.StatusUnauthorized, errors.New("you are not authorized"))
			return
		}

		utils.SuccessJSON(c, http.StatusOK, "success")
		return
	}

	utils.ErrorJSON(c, http.StatusUnauthorized, errors.New("you are not authorized"))
}
