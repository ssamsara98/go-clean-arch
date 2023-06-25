package dto

type GetUserByIDParams struct {
	ID string `uri:"userId" binding:"required"`
}
