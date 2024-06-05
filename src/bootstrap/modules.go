package bootstrap

import (
	"go-clean-arch/src/api/controllers"
	"go-clean-arch/src/api/middlewares"
	"go-clean-arch/src/api/routes"
	"go-clean-arch/src/api/services"
	"go-clean-arch/src/helpers"
	"go-clean-arch/src/infrastructure"
	"go-clean-arch/src/lib"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	lib.Module,
	infrastructure.Module,
	helpers.Module,
	services.Module,
	controllers.Module,
	middlewares.Module,
	routes.Module,
)
