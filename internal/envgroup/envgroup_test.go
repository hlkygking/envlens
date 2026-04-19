package envgroup

import (
	"testing"
)

func findGroup(t *testing.T, r Result, name string) Group {
	t.Helper()
	for _, g := range r.Groups {
		if g.Name == name {
			return g
		}
	}
	t.Fatalf("group %q not found", name)
	return Group{}
}

func TestApply_Basic(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_NAME": "envlens",
		"LOG_LEVEL": "info",
	}
	rules := []Rule{
		{Name: "database", Pattern: "^DB_"},
		{Name: "app", Pattern: "^APP_"},
	}
	r, err := Apply(env, rules)
	if err != nil {
		t.Fatal(err)
	}
	db := findGroup(t, r, "database")
	if len(db.Keys) != 2 {
		t.Errorf("expected 2 db keys, got %d", len(db.Keys))
	}
	app := findGroup(t, r, "app")
	if len(app.Keys) != 1 {
		t.Errorf("expected 1 app key, got %d", len(app.Keys))
	}
}

func TestApply_Ungrouped(t *testing.T) {
	env := map[string]string{"FOO": "1", "BAR": "2"}
	r, err := Apply(env, []Rule{{Name: "foo", Pattern: "^FOO"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Ungrouped) != 1 || r.Ungrouped[0] != "BAR" {
		t.Errorf("expected BAR ungrouped, got %v", r.Ungrouped)
	}
}

func TestApply_InvalidPattern(t *testing.T) {
	env := map[string]string{"A": "1"}
	_, err := Apply(env, []Rule{{Name: "bad", Pattern: "[invalid"}})
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestSummary_Counts(t *testing.T) {
	r := Result{
		Groups:    []Group{{Name: "db", Keys: []string{"A", "B"}}, {Name: "app", Keys: []string{"C"}}},
		Ungrouped: []string{"X", "Y", "Z"},
	}
	sm := Summary(r)
	if sm["db"] != 2 || sm["app"] != 1 || sm["ungrouped"] != 3 {
		t.Errorf("unexpected summary: %v", sm)
	}
}
