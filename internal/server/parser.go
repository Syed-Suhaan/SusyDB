package server

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
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

// ParseRESP parses a RESP array from the reader.
// Format: *<count>\r\n$<len>\r\n<content>\r\n...
func ParseRESP(reader *bufio.Reader) ([]string, error) {
	// Read array header: *<count>\r\n
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimSpace(line)
	if len(line) == 0 || line[0] != '*' {
		return nil, fmt.Errorf("invalid RESP array header: %s", line)
	}

	count, err := strconv.Atoi(line[1:])
	if err != nil {
		return nil, fmt.Errorf("invalid RESP array count: %v", err)
	}

	args := make([]string, 0, count)
	for i := 0; i < count; i++ {
		// Read bulk string header: $<length>\r\n
		line, err = reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] != '$' {
			return nil, fmt.Errorf("invalid RESP bulk string header: %s", line)
		}

		length, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, fmt.Errorf("invalid RESP bulk string length: %v", err)
		}

		// Read data (plus \r\n)
		// We read length + 2 bytes
		data := make([]byte, length+2)
		_, err = io.ReadFull(reader, data)
		if err != nil {
			return nil, err
		}

		// Trim \r\n from data
		args = append(args, string(data[:length]))
	}

	return args, nil
}
