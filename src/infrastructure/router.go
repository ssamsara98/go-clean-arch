package infrastructure

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	lm "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/ssamsara98/go-clean-arch/src/lib"
	"github.com/ssamsara98/go-clean-arch/src/utils"
)

const idleTimeout = 5 * time.Second

func limiterConfig() limiter.Config {
	return limiter.Config{
		Max:        100,
		Expiration: 10 * time.Second,
	}
}

func lmConfig(logger io.Writer) lm.Config {
	return lm.Config{
		Output: logger,
		Format: fmt.Sprintf(
			"pid:${%s} | ${%s} | ${%s} | ${%s} | ${%s} | ${%s} | ${%s}",
			lm.TagPid, lm.TagStatus, lm.TagLatency, lm.TagIP, lm.TagMethod, lm.TagPath, lm.TagError,
		),
		DisableColors: true,
	}
}

/*
Router -> Fiber Router
*/

type Router struct {
	*fiber.App
}

/*
NewRouter : all the routes are defined here
*/

func NewRouter(
	env *lib.Env,
	logger *lib.Logger,
) *Router {

	app := fiber.New(fiber.Config{
		IdleTimeout:  idleTimeout,
		ErrorHandler: fiberErrorHandler,
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
		BodyLimit:    5 * 1024 * 1024, // 5 MB
	})

	/* MaxMultipartMemory */

	/* Middlewares */
	app.Use(idempotency.New())
	app.Use(recover.New())
	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(limiter.New(limiterConfig()))
	app.Use(lm.New(lmConfig(logger.GetFiberLogger())))
	app.Use(etag.New())
	app.Use(healthcheck.New())

	router := &Router{app}
	return router
}

func fiberErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}
	return utils.ErrorJSON(c, err, code)
}
