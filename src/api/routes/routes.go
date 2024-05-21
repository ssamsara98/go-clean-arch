package routes

import (
	"errors"
	"go-clean-arch/src/api/middlewares"
	"go-clean-arch/src/infrastructure"
	"go-clean-arch/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(NewRoutes),
	fx.Provide(NewAppRoutes),
	fx.Provide(NewUsersRoutes),
)

// Route interface
type IRoute interface {
	Setup()
}

// Routes contains multiple routes
type Routes struct {
	handler             *infrastructure.Router
	rateLimitMiddleware *middlewares.RateLimitMiddleware
	appRoutes           *AppRoutes
	usersRoutes         *UsersRoutes
}

// NewRoutes sets up routes
func NewRoutes(
	handler *infrastructure.Router,
	rateLimitMiddleware *middlewares.RateLimitMiddleware,
	appRoutes *AppRoutes,
	usersRoutes *UsersRoutes,
) *Routes {
	return &Routes{
		handler,
		rateLimitMiddleware,
		appRoutes,
		usersRoutes,
	}
}

// Setup all the route
func (r *Routes) Setup() {
	// for _, route := range r {
	// 	route.Setup()
	// }
	r.handler.Use(r.rateLimitMiddleware.Handle())

	root := r.handler.Group("")
	apiV1 := r.handler.Group("v1")

	r.appRoutes.Run(root)
	r.usersRoutes.Run(apiV1)

	// Not Found route
	r.handler.NoRoute(func(c *gin.Context) {
		utils.ErrorJSON(c, http.StatusNotFound, errors.New("not found"))
	})
}
