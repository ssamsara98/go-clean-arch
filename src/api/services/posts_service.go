package services

import (
	"go-clean-arch/src/api/dto"
	"go-clean-arch/src/infrastructure"
	"go-clean-arch/src/lib"
	"go-clean-arch/src/models"

	"gorm.io/gorm"
)

type PostsService struct {
	logger          *lib.Logger
	db              *infrastructure.Database
	paginationScope *gorm.DB
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

// PaginationScope
func (p *PostsService) SetPaginationScope(scope func(*gorm.DB) *gorm.DB) *PostsService {
	p.paginationScope = p.db.WithTrx(p.db.Scopes(scope)).DB
	return p
}

func (p *PostsService) GetPostList() (*[]models.Post, *int64, error) {
	var items []models.Post
	var count int64

	err := p.db.WithTrx(p.paginationScope).Where(&models.Post{IsPublished: true}).Order("id DESC").Find(&items).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return nil, nil, err
	}

	return &items, &count, nil
}

func (p *PostsService) GetPostById(uri *dto.GetPostByIDParams) (post models.Post, err error) {
	return post, p.db.First(&post, "id = ?", uri.ID).Error
}

func (p *PostsService) CreatePost(user *models.User, body *dto.CreatePostDto) (*models.Post, error) {
	post := models.Post{
		AuthorId:    user.ID,
		Title:       body.Title,
		Content:     body.Content,
		IsPublished: body.IsPublished,
	}

	err := p.db.Create(&post).Error
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (p *PostsService) UpdatePost(user *models.User, uri *dto.GetPostByIDParams, body *dto.UpdatePostDto) {
	var post models.Post
	p.db.Where("id = ?", uri.ID).Where("author_id = ?", user.ID).First(&post)

	if body.Title != nil {
		post.Title = *body.Title
	}
	if body.Content != nil {
		post.Content = *body.Content
	}

	p.db.Save(&post)
}

func (p *PostsService) PublishPost(post *models.Post, uri *dto.GetPostByIDParams, body *dto.PublishPostDto) {
	if body.IsPublished != nil {
		post.IsPublished = *body.IsPublished
	}

	p.db.Save(&post)
}

func (p *PostsService) DeletePost(post *models.Post, user *models.User, uri *dto.GetPostByIDParams) {
	p.db.Where("id = ?", uri.ID).Delete(&post)
}
