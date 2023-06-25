package commands

import (
	"go-clean-arch/infrastructure"
	"go-clean-arch/lib"
	"go-clean-arch/src/middlewares"
	"go-clean-arch/src/routes"
	"time"

	"github.com/spf13/cobra"
)

// ServeCommand test command
type ServeCommand struct {
}

func (s *ServeCommand) Short() string {
	return "Serve Application"
}

func (s *ServeCommand) Setup(_ *cobra.Command) {}

func (s *ServeCommand) Run() lib.CommandRunner {
	return func(
		env *lib.Env,
		logger lib.Logger,
		database infrastructure.Database,
		middleware middlewares.Middlewares,
		routes routes.Routes,
		router infrastructure.Router,
	) {
		logger.Info(`+-----------------------+`)
		logger.Info(`| GO CLEAN ARCHITECTURE |`)
		logger.Info(`+-----------------------+`)

		// Using time zone as specified in env file
		loc, _ := time.LoadLocation(env.TimeZone)
		time.Local = loc

		middleware.Setup()
		routes.Setup()
		// seeds.Setup()

		// if env.Environment != "local" && env.SentryDSN != "" {
		// 	err := sentry.Init(sentry.ClientOptions{
		// 		Dsn:              env.SentryDSN,
		// 		AttachStacktrace: true,
		// 	})
		// 	if err != nil {
		// 		logger.Error("sentry initialization failed")
		// 		logger.Error(err.Error())
		// 	}
		// }

		logger.Info("Running server")
		if env.ServerPort == "" {
			if err := router.Run(); err != nil {
				logger.Fatal(err)
				return
			}
		} else {
			if err := router.Run(":" + env.ServerPort); err != nil {
				logger.Fatal(err)
				return
			}
		}
	}
}

func NewServeCommand() *ServeCommand {
	return &ServeCommand{}
}
