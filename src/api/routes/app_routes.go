package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ssamsara98/go-clean-arch/src/api/controllers"
	"github.com/ssamsara98/go-clean-arch/src/api/middlewares"
	"github.com/ssamsara98/go-clean-arch/src/constants"
)

type AppRoutes struct {
	appController           *controllers.AppController
	jwtAuthMiddleware       *middlewares.JWTAuthMiddleware
	dbTransactionMiddleware *middlewares.DBTransactionMiddleware
}

func NewAppRoutes(
	appController *controllers.AppController,
	jwtAuthMiddleware *middlewares.JWTAuthMiddleware,
	dbTransactionMiddleware *middlewares.DBTransactionMiddleware,
) *AppRoutes {
	return &AppRoutes{
		appController,
		jwtAuthMiddleware,
		dbTransactionMiddleware,
	}
}

func (app AppRoutes) Run(handler fiber.Router) {
	handler.Get("", app.appController.Home)
	handler.Post("register", app.dbTransactionMiddleware.Handle(), app.appController.Register)
	handler.Post("login", app.appController.Login)
	handler.Get("me", app.jwtAuthMiddleware.Handle(constants.TokenAccess, true), app.appController.Me)
	handler.Patch("me", app.jwtAuthMiddleware.Handle(constants.TokenAccess, true), app.appController.UpdateProfile)
	handler.Get("token/check", app.jwtAuthMiddleware.Handle(constants.TokenAccess, false), app.appController.TokenCheck)
	handler.Get("token/refresh", app.jwtAuthMiddleware.Handle(constants.TokenRefresh, true), app.appController.TokenRefresh)
}
