package middlewares

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ssamsara98/go-clean-arch/src/constants"
	"github.com/ssamsara98/go-clean-arch/src/lib"
)

type PaginationMiddleware struct {
	logger *lib.Logger
}

func NewPaginationMiddleware(logger *lib.Logger) *PaginationMiddleware {
	return &PaginationMiddleware{logger: logger}
}

func (p PaginationMiddleware) Handle() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		p.logger.Debug("setting up pagination middleware")

		limit, err := strconv.ParseInt(c.Query("limit"), 10, 0)
		if err != nil || limit < 5 {
			limit = 10
		}

		page, err := strconv.ParseInt(c.Query("page"), 10, 0)
		if err != nil || page < 1 {
			page = 1
		}

		c.Locals(constants.Limit, limit)
		c.Locals(constants.Page, page)

		return c.Next()
	}
}

func (p PaginationMiddleware) HandleCursor() fiber.Handler {
	p.logger.Debug("setting up cursor pagination middleware")

	return func(c *fiber.Ctx) (err error) {
		limit, err := strconv.ParseInt(c.Query("limit"), 10, 0)
		if err != nil || limit < 5 {
			limit = 10
		}

		cursor, err := strconv.ParseInt(c.Query("cursor"), 10, 0)
		if err != nil {
			cursor = 0
		}

		c.Locals(constants.Limit, limit)
		c.Locals(constants.Cursor, cursor)

		return c.Next()
	}
}
