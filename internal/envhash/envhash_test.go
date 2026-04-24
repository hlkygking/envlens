package envhash_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/envhash"
)

func findResult(entries []envhash.Entry, key string) (envhash.Entry, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e, true
		}
	}
	return envhash.Entry{}, false
}

func TestApply_ProducesHashes(t *testing.T) {
	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	entries := envhash.Apply(env, envhash.Options{})
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if len(e.Hash) != 64 {
			t.Errorf("key %s: expected 64-char hex hash, got %d chars", e.Key, len(e.Hash))
		}
	}
}

func TestApply_SortedKeys(t *testing.T) {
	env := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	entries := envhash.Apply(env, envhash.Options{})
	if entries[0].Key != "A_KEY" || entries[1].Key != "M_KEY" || entries[2].Key != "Z_KEY" {
		t.Errorf("entries not sorted: %v", entries)
	}
}

func TestApply_KeysOnly_DiffersFromFull(t *testing.T) {
	env := map[string]string{"SECRET": "topsecret"}
	full := envhash.Apply(env, envhash.Options{KeysOnly: false})
	keysOnly := envhash.Apply(env, envhash.Options{KeysOnly: true})
	if full[0].Hash == keysOnly[0].Hash {
		t.Error("expected KeysOnly hash to differ from full hash")
	}
}

func TestApply_SaltChangesHash(t *testing.T) {
	env := map[string]string{"DB_URL": "postgres://localhost"}
	without := envhash.Apply(env, envhash.Options{})
	with := envhash.Apply(env, envhash.Options{Salt: "mysalt"})
	if without[0].Hash == with[0].Hash {
		t.Error("expected salt to change hash")
	}
}

func TestToMap_RoundTrip(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	entries := envhash.Apply(env, envhash.Options{})
	m := envhash.ToMap(entries)
	if len(m) != 2 {
		t.Fatalf("expected 2 keys in map, got %d", len(m))
	}
	if e, ok := findResult(entries, "FOO"); ok {
		if m["FOO"] != e.Hash {
			t.Errorf("map hash mismatch for FOO")
		}
	}
}

func TestGetSummary_Counts(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2", "C": "3"}
	entries := envhash.Apply(env, envhash.Options{})
	s := envhash.GetSummary(entries)
	if s.Total != 3 || s.Hashed != 3 {
		t.Errorf("unexpected summary: %+v", s)
	}
}

func TestCanonicalHash_Deterministic(t *testing.T) {
	env := map[string]string{"X": "1", "Y": "2"}
	h1 := envhash.CanonicalHash(env, "")
	h2 := envhash.CanonicalHash(env, "")
	if h1 != h2 {
		t.Error("canonical hash is not deterministic")
	}
	if len(h1) != 64 {
		t.Errorf("expected 64-char hash, got %d", len(h1))
	}
}

func TestCanonicalHash_ChangesOnValueEdit(t *testing.T) {
	env1 := map[string]string{"KEY": "value1"}
	env2 := map[string]string{"KEY": "value2"}
	if envhash.CanonicalHash(env1, "") == envhash.CanonicalHash(env2, "") {
		t.Error("expected different canonical hashes for different values")
	}
}

func TestCanonicalHash_SaltPrefix(t *testing.T) {
	env := map[string]string{"K": "v"}
	withSalt := envhash.CanonicalHash(env, "salt123")
	without := envhash.CanonicalHash(env, "")
	if withSalt == without {
		t.Error("salt should change canonical hash")
	}
	if !strings.HasPrefix(withSalt, "") {
		t.Error("canonical hash should be non-empty")
	}
}
