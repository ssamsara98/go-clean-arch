package middlewares

import (
	"go-clean-arch/constants"
	"go-clean-arch/lib"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PaginationMiddleware struct {
	logger lib.Logger
}

func NewPaginationMiddleware(logger lib.Logger) *PaginationMiddleware {
	return &PaginationMiddleware{logger: logger}
}

func (p PaginationMiddleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			p.logger.Info("setting up pagination middleware")

			limit, err := strconv.ParseInt(c.QueryParam("limit"), 10, 0)
			if err != nil {
				limit = 10
			}

			page, err := strconv.ParseInt(c.QueryParam("page"), 10, 0)
			if err != nil {
				page = 1
			}

			c.Set(constants.Limit, limit)
			c.Set(constants.Page, page)

			return next(c)
		}
	}
}
