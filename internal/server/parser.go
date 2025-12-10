package server

import (
	"strings"
)

// parseCommand parses a raw command string into parts, respecting quoted strings.
// This allows values like: SET key "value with spaces"
// Returns a slice of arguments with quotes removed from quoted args.
func parseCommand(input string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false

	for i := 0; i < len(input); i++ {
		char := input[i]

		if char == '"' {
			inQuotes = !inQuotes
			// Don't include the quote character in the output
			continue
		}

		if char == ' ' && !inQuotes {
			// End of current token
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
