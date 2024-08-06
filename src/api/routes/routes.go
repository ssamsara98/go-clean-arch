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
	fx.Provide(NewPostsRoutes),
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
	postsRoutes         *PostsRoutes
}

// NewRoutes sets up routes
func NewRoutes(
	handler *infrastructure.Router,
	rateLimitMiddleware *middlewares.RateLimitMiddleware,
	appRoutes *AppRoutes,
	usersRoutes *UsersRoutes,
	postsRoutes *PostsRoutes,
) *Routes {
	return &Routes{
		handler,
		rateLimitMiddleware,
		appRoutes,
		usersRoutes,
		postsRoutes,
	}
}

// Setup all the route
func (r Routes) Setup() {
	// for _, route := range r {
	// 	route.Setup()
	// }
	r.handler.Use(r.rateLimitMiddleware.Handle())

	root := r.handler.Group("")
	apiV1 := r.handler.Group("v1")

	r.appRoutes.Run(root)
	r.usersRoutes.Run(apiV1)
	r.postsRoutes.Run(apiV1)

	// Not Found route
	r.handler.NoRoute(func(c *gin.Context) {
		utils.ErrorJSON(c, errors.New("not found"), http.StatusNotFound)
	})
}
