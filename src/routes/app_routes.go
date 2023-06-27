package routes

import (
	"go-clean-arch/infrastructure"
	"go-clean-arch/src/controllers"
	"go-clean-arch/src/middlewares"
)

type AppRoutes struct {
	appController     *controllers.AppController
	jwtAuthMiddleware *middlewares.JWTAuthMiddleware
}

func NewAppRoutes(
	appController *controllers.AppController,
	jwtAuthMiddleware *middlewares.JWTAuthMiddleware,
) *AppRoutes {
	return &AppRoutes{
		appController,
		jwtAuthMiddleware,
	}
}

func (app AppRoutes) Run(handler infrastructure.Router) {
	handler.GET("/", app.appController.Home)
	handler.POST("/register", app.appController.Register)
	handler.POST("/login", app.appController.Login)
	handler.GET("/me", app.jwtAuthMiddleware.Handle(), app.appController.Me)
	handler.PATCH("/me", app.jwtAuthMiddleware.Handle(), app.appController.UpdateProfile)
}
