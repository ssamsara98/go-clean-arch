package bootstrap

import (
	"github.com/ssamsara98/go-clean-arch/src/api/controllers"
	"github.com/ssamsara98/go-clean-arch/src/api/middlewares"
	"github.com/ssamsara98/go-clean-arch/src/api/routes"
	"github.com/ssamsara98/go-clean-arch/src/api/services"
	"github.com/ssamsara98/go-clean-arch/src/helpers"
	"github.com/ssamsara98/go-clean-arch/src/infrastructure"
	"github.com/ssamsara98/go-clean-arch/src/lib"
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
