package http

import (
	"fibre_rate_limit_service/internal/limiters"
	"fibre_rate_limit_service/internal/policies"

	"github.com/gofiber/fiber/v2"
)

// CheckHandler validates a request against policy and limiter
func CheckHandler(c *fiber.Ctx, lm *limiters.Manager, pe *policies.Evaluator) error {
	// Example: use headers to identify client
	clientID := c.Get("X-Client-ID")
	if clientID == "" {
		clientID = "anonymous"
	}

	// Route name could be extracted from path
	route := c.Path()

	// Step 1: Evaluate policy
	policyResult := pe.Evaluate(clientID, route, c)
	if !policyResult.Allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"allowed": false,
			"reason":  policyResult.Reason,
		})
	}

	// Step 2: Apply limiter
	l, ok := lm.GetLimiter(route)
	if !ok {
		// If no limiter defined, allow by default
		return c.JSON(fiber.Map{
			"allowed": true,
			"reason":  "no limiter configured for this route",
		})
	}

	res := l.Check(clientID)
	if !res.Allowed {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"allowed":   false,
			"remaining": res.Remaining,
			"reset_at":  res.ResetAt,
			"reason":    "rate limit exceeded",
		})
	}

	// Allowed
	return c.JSON(fiber.Map{
		"allowed":   true,
		"remaining": res.Remaining,
		"reset_at":  res.ResetAt,
	})
}
