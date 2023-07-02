package infrastructure

import (
	"go-clean-arch/lib"
	"go-clean-arch/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Router -> Echo Router
type Router struct {
	*echo.Echo
}

func NewRouter(
	env *lib.Env,
	logger lib.Logger,
) Router {
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			req := c.Request()
			logger.Infof("method=%s uri=%s status=%d remote_address=%s user_agent=%s", req.Method, v.URI, v.Status, req.RemoteAddr, req.UserAgent())
			return nil
		},
	}))

	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	e.GET("/health-check", func(c echo.Context) error {
		return utils.SuccessJSON(c, http.StatusOK, "clean architecture 📺 API Up and Running")
	})

	return Router{e}
}
