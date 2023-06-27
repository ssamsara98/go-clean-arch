package bootstrap

import (
	"go-clean-arch/infrastructure"
	"go-clean-arch/lib"
	"go-clean-arch/src/controllers"
	"go-clean-arch/src/middlewares"
	"go-clean-arch/src/routes"
	"go-clean-arch/src/services"

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
