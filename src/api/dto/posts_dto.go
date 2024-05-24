package dto

type GetPostByIDParams struct {
	ID string `uri:"postId" binding:"required"`
}

type CreatePostDto struct {
	Title       string `form:"title" binding:"required"`
	Content     string `form:"content" binding:"required"`
	IsPublished bool   `form:"isPublished"`
}

type UpdatePostDto struct {
	Title   *string `form:"title"`
	Content *string `form:"content"`
}

type PublishPostDto struct {
	IsPublished *bool `form:"isPublished"`
}

type AddPostCommentDto struct {
	Content string `form:"content"`
}
