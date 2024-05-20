package controllers

import (
	"go-clean-arch/src/api/dto"
	"go-clean-arch/src/api/services"
	"go-clean-arch/src/lib"
	"go-clean-arch/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
	logger       *lib.Logger
	usersService *services.UsersService
}

func NewUsersController(
	logger *lib.Logger,
	usersService *services.UsersService,
) *UsersController {
	return &UsersController{
		logger,
		usersService,
	}
}

func (u *UsersController) GetUserList(c *gin.Context) {
	users, err := u.usersService.SetPaginationScope(utils.Paginate(c)).GetUserList()
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}

	utils.JSONWithPagination(c, http.StatusOK, users)
}

func (u *UsersController) GetUserByID(c *gin.Context) {
	var uri dto.GetUserByIDParams
	err := c.ShouldBindUri(&uri)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	user, err := u.usersService.GetUserByID(&uri)
	if err != nil {
		utils.ErrorJSON(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
