package controllers

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/ssamsara98/go-clean-arch/src/api/dto"
	"github.com/ssamsara98/go-clean-arch/src/api/services"
	"github.com/ssamsara98/go-clean-arch/src/lib"
	"github.com/ssamsara98/go-clean-arch/src/models"
	"github.com/ssamsara98/go-clean-arch/src/utils"
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

func (p PostsController) GetPostList(c *fiber.Ctx) error {
	limit, page := utils.GetPaginationQuery(c)
	items, count, err := p.postsService.GetPostList(limit, page)
	if err != nil {
		utils.ErrorJSON(c, err)
		return err
	}

	resp := utils.CreatePagination(items, count, limit, page)
	return utils.SuccessJSON(c, resp)
}

func (p PostsController) GetPostById(c *fiber.Ctx) error {
	uri, err := utils.BindParams[dto.GetPostByIDParams](c)
	if err != nil {
		return err
	}

	user, err := p.postsService.GetPostById(uri)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusNotFound)
		return err
	}

	return utils.SuccessJSON(c, user)
}

func (p PostsController) CreatePost(c *fiber.Ctx) error {
	body, err := utils.BindBody[dto.CreatePostDto](c)
	if err != nil {
		return err
	}

	user, _ := utils.GetUser[models.User](c)
	result, err := p.postsService.CreatePost(user, body)
	if err != nil {
		utils.ErrorJSON(c, err)
		return err
	}

	return utils.SuccessJSON(c, result)
}

func (p PostsController) UpdatePost(c *fiber.Ctx) error {
	user, _ := utils.GetUser[models.User](c)

	uri, err := utils.BindParams[dto.GetPostByIDParams](c)
	if err != nil {
		return err
	}

	body, err := utils.BindBody[dto.UpdatePostDto](c)
	if err != nil {
		return err
	}

	post, err := p.postsService.GetPostById(uri)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusNotFound)
		return err
	}
	if *post.AuthorID != user.ID {
		utils.ErrorJSON(c, errors.New("author_id != user.id"), http.StatusForbidden)
		return err
	}

	p.postsService.UpdatePost(user, uri, body)

	return utils.SuccessJSON(c, fiber.Map{}, http.StatusNoContent)
}

func (p PostsController) PublishPost(c *fiber.Ctx) error {
	user, _ := utils.GetUser[models.User](c)

	uri, err := utils.BindParams[dto.GetPostByIDParams](c)
	if err != nil {
		return err
	}

	body, err := utils.BindBody[dto.PublishPostDto](c)
	if err != nil {
		return err
	}

	post, err := p.postsService.GetPostById(uri)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusNotFound)
		return err
	}
	if *post.AuthorID != user.ID {
		utils.ErrorJSON(c, errors.New("author_id != user.id"), http.StatusForbidden)
		return err
	}

	p.postsService.PublishPost(post, uri, body)

	return utils.SuccessJSON(c, fiber.Map{}, http.StatusNoContent)
}

func (p PostsController) DeletePost(c *fiber.Ctx) error {
	user, _ := utils.GetUser[models.User](c)

	uri, err := utils.BindParams[dto.GetPostByIDParams](c)
	if err != nil {
		return err
	}

	post, err := p.postsService.GetPostById(uri)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusNotFound)
		return err
	}
	if *post.AuthorID != user.ID {
		utils.ErrorJSON(c, errors.New("author_id != user.id"), http.StatusForbidden)
		return err
	}

	p.postsService.DeletePost(post, user, uri)

	return utils.SuccessJSON(c, http.StatusNoContent, fiber.Map{})
}
