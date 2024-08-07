package routes

import (
	"go-clean-arch/src/api/controllers"
	"go-clean-arch/src/api/middlewares"
	"go-clean-arch/src/lib"

	"github.com/gin-gonic/gin"
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

func (u UsersRoutes) Run(handler *gin.RouterGroup) {
	router := handler.Group("users")

	router.GET("/", u.paginationMiddleware.Handle(), u.usersController.GetUserList)
	router.GET("/u/:userId", u.usersController.GetUserByID)
}
