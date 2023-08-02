package routes

import (
	"go-clean-arch/api/controllers"
	"go-clean-arch/api/middlewares"
	"go-clean-arch/infrastructure"
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

func (app AppRoutes) Run(handler infrastructure.Router) {
	handler.GET("/", app.appController.Home)
	handler.POST("/register", app.dbTransactionMiddleware.Handle(), app.appController.Register)
	handler.POST("/login", app.appController.Login)
	handler.GET("/me", app.jwtAuthMiddleware.Handle(), app.appController.Me)
	handler.PATCH("/me", app.jwtAuthMiddleware.Handle(), app.appController.UpdateProfile)
	handler.GET("/token-check", app.appController.TokenCheck)
	handler.GET("/token-renew", app.appController.TokenRenew)
}
