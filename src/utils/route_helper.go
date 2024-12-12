package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ssamsara98/go-clean-arch/src/constants"
)

var validate = validator.New()

func checkValidation[T any](c *fiber.Ctx, obj T) error {
	if err := validate.Struct(obj); err != nil {
		errors := []string{}
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("'%s' = '%v' | [%s]",
				err.Field(), err.Value(), err.Tag()))
		}
		return fiber.NewError(fiber.StatusUnprocessableEntity, strings.Join(errors, "\n"))
	}
	return nil
}

func BindBody[T any](c *fiber.Ctx) (*T, error) {
	body := new(T)
	if err := c.BodyParser(body); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := checkValidation(c, body); err != nil {
		return nil, err
	}

	return body, nil
}

func BindParams[T any](c *fiber.Ctx) (*T, error) {
	params := new(T)
	if err := c.ParamsParser(params); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := checkValidation(c, params); err != nil {
		return nil, err
	}
	return params, nil
}

func GetUser[T any](c *fiber.Ctx) (*T, bool) {
	user, boolean := c.Locals(constants.User).(*T)
	return user, boolean
}
