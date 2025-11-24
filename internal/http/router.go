package http

import (
	"fibre_rate_limit_service/internal/limiters"
	"fibre_rate_limit_service/internal/policies"
	"fibre_rate_limit_service/internal/storage"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App, lm *limiters.Manager, pe *policies.Evaluator, store *storage.ShardedMap) {
	api := app.Group("/")

	// /check endpoint
	api.Post("/check", func(c *fiber.Ctx) error {
		return CheckHandler(c, lm, pe)
	})

	// Admin endpoints
	admin := api.Group("/admin")
	admin.Post("/limiters", func(c *fiber.Ctx) error {
		return AdminLimitersHandler(c, lm, store) // pass store
	})
	admin.Post("/policies", func(c *fiber.Ctx) error {
		return AdminPoliciesHandler(c, pe)
	})
	admin.Get("/snapshot", func(c *fiber.Ctx) error {
		return SnapshotHandler(c, lm)
	})
}
