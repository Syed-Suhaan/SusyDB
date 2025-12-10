package server

import (
	"strings"
)

// parseCommand parses a raw command string into parts, respecting quoted strings.
// Supports:
// - Basic quoting: SET key "value with spaces"
// - Escape sequences: SET key "value with \"quotes\" inside"
// - Backslash escape: \\ for literal backslash, \" for literal quote
func parseCommand(input string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	escaped := false

	for i := 0; i < len(input); i++ {
		char := input[i]

		// Handle escape sequences
		if escaped {
			// Write the escaped character literally
			current.WriteByte(char)
			escaped = false
			continue
		}

		if char == '\\' {
			escaped = true
			continue
		}

		if char == '"' {
			inQuotes = !inQuotes
			continue
		}

		if char == ' ' && !inQuotes {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteByte(char)
	}

	// Don't forget the last token
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}
