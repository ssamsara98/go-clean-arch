package dto

type GetUserByIDParams struct {
	ID string `param:"userId" validate:"required"`
}
