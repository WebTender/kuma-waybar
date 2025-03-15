package main

import (
	"reflect"
	"testing"
)

func TestParseDotEnv(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:  "simple key-value",
			input: "FOO=bar",
			expected: map[string]string{
				"FOO": "bar",
			},
		},
		{
			name:  "quoted value",
			input: `FOO="bar"`,
			expected: map[string]string{
				"FOO": "bar",
			},
		},
		{
			name: "ignore comments and empty lines",
			input: `
# This is a comment
FOO=bar

# Another comment
BAZ=qux
`,
			expected: map[string]string{
				"FOO": "bar",
				"BAZ": "qux",
			},
		},
		{
			name: "ignore invalid lines",
			input: `FOO=bar
INVALID_LINE
BAZ=qux`,
			expected: map[string]string{
				"FOO": "bar",
				"BAZ": "qux",
			},
		},
		{
			name:  "trim spaces",
			input: "   FOO   =   bar   ",
			expected: map[string]string{
				"FOO": "bar",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDotEnv(tt.input)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
