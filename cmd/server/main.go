package main

import (
	"time"

	"fibre_rate_limit_service/internal/http"
	"fibre_rate_limit_service/internal/limiters"
	"fibre_rate_limit_service/internal/policies"
	"fibre_rate_limit_service/internal/storage"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1️⃣ Create Fiber app
	app := fiber.New()

	// 2️⃣ Create sharded storage for token buckets
	store := storage.NewShardedMap(16, 10*time.Second, 5*time.Second)
	defer store.Close()

	// 3️⃣ Create limiter manager
	lm := limiters.NewManager()

	// 4️⃣ Create policy evaluator
	pe := policies.NewEvaluator()

	// 5️⃣ Add a sample policy rule for /check
	pe.AddRule("/check", policies.Rule{
		Header: "X-Secret",
		Value:  "123",
	})

	// 6️⃣ Create a token bucket limiter for /check
	tb := limiters.NewTokenBucket(limiters.TokenBucketConfig{
		Name:        "/check",
		Capacity:    5,
		RefillRate:  1,
		RefillEvery: 2 * time.Second,
		TTL:         30 * time.Second,
	}, store)

	// 7️⃣ Attach limiter to route
	lm.SetLimiter("/check", tb)

	// 8️⃣ Setup routes
	http.SetupRouter(app, lm, pe, store)

	// 9️⃣ Start server
	app.Listen(":8080")
}
