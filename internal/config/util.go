package config

// SafeString returns a fallback value if input is empty.
func SafeString(val, fallback string) string {
	if val == "" {
		return fallback
	}
	return val
}
