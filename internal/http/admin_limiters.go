package http

import (
	"encoding/json"
	"time"

	"fibre_rate_limit_service/internal/limiters"
	"fibre_rate_limit_service/internal/storage"

	"github.com/gofiber/fiber/v2"
)

// LimiterRequest represents the JSON body for creating/updating a limiter
type LimiterRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // e.g., "token-bucket"
	Capacity    int64  `json:"capacity"`
	RefillRate  int64  `json:"refill_rate"`
	RefillEvery int64  `json:"refill_every"` // seconds
	TTL         int64  `json:"ttl"`          // seconds
}

// AdminLimitersHandler handles POST /admin/limiters
func AdminLimitersHandler(c *fiber.Ctx, lm *limiters.Manager, store *storage.ShardedMap) error {
	var req LimiterRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Currently only support Token-Bucket
	if req.Type != "token-bucket" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unsupported limiter type",
		})
	}

	// Create TokenBucket
	tb := limiters.NewTokenBucket(limiters.TokenBucketConfig{
		Name:        req.Name,
		Capacity:    int(req.Capacity),
		RefillRate:  int(req.RefillRate),
		RefillEvery: time.Duration(req.RefillEvery) * time.Second,
		TTL:         time.Duration(req.TTL) * time.Second,
	}, store)

	// Hot-reload: add or replace limiter
	lm.SetLimiter(req.Name, tb)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "limiter added/updated successfully",
	})
}
