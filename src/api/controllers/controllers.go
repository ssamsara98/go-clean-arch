package controllers

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewAppController),
	fx.Provide(NewUsersController),
	fx.Provide(NewPostsController),
)
