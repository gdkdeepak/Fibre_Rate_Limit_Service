package http

import (
	//"encoding/json"
	"fibre_rate_limit_service/internal/policies"

	"github.com/gofiber/fiber/v2"
)

// PolicyRuleRequest represents a single policy rule
type PolicyRequest struct {
    Route  string `json:"route"`
    Header string `json:"header"` // single header name
    Value  string `json:"value"`  // value that header must match
}


// AdminPoliciesHandler handles POST /admin/policies
func AdminPoliciesHandler(c *fiber.Ctx, pe *policies.Evaluator) error {
    var req PolicyRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }

    // Add rule to Evaluator
    pe.AddRule(req.Route, policies.Rule{
        Header: req.Header,
        Value:  req.Value,
    })

    return c.JSON(fiber.Map{
        "message": "policy added/updated successfully",
    })
}
