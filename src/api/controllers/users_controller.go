package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/ssamsara98/go-clean-arch/src/api/dto"
	"github.com/ssamsara98/go-clean-arch/src/api/services"
	"github.com/ssamsara98/go-clean-arch/src/lib"
	"github.com/ssamsara98/go-clean-arch/src/utils"
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

func (u UsersController) GetUserList(c *fiber.Ctx) error {
	limit, page := utils.GetPaginationQuery(c)
	items, count, err := u.usersService.GetUserList(limit, page)
	if err != nil {
		utils.ErrorJSON(c, err)
		return err
	}

	resp := utils.CreatePagination(items, count, limit, page)
	return utils.SuccessJSON(c, resp)
}
func (u UsersController) GetUserListCursor(c *fiber.Ctx) error {
	limit, cursor := utils.GetPaginationCursorQuery(c)
	items, err := u.usersService.GetUserListCursor(limit, cursor)
	if err != nil {
		return err
	}

	resp := utils.CreatePaginationCursor(items, limit, cursor)
	return utils.SuccessJSON(c, resp)
}

func (u UsersController) GetUserByID(c *fiber.Ctx) error {
	uri, err := utils.BindParams[dto.GetUserByIDParams](c)
	if err != nil {
		return err
	}

	user, err := u.usersService.GetUserByID(uri)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusNotFound)
		return err
	}

	return utils.SuccessJSON(c, user)
}
