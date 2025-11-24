package middleware

import "github.com/gofiber/fiber/v2"

func RequestContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// add metadata if needed
		return c.Next()
	}
}
