package commands

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/ssamsara98/go-clean-arch/src/api/middlewares"
	"github.com/ssamsara98/go-clean-arch/src/api/routes"
	"github.com/ssamsara98/go-clean-arch/src/constants"
	"github.com/ssamsara98/go-clean-arch/src/infrastructure"
	"github.com/ssamsara98/go-clean-arch/src/lib"
	"go.uber.org/fx"
)

/*
ServeCommand test command
*/

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
		if env.Environment == constants.Production {
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

		// if (env.Environment != constants.Local && env.Environment != constants.Development) && env.SentryDSN != "" {
		// 	err := sentry.Init(sentry.ClientOptions{
		// 		Dsn:              env.SentryDSN,
		// 		AttachStacktrace: true,
		// 	})
		// 	if err != nil {
		// 		logger.Error("sentry initialization failed")
		// 		logger.Error(err.Error())
		// 	}
		// }

		port := ":8080"
		if env.Port != "" {
			port = ":" + env.Port
		}
		logger.Info("Running server")

		// /* Using Graceful Shutdown */
		// go func() {
		// 	go func() {
		// 		if err := router.Listen(port); err != nil {
		// 			logger.Panic(err)
		// 		}
		// 	}()

		// 	quit := make(chan os.Signal, 1)                                    // Create channel to signify a signal being sent
		// 	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel
		// 	<-quit                                                             // This blocks the main thread until an interrupt is received

		// 	logger.Info("Gracefully shutting down...")
		// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		// 	defer cancel()
		// 	if err := router.ShutdownWithContext(ctx); err != nil {
		// 		logger.Panic(err)
		// 	}
		// 	select {
		// 	case d := <-ctx.Done():
		// 		logger.Infof("timeout of 5 seconds. %s", d)
		// 	}
		// 	logger.Info("Fiber was successful shutdown.")
		// }()

		/* Using Lifecycle */
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) (err error) {
				logger.Info("Starting HTTP server at", port)
				go func() {
					err = router.Listen(port)
					if err != nil {
						logger.Panic(err)
					}
				}()
				return err
			},
			OnStop: func(ctx context.Context) (err error) {
				logger.Info("Gracefully shutting down...")
				err = router.ShutdownWithContext(ctx)
				logger.Info("Server was successful shutdown.")
				return err
			},
		})
	}
}

func NewServeCommand() *ServeCommand {
	return &ServeCommand{}
}
