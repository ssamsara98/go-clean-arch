package commands

import (
	"context"
	"go-clean-arch/src/api/middlewares"
	"go-clean-arch/src/api/routes"
	"go-clean-arch/src/infrastructure"
	"go-clean-arch/src/lib"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// ServeCommand test command
type ServeCommand struct {
}

func (s ServeCommand) Short() string {
	return "Serve Application"
}

func (s ServeCommand) Setup(_ *cobra.Command) {}

func (s ServeCommand) Run() lib.CommandRunner {
	return func(
		env *lib.Env,
		logger *lib.Logger,
		database *infrastructure.Database,
		router *infrastructure.Router,
		middleware *middlewares.Middlewares,
		routes *routes.Routes,
		lc fx.Lifecycle,
	) {
		if env.Environment == "production" {
			logger.Info(`+-------PRODUCTION-------+`)
		}
		logger.Info(`+------------------------+`)
		logger.Info(`| GO CLEAN ARCHITECTURE  |`)
		logger.Info(`+------------------------+`)

		// Using time zone as specified in env file
		loc, _ := time.LoadLocation(env.TimeZone)
		time.Local = loc

		middleware.Setup()
		routes.Setup()

		// if (env.Environment != "local" && env.Environment != "development") && env.SentryDSN != "" {
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

		// /* Using router.Run */
		// if env.Port == "" {
		// 	if err := router.Run(); err != nil {
		// 		logger.Panic(err)
		// 		return
		// 	}
		// } else {
		// 	if err := router.Run(":" + env.Port); err != nil {
		// 		logger.Panic(err)
		// 		return
		// 	}
		// }

		/* Using Lifecycle */
		address := ":8080"
		if env.Port != "" {
			address = ":" + env.Port
		}
		server := &http.Server{Addr: address, Handler: router}
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) (err error) {
				logger.Info("Starting HTTP server at", server.Addr)
				go func() {
					err = server.ListenAndServe()
					if err != nil && err != http.ErrServerClosed {
						logger.Panic(err)
					}
				}()
				return err
			},
			OnStop: func(ctx context.Context) (err error) {
				logger.Info("Gracefully shutting down...")
				err = server.Shutdown(ctx)
				logger.Info("Server was successful shutdown.")
				return err
			},
		})

	}
}

func NewServeCommand() *ServeCommand {
	return &ServeCommand{}
}
