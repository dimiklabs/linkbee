package util

import "strings"

// SanitizeString trims whitespace and removes null bytes from a string.
func SanitizeString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\x00", "")
	return s
}
