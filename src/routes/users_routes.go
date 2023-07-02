package routes

import (
	"go-clean-arch/lib"
	"go-clean-arch/src/controllers"
	"go-clean-arch/src/middlewares"

	"github.com/labstack/echo/v4"
)

type UsersRoutes struct {
	logger               lib.Logger
	paginationMiddleware *middlewares.PaginationMiddleware
	usersController      *controllers.UsersController
}

func NewUsersRoutes(
	logger lib.Logger,
	paginationMiddleware *middlewares.PaginationMiddleware,
	usersController *controllers.UsersController,
) *UsersRoutes {
	return &UsersRoutes{
		logger,
		paginationMiddleware,
		usersController,
	}
}

func (u UsersRoutes) Run(handler *echo.Group) {
	router := handler.Group("/users")

	router.GET("/", u.usersController.GetUserList, u.paginationMiddleware.Handle())
	router.GET("/:userId", u.usersController.GetUserByID)
}
