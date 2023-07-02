package commands

import (
	"context"
	"go-clean-arch/infrastructure"
	"go-clean-arch/lib"
	"go-clean-arch/src/middlewares"
	"go-clean-arch/src/routes"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
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
		e infrastructure.Router,
		lc fx.Lifecycle,
	) {
		logger.Info(`+-----------------------+`)
		logger.Info(`| GO CLEAN ARCHITECTURE |`)
		logger.Info(`+-----------------------+`)

		// Using time zone as specified in env file
		loc, _ := time.LoadLocation(env.TimeZone)
		time.Local = loc

		middleware.Setup()
		routes.Setup()

		logger.Info("Running server")
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				// Start server
				go func() {
					var err error
					if env.ServerPort != "" {
						err = e.Start(":" + env.ServerPort)
					} else {
						err = e.Start(":8080")
					}

					if err != nil && err != http.ErrServerClosed {
						e.Logger.Fatal("shutting down the server")
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return e.Shutdown(ctx)
			},
		})
	}
}

func NewServeCommand() *ServeCommand {
	return &ServeCommand{}
}
