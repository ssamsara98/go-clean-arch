package bootstrap

import (
	"go-clean-arch/api/controllers"
	"go-clean-arch/api/middlewares"
	"go-clean-arch/api/routes"
	"go-clean-arch/api/services"
	"go-clean-arch/infrastructure"
	"go-clean-arch/lib"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	lib.Module,
	infrastructure.Module,
	services.Module,
	controllers.Module,
	middlewares.Module,
	routes.Module,
)
