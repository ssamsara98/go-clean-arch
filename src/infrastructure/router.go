package infrastructure

import (
	"fmt"
	"go-clean-arch/src/constants"
	"go-clean-arch/src/lib"
	"go-clean-arch/src/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router -> Gin Router
type Router struct {
	*gin.Engine
}

func debugPrintRouteFunc(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	debugPrint := func(format string, values ...any) {
		if gin.IsDebugging() {
			fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] "+format, values...)
		}
	}
	debugPrint("%-6s %-25s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
}

func logFormatter(params gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if params.IsOutputColor() {
		statusColor = params.StatusCodeColor()
		methodColor = params.MethodColor()
		resetColor = params.ResetColor()
	}

	if params.Latency > time.Minute {
		params.Latency = params.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("[GIN] |%s %3d %s| %10v | %15s |%s %-7s %s| %#v | %s",
		statusColor, params.StatusCode, resetColor,
		params.Latency,
		params.ClientIP,
		methodColor, params.Method, resetColor,
		params.Path,
		params.ErrorMessage,
	)
}

// NewRouter : all the routes are defined here
func NewRouter(
	env *lib.Env,
	logger *lib.Logger,
) *Router {

	// if (env.Environment != constants.Local && env.Environment != constants.Development) && env.SentryDSN != "" {
	// 	if err := sentry.Init(sentry.ClientOptions{
	// 		Dsn:         env.SentryDSN,
	// 		Environment: `clean-backend-` + env.Environment,
	// 	}); err != nil {
	// 		logger.Infof("Sentry initialization failed: %v\n", err)
	// 	}
	// }

	if env.Environment == constants.Production {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
	} else {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
	}
	gin.DefaultWriter = logger.GetGinLogger()
	gin.DebugPrintRouteFunc = debugPrintRouteFunc

	httpRouter := gin.New()
	httpRouter.MaxMultipartMemory = env.MaxMultipartMemory

	httpRouter.Use(gin.LoggerWithFormatter(logFormatter))
	httpRouter.Use(gin.Recovery())
	httpRouter.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	// // Attach sentry middleware
	// httpRouter.Use(sentrygin.New(sentrygin.Options{
	// 	Repanic: true,
	// }))

	httpRouter.GET("/health-check", func(c *gin.Context) {
		utils.SuccessJSON(c, "clean architecture 📺 API Up and Running")
	})

	router := &Router{
		httpRouter,
	}
	return router
}
