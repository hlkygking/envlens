package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envlens/envlens/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(p, []byte(content), 0644))
	return p
}

func TestParseFile_Basic(t *testing.T) {
	path := writeTemp(t, "APP_ENV=production\nDB_HOST=localhost\n")
	env, err := parser.ParseFile(path)
	require.NoError(t, err)
	assert.Equal(t, "production", env["APP_ENV"])
	assert.Equal(t, "localhost", env["DB_HOST"])
}

func TestParseFile_IgnoresComments(t *testing.T) {
	path := writeTemp(t, "# comment\nKEY=value\n")
	env, err := parser.ParseFile(path)
	require.NoError(t, err)
	assert.Len(t, env, 1)
	assert.Equal(t, "value", env["KEY"])
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTemp(t, `SECRET="my secret"
TOKEN='abc123'`)
	env, err := parser.ParseFile(path)
	require.NoError(t, err)
	assert.Equal(t, "my secret", env["SECRET"])
	assert.Equal(t, "abc123", env["TOKEN"])
}

func TestParseFile_EmptyValue(t *testing.T) {
	path := writeTemp(t, "EMPTY=\n")
	env, err := parser.ParseFile(path)
	require.NoError(t, err)
	assert.Equal(t, "", env["EMPTY"])
}

func TestParseFile_MissingFile(t *testing.T) {
	_, err := parser.ParseFile("/nonexistent/.env")
	assert.Error(t, err)
}

func TestParseFile_InvalidLine(t *testing.T) {
	path := writeTemp(t, "NOEQUALSIGN\n")
	_, err := parser.ParseFile(path)
	assert.Error(t, err)
}
