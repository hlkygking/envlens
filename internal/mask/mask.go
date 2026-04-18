package mask

import "strings"

// sensitivePatterns mirrors audit package heuristics
var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "APIKEY",
	"PRIVATE", "CREDENTIAL", "AUTH", "PWD",
}

// IsSensitive returns true if the key looks like a secret.
func IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// MaskValue replaces a sensitive value with asterisks, preserving length up to 3 chars.
func MaskValue(value string) string {
	if len(value) == 0 {
		return ""
	}
	if len(value) <= 3 {
		return strings.Repeat("*", len(value))
	}
	return string(value[0]) + strings.Repeat("*", len(value)-2) + string(value[len(value)-1])
}

// MaskMap returns a copy of the map with sensitive values masked.
func MaskMap(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if IsSensitive(k) {
			out[k] = MaskValue(v)
		} else {
			out[k] = v
		}
	}
	return out
}
