package services

import (
	"github.com/ssamsara98/go-clean-arch/src/api/dto"
	"github.com/ssamsara98/go-clean-arch/src/infrastructure"
	"github.com/ssamsara98/go-clean-arch/src/lib"
	"github.com/ssamsara98/go-clean-arch/src/models"
	"github.com/ssamsara98/go-clean-arch/src/utils"
)

type PostsService struct {
	logger *lib.Logger
	db     *infrastructure.Database
}

func NewPostsService(
	logger *lib.Logger,
	db *infrastructure.Database,
) *PostsService {
	return &PostsService{
		logger: logger,
		db:     db,
	}
}

func (p PostsService) GetPostList(limit, page *int64) (*[]models.Post, *int64, error) {
	items := new([]models.Post)
	count := new(int64)

	p.db = p.db.SetHandle(p.db.Scopes(utils.Paginate(limit, page)))
	err := p.db.Where(&models.Post{IsPublished: true}).Order("id DESC").Find(items).Offset(-1).Limit(-1).Count(count).Error
	if err != nil {
		return nil, nil, err
	}

	return items, count, nil
}

func (p PostsService) GetPostById(uri *dto.GetPostByIDParams) (*models.Post, error) {
	post := new(models.Post)
	return post, p.db.Preload("Author").First(post, "id = ?", uri.ID).Error
}

func (p PostsService) CreatePost(user *models.User, body *dto.CreatePostDto) (*models.Post, error) {
	post := &models.Post{
		AuthorID:    &user.ID,
		Title:       body.Title,
		Content:     body.Content,
		IsPublished: body.IsPublished,
	}

	err := p.db.Create(post).Error
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (p PostsService) UpdatePost(user *models.User, uri *dto.GetPostByIDParams, body *dto.UpdatePostDto) {
	post := new(models.Post)
	p.db.Where("id = ?", uri.ID).Where("author_id = ?", user.ID).First(post)

	if body.Title != nil {
		post.Title = *body.Title
	}
	if body.Content != nil {
		post.Content = *body.Content
	}

	p.db.Save(post)
}

func (p PostsService) PublishPost(post *models.Post, uri *dto.GetPostByIDParams, body *dto.PublishPostDto) {
	if body.IsPublished != nil {
		post.IsPublished = *body.IsPublished
	}

	p.db.Save(post)
}

func (p PostsService) DeletePost(post *models.Post, user *models.User, uri *dto.GetPostByIDParams) {
	p.db.Where("id = ?", uri.ID).Delete(post)
}
