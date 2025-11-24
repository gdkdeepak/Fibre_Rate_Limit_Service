package http

import (
	"fibre_rate_limit_service/internal/limiters"

	"github.com/gofiber/fiber/v2"
)

// SnapshotHandler returns the internal memory state of limiters
func SnapshotHandler(c *fiber.Ctx, lm *limiters.Manager) error {
	snapshot := make(map[string]interface{})

	for _, name := range lm.ListLimiters() {
		l, _ := lm.GetLimiter(name)

		// Only supporting TokenBucket for snapshot in memory
		if tb, ok := l.(*limiters.TokenBucket); ok {
			tbSnapshot := tb.StoreSnapshot()
			snapshot[name] = tbSnapshot
		}
	}

	return c.JSON(snapshot)
}
