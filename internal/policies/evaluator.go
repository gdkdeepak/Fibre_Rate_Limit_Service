package policies

import (
	"github.com/gofiber/fiber/v2"
	"sync"
)

// Result represents the outcome of policy evaluation
type Result struct {
	Allowed bool
	Reason  string
}

// Rule defines a simple header-based rule
type Rule struct {
	Header string
	Value  string
}

// Evaluator stores all rules for routes
type Evaluator struct {
	mu    sync.RWMutex
	rules map[string][]Rule // route -> list of rules
}

// NewEvaluator creates a new evaluator
func NewEvaluator() *Evaluator {
	return &Evaluator{
		rules: make(map[string][]Rule),
	}
}

// AddRule adds a rule for a route
func (e *Evaluator) AddRule(route string, r Rule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules[route] = append(e.rules[route], r)
}

func (e *Evaluator) SetRule(route string, rule Rule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules[route] = []Rule{rule} // ✅ wrap in slice
}

// Evaluate checks the rules for a client request
func (e *Evaluator) Evaluate(clientID string, route string, c *fiber.Ctx) Result {
	e.mu.RLock()
	defer e.mu.RUnlock()

	rules := e.rules[route]
	for _, r := range rules {
		if c.Get(r.Header) != r.Value {
			return Result{
				Allowed: false,
				Reason:  "Header " + r.Header + " must equal " + r.Value,
			}
		}
	}

	return Result{Allowed: true}
}
