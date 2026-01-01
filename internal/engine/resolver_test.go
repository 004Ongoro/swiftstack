package engine

import "testing"

func TestResolveVersion(t *testing.T) {
	tests := []struct {
		base     string
		slice    string
		expected string
	}{
		{"18.2.0", "19.0.0", "19.0.0"},
		{"^16.0.0", "15.0.0", "^16.0.0"},
		{"latest", "1.0.0", "1.0.0"},
	}

	for _, tt := range tests {
		result := ResolveVersion(tt.base, tt.slice)
		if result != tt.expected {
			t.Errorf("ResolveVersion(%s, %s) = %s; want %s", tt.base, tt.slice, result, tt.expected)
		}
	}
}