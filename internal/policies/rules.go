package policies

import "strings"

// HeaderEqualsPolicy checks if a metadata key equals an expected value.
type HeaderEqualsPolicy struct {
	name  string
	Key   string
	Value string
}

// NewHeaderEqualsPolicy creates a new policy.
func NewHeaderEqualsPolicy(name, key, value string) *HeaderEqualsPolicy {
	return &HeaderEqualsPolicy{
		name:  name,
		Key:   key,
		Value: value,
	}
}

// Name returns the policy name.
func (p *HeaderEqualsPolicy) Name() string {
	return p.name
}

// Evaluate checks if the metadata contains the required key/value.
func (p *HeaderEqualsPolicy) Evaluate(meta map[string]string) bool {
	val, ok := meta[p.Key]
	if !ok {
		return false
	}
	return strings.EqualFold(val, p.Value)
}
