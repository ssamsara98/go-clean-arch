package commands

import (
	"context"
	"go-clean-arch/api/middlewares"
	"go-clean-arch/api/routes"
	"go-clean-arch/infrastructure"
	"go-clean-arch/lib"
	"net"
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
		logger *lib.Logger,
		database *infrastructure.Database,
		middleware *middlewares.Middlewares,
		routes *routes.Routes,
		router *infrastructure.Router,
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
		// --- using router.Run
		// if env.ServerPort == "" {
		// 	if err := router.Run(); err != nil {
		// 		logger.Fatal(err)
		// 		return
		// 	}
		// } else {
		// 	if err := router.Run(":" + env.ServerPort); err != nil {
		// 		logger.Fatal(err)
		// 		return
		// 	}
		// }

		// --- using lifecycle
		var srv *http.Server
		if env.ServerPort != "" {
			srv = &http.Server{Addr: ":" + env.ServerPort, Handler: router}
		} else {
			srv = &http.Server{Handler: router}
		}
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				ln, err := net.Listen("tcp", srv.Addr)
				if err != nil {
					return err
				}
				logger.Info("Starting HTTP server at", srv.Addr)
				go func() {
					err := srv.Serve(ln)
					if err != nil {
						logger.Error(err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return srv.Shutdown(ctx)
			},
		})
	}
}

func NewServeCommand() *ServeCommand {
	return &ServeCommand{}
}
