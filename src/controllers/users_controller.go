package controllers

import (
	"go-clean-arch/lib"
	"go-clean-arch/src/dto"
	"go-clean-arch/src/services"
	"go-clean-arch/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UsersController struct {
	logger       lib.Logger
	usersService *services.UsersService
}

func NewUsersController(
	logger lib.Logger,
	usersService *services.UsersService,
) *UsersController {
	return &UsersController{
		logger,
		usersService,
	}
}

func (u UsersController) GetUserList(c echo.Context) error {
	users, err := u.usersService.SetPaginationScope(utils.Paginate(c)).GetUserList()
	if err != nil {
		return utils.ErrorJSON(c, http.StatusInternalServerError, err)
	}

	return utils.JSONWithPagination(c, http.StatusOK, users)
}

func (u UsersController) GetUserByID(c echo.Context) error {
	var uri dto.GetUserByIDParams
	err := c.Bind(&uri)
	if err != nil {
		return utils.ErrorJSON(c, http.StatusBadRequest, err)
	}

	user, err := u.usersService.GetUserByID(&uri)
	if err != nil {
		return utils.ErrorJSON(c, http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, user)
}
