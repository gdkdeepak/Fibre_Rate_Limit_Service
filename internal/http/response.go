package http

import "github.com/gofiber/fiber/v2"

// JSONResponse creates a consistent JSON response wrapper.
func JSONResponse(c *fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"data": data,
	})
}

// JSONError creates a consistent JSON error wrapper.
func JSONError(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": msg,
	})
}
