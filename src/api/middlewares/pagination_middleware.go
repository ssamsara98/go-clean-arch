package middlewares

import (
	"go-clean-arch/src/constants"
	"go-clean-arch/src/lib"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginationMiddleware struct {
	logger *lib.Logger
}

func NewPaginationMiddleware(logger *lib.Logger) *PaginationMiddleware {
	return &PaginationMiddleware{logger: logger}
}

func (p *PaginationMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		p.logger.Debug("setting up pagination middleware")

		limit, err := strconv.ParseInt(c.Query("limit"), 10, 0)
		if err != nil {
			limit = 10
		}

		page, err := strconv.ParseInt(c.Query("page"), 10, 0)
		if err != nil {
			page = 1
		}

		c.Set(constants.Limit, limit)
		c.Set(constants.Page, page)

		c.Next()
	}
}
