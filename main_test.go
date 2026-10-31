package main

import "testing"

func TestNormalizeTarget(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
		valid bool
	}{
		{name: "IPv4", input: "8.8.8.8", want: "8.8.8.8", valid: true},
		{name: "hostname", input: "example.com", want: "example.com", valid: true},
		{name: "URL host", input: "https://example.com", want: "example.com", valid: true},
		{name: "path", input: "example.com/path", valid: false},
		{name: "empty input", input: " ", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, valid := normalizeTarget(tt.input)
			if got != tt.want || valid != tt.valid {
				t.Fatalf("normalizeTarget(%q) = (%q, %t), want (%q, %t)", tt.input, got, valid, tt.want, tt.valid)
			}
		})
	}
}
