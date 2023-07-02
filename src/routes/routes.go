package routes

import (
	"go-clean-arch/infrastructure"

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
	handler     infrastructure.Router
	appRoutes   *AppRoutes
	usersRoutes *UsersRoutes
}

// NewRoutes sets up routes
func NewRoutes(
	handler infrastructure.Router,
	appRoutes *AppRoutes,
	usersRoutes *UsersRoutes,
) Routes {
	return Routes{
		handler,
		appRoutes,
		usersRoutes,
	}
}

// Setup all the route
func (r Routes) Setup() {
	apiV1 := r.handler.Group("/v1")

	r.appRoutes.Run(r.handler)
	r.usersRoutes.Run(apiV1)

	// Not Found route
	// r.handler.NoRoute(func(c *gin.Context) {
	// 	utils.ErrorJSON(c, http.StatusNotFound, errors.New("not found"))
	// })
}
