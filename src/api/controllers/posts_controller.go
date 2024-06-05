package controllers

import (
	"errors"
	"go-clean-arch/src/api/dto"
	"go-clean-arch/src/api/services"
	"go-clean-arch/src/lib"
	"go-clean-arch/src/models"
	"go-clean-arch/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostsController struct {
	logger       *lib.Logger
	postsService *services.PostsService
}

func NewPostsController(
	logger *lib.Logger,
	postsService *services.PostsService,
) *PostsController {
	return &PostsController{
		logger,
		postsService,
	}
}

func (p *PostsController) GetPostList(c *gin.Context) {
	limit, page := utils.GetPaginationQuery(c)
	items, count, err := p.postsService.SetPaginationScope(utils.Paginate(limit, page)).GetPostList()
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}

	resp := utils.CreatePagination(items, count, limit, page)
	utils.SuccessJSON(c, http.StatusOK, resp)
}

func (p *PostsController) GetPostById(c *gin.Context) {
	var uri dto.GetPostByIDParams
	err := c.ShouldBindUri(&uri)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	user, err := p.postsService.GetPostById(&uri)
	if err != nil {
		utils.ErrorJSON(c, http.StatusNotFound, err)
		return
	}

	utils.SuccessJSON(c, http.StatusOK, user)
}

func (p *PostsController) CreatePost(c *gin.Context) {
	body, err := utils.BindBody[dto.CreatePostDto](c)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}
	user, _ := utils.GetUser[models.User](c)
	result, err := p.postsService.CreatePost(user, body)
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessJSON(c, http.StatusOK, result)
}

func (p *PostsController) UpdatePost(c *gin.Context) {
	user, _ := utils.GetUser[models.User](c)

	uri, err := utils.BindUri[dto.GetPostByIDParams](c)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	body, err := utils.BindBody[dto.UpdatePostDto](c)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	post, err := p.postsService.GetPostById(uri)
	if err != nil {
		utils.ErrorJSON(c, http.StatusNotFound, err)
		return
	}
	if post.AuthorID != user.ID {
		utils.ErrorJSON(c, http.StatusForbidden, errors.New("author_id != user.id"))
		return
	}

	p.postsService.UpdatePost(user, uri, body)

	c.JSON(http.StatusNoContent, gin.H{})
}

func (p *PostsController) PublishPost(c *gin.Context) {
	user, _ := utils.GetUser[models.User](c)

	uri, err := utils.BindUri[dto.GetPostByIDParams](c)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	body, err := utils.BindBody[dto.PublishPostDto](c)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	post, err := p.postsService.GetPostById(uri)
	if err != nil {
		utils.ErrorJSON(c, http.StatusNotFound, err)
		return
	}
	if post.AuthorID != user.ID {
		utils.ErrorJSON(c, http.StatusForbidden, errors.New("author_id != user.id"))
		return
	}

	p.postsService.PublishPost(&post, uri, body)

	c.JSON(http.StatusNoContent, gin.H{})
}

func (p *PostsController) DeletePost(c *gin.Context) {
	user, _ := utils.GetUser[models.User](c)

	uri, err := utils.BindUri[dto.GetPostByIDParams](c)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	post, err := p.postsService.GetPostById(uri)
	if err != nil {
		utils.ErrorJSON(c, http.StatusNotFound, err)
		return
	}
	if post.AuthorID != user.ID {
		utils.ErrorJSON(c, http.StatusForbidden, errors.New("author_id != user.id"))
		return
	}

	p.postsService.DeletePost(&post, user, uri)

	c.JSON(http.StatusNoContent, gin.H{})
}
