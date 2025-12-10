package server

import (
	"reflect"
	"testing"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "simple command",
			input: "SET key value",
			want:  []string{"SET", "key", "value"},
		},
		{
			name:  "quoted value with spaces",
			input: `SET key "value with spaces"`,
			want:  []string{"SET", "key", "value with spaces"},
		},
		{
			name:  "json value",
			input: `SET session:1 "{user:suhaan,role:admin}"`,
			want:  []string{"SET", "session:1", `{user:suhaan,role:admin}`},
		},
		{
			name:  "multiple quoted args",
			input: `HSET "my key" "my field" "my value"`,
			want:  []string{"HSET", "my key", "my field", "my value"},
		},
		{
			name:  "single word",
			input: "PING",
			want:  []string{"PING"},
		},
		{
			name:  "empty input",
			input: "",
			want:  nil,
		},
		{
			name:  "extra spaces",
			input: "GET   key",
			want:  []string{"GET", "key"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCommand(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
