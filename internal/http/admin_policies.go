package http

import (
	"encoding/json"
	"fibre_rate_limit_service/internal/policies"

	"github.com/gofiber/fiber/v2"
)

// PolicyRuleRequest represents a single policy rule
type PolicyRuleRequest struct {
	Route     string `json:"route"`
	HeaderKey string `json:"header_key"`
	HeaderVal string `json:"header_value"`
}

// AdminPoliciesHandler handles POST /admin/policies
func AdminPoliciesHandler(c *fiber.Ctx, pe *policies.Evaluator) error {
	var rules []PolicyRuleRequest
	if err := json.Unmarshal(c.Body(), &rules); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	for _, r := range rules {
		// Use AddRule if you want multiple rules per route
		pe.AddRule(r.Route, policies.Rule{
			Header: r.HeaderKey,
			Value:  r.HeaderVal,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "policies added/updated successfully",
	})
}
