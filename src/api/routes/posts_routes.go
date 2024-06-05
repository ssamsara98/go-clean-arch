package routes

import (
	"go-clean-arch/src/api/controllers"
	"go-clean-arch/src/api/middlewares"
	"go-clean-arch/src/constants"
	"go-clean-arch/src/lib"

	"github.com/gin-gonic/gin"
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

func (p *PostsRoutes) Run(handler *gin.RouterGroup) {
	router := handler.Group("posts")

	router.GET("/", p.paginationMiddleware.Handle(), p.postsController.GetPostList)
	router.GET("/p/:postId", p.postsController.GetPostById)

	router.Use(p.jwtAuthMiddleware.Handle(constants.TokenAccess, true))
	router.POST("/", p.postsController.CreatePost)
	router.PATCH("/p/:postId", p.postsController.UpdatePost)
	router.PATCH("/p/:postId/publish", p.postsController.PublishPost)
	router.DELETE("/p/:postId", p.postsController.DeletePost)
}
