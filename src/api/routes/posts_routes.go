package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ssamsara98/go-clean-arch/src/api/controllers"
	"github.com/ssamsara98/go-clean-arch/src/api/middlewares"
	"github.com/ssamsara98/go-clean-arch/src/constants"
	"github.com/ssamsara98/go-clean-arch/src/lib"
)

type PostsRoutes struct {
	logger               *lib.Logger
	paginationMiddleware *middlewares.PaginationMiddleware
	jwtAuthMiddleware    *middlewares.JWTAuthMiddleware
	postsController      *controllers.PostsController
}

func NewPostsRoutes(
	logger *lib.Logger,
	paginationMiddleware *middlewares.PaginationMiddleware,
	jwtAuthMiddleware *middlewares.JWTAuthMiddleware,
	postsController *controllers.PostsController,
) *PostsRoutes {
	return &PostsRoutes{
		logger,
		paginationMiddleware,
		jwtAuthMiddleware,
		postsController,
	}
}

func (p PostsRoutes) Run(handler fiber.Router) {
	router := handler.Group("posts")

	router.Get("", p.paginationMiddleware.Handle(), p.postsController.GetPostList)
	router.Get("p/:postId", p.postsController.GetPostById)

	router.Use(p.jwtAuthMiddleware.Handle(constants.TokenAccess, true))
	router.Post("", p.postsController.CreatePost)
	router.Patch("p/:postId", p.postsController.UpdatePost)
	router.Patch("p/:postId/publish", p.postsController.PublishPost)
	router.Delete("p/:postId", p.postsController.DeletePost)
}
