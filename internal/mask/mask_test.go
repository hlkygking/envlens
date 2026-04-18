package mask

import "testing"

func TestIsSensitive(t *testing.T) {
	cases := []struct {
		key      string
		want     bool
	}{
		{"DB_PASSWORD", true},
		{"API_KEY", true},
		{"AUTH_TOKEN", true},
		{"DATABASE_URL", false},
		{"PORT", false},
		{"PRIVATE_KEY", true},
		{"APP_NAME", false},
	}
	for _, c := range cases {
		got := IsSensitive(c.key)
		if got != c.want {
			t.Errorf("IsSensitive(%q) = %v, want %v", c.key, got, c.want)
		}
	}
}

func TestMaskValue(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"ab", "**"},
		{"abc", "***"},
		{"secret", "s****t"},
		{"x", "*"},
	}
	for _, c := range cases {
		got := MaskValue(c.input)
		if got != c.want {
			t.Errorf("MaskValue(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestMaskMap(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "hunter2",
		"APP_NAME":    "myapp",
		"API_KEY":     "abc123",
	}
	masked := MaskMap(env)
	if masked["APP_NAME"] != "myapp" {
		t.Errorf("non-sensitive value should be unchanged")
	}
	if masked["DB_PASSWORD"] == "hunter2" {
		t.Errorf("DB_PASSWORD should be masked")
	}
	if masked["API_KEY"] == "abc123" {
		t.Errorf("API_KEY should be masked")
	}
	// original map should be untouched
	if env["DB_PASSWORD"] != "hunter2" {
		t.Errorf("original map should not be modified")
	}
}
