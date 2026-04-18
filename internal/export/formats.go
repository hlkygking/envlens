package export

// SupportedFormats returns all valid export format strings.
func SupportedFormats() []Format {
	return []Format{FormatDotenv, FormatJSON, FormatShell}
}

// ParseFormat converts a string to a Format, returning an error if invalid.
func ParseFormat(s string) (Format, error) {
	f := Format(s)
	for _, v := range SupportedFormats() {
		if f == v {
			return f, nil
		}
	}
	return "", fmt.Errorf("unknown format %q; supported: dotenv, json, shell")
}
