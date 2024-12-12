package dto

type GetPostByIDParams struct {
	ID string `params:"postId" validate:"required,number"`
}

type CreatePostDto struct {
	Title       string `json:"title" validate:"required"`
	Content     string `json:"content" validate:"required"`
	IsPublished bool   `json:"isPublished"`
}

type UpdatePostDto struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

type PublishPostDto struct {
	IsPublished *bool `json:"isPublished"`
}

type AddPostCommentDto struct {
	Content string `json:"content"`
}
