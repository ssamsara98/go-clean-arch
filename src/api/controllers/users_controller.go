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

func (u UsersController) GetUserList(c *gin.Context) {
	limit, page := utils.GetPaginationQuery(c)
	items, count, err := u.usersService.SetPaginationScope(utils.Paginate(limit, page)).GetUserList()
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}

	resp := utils.CreatePagination(items, count, limit, page)
	utils.SuccessJSON(c, http.StatusOK, resp)
}

func (u UsersController) GetUserByID(c *gin.Context) {
	uri, err := utils.BindUri[dto.GetUserByIDParams](c)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	user, err := u.usersService.GetUserByID(uri)
	if err != nil {
		utils.ErrorJSON(c, http.StatusNotFound, err)
		return
	}

	utils.SuccessJSON(c, http.StatusOK, user)
}
