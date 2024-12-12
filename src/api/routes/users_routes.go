package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ssamsara98/go-clean-arch/src/api/controllers"
	"github.com/ssamsara98/go-clean-arch/src/api/middlewares"
	"github.com/ssamsara98/go-clean-arch/src/lib"
)

type UsersRoutes struct {
	logger               *lib.Logger
	paginationMiddleware *middlewares.PaginationMiddleware
	usersController      *controllers.UsersController
}

func NewUsersRoutes(
	logger *lib.Logger,
	paginationMiddleware *middlewares.PaginationMiddleware,
	usersController *controllers.UsersController,
) *UsersRoutes {
	return &UsersRoutes{
		logger,
		paginationMiddleware,
		usersController,
	}
}

func (u UsersRoutes) Run(handler fiber.Router) {
	router := handler.Group("users")

	router.Get("", u.paginationMiddleware.Handle(), u.usersController.GetUserList)
	router.Get("cursor", u.paginationMiddleware.HandleCursor(), u.usersController.GetUserListCursor)
	router.Get("u/:userId", u.usersController.GetUserByID)
}
