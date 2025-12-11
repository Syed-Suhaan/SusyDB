package server

import (
	"bufio"
	"strings"
	"testing"
)

func TestParseRESP(t *testing.T) {
	// *2\r\n$3\r\nSET\r\n$3\r\nfoo\r\n
	input := "*2\r\n$3\r\nSET\r\n$3\r\nfoo\r\n"
	reader := bufio.NewReader(strings.NewReader(input))

	args, err := ParseRESP(reader)
	if err != nil {
		t.Fatalf("ParseRESP failed: %v", err)
	}

	if len(args) != 2 {
		t.Fatalf("Expected 2 args, got %d", len(args))
	}

	if args[0] != "SET" {
		t.Errorf("Expected arg[0] 'SET', got '%s'", args[0])
	}
	if args[1] != "foo" {
		t.Errorf("Expected arg[1] 'foo', got '%s'", args[1])
	}
}

func TestParseRESP_InvalidHeader(t *testing.T) {
	input := "$3\r\nSET\r\n" // Not an array
	reader := bufio.NewReader(strings.NewReader(input))

	_, err := ParseRESP(reader)
	if err == nil {
		t.Error("Expected error for non-array input, got nil")
	}
}
