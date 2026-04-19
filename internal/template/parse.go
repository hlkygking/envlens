package template

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ParseFile reads a template definition file where each line is KEY=template.
// Lines starting with '#' and blank lines are ignored.
func ParseFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("template: open %s: %w", path, err)
	}
	defer f.Close()

	templates := make(map[string]string)
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 1 {
			return nil, fmt.Errorf("template: line %d: invalid format %q", lineNum, line)
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		templates[key] = val
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("template: scan %s: %w", path, err)
	}
	return templates, nil
}
