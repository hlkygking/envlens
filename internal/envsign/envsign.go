package envsign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// Entry holds a key, its value, and the computed signature.
type Entry struct {
	Key       string
	Value     string
	Signature string
	Valid     bool
}

// Summary holds counts after verification.
type Summary struct {
	Total   int
	Valid   int
	Invalid int
}

// Sign computes HMAC-SHA256 signatures for each key=value pair using secret.
func Sign(env map[string]string, secret string) []Entry {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		v := env[k]
		sig := computeHMAC(k, v, secret)
		entries = append(entries, Entry{
			Key:       k,
			Value:     v,
			Signature: sig,
			Valid:     true,
		})
	}
	return entries
}

// Verify checks each env entry against the provided signatures map (key -> expected sig).
func Verify(env map[string]string, signatures map[string]string, secret string) []Entry {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		v := env[k]
		computed := computeHMAC(k, v, secret)
		expected, exists := signatures[k]
		valid := exists && hmac.Equal([]byte(computed), []byte(expected))
		entries = append(entries, Entry{
			Key:       k,
			Value:     v,
			Signature: computed,
			Valid:     valid,
		})
	}
	return entries
}

// GetSummary returns counts from a slice of entries.
func GetSummary(entries []Entry) Summary {
	s := Summary{Total: len(entries)}
	for _, e := range entries {
		if e.Valid {
			s.Valid++
		} else {
			s.Invalid++
		}
	}
	return s
}

func computeHMAC(key, value, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(fmt.Sprintf("%s=%s", strings.TrimSpace(key), strings.TrimSpace(value))))
	return hex.EncodeToString(h.Sum(nil))
}
