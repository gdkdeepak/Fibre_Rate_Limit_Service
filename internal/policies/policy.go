package policies

// Policy describes the evaluation rule interface.
type Policy interface {
	Name() string
	Evaluate(meta map[string]string) bool
}

// SimpleHeaderPolicy — example implementation
type SimpleHeaderPolicy struct {
	Header string
	Value  string
}

func (p *SimpleHeaderPolicy) Name() string {
	return "simple_header_policy"
}

func (p *SimpleHeaderPolicy) Evaluate(meta map[string]string) bool {
	v, ok := meta[p.Header]
	if !ok {
		return false
	}
	return v == p.Value
}
