package utils

import (
	"github.com/gofiber/fiber/v2"
)

func ErrorJSON(c *fiber.Ctx, err error, opts ...any) error {
	statusCode := fiber.StatusInternalServerError
	if len(opts) > 0 {
		statusCode, _ = opts[0].(int)
	}
	resp := fiber.Map{
		"status":     "error",
		"statusCode": statusCode,
		"message":    err.Error(),
		"error":      err,
	}
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Status(statusCode).JSON(resp)
}

func SuccessJSON(c *fiber.Ctx, data any, opts ...any) error {
	statusCode := fiber.StatusOK
	if len(opts) > 0 {
		statusCode, _ = opts[0].(int)
	}
	resp := fiber.Map{
		"status":     "success",
		"statusCode": statusCode,
		"result":     data,
	}
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Status(statusCode).JSON(resp)
}
