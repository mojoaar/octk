package cmd

import (
	"testing"
)

func TestDisplayWidth(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"", 0},
		{"hello", 5},
		{"こんにちは", 5},
		{"abc日本語", 6},
		{"█", 1},
		{"╔═╗", 3},
		{"██████╗", 7},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := displayWidth(tt.input)
			if got != tt.want {
				t.Errorf("displayWidth(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}
