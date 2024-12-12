package dto

type GetUserByIDParams struct {
	ID string `params:"userId" validate:"required,number"`
}
