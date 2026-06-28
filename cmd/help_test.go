package cmd

import (
	"strings"
	"testing"
)

func TestFrameline(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantLen int // should equal innerWidth + 2 for borders
	}{
		{"empty", "", innerWidth + 2},
		{"short", "hello", innerWidth + 2},
		{"exact size", strings.Repeat("x", innerWidth), innerWidth + 2},
		{"ascii art", "   ██████╗  ██████╗████████╗██╗  ██╗", innerWidth + 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := frameline(tt.content)
			if displayWidth(got) != tt.wantLen {
				t.Errorf("frameline(%q) displayWidth = %d, want %d", tt.content, displayWidth(got), tt.wantLen)
			}
			if !strings.HasPrefix(got, "│") || !strings.HasSuffix(got, "│") {
				t.Errorf("frameline(%q) = %q, missing border chars", tt.content, got)
			}
		})
	}
}

func TestCenter(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		width    int
		wantLeft int
	}{
		{"short", "hello", 20, 7},
		{"exact", "hello", 5, 0},
		{"longer", "hello world", 5, 0},
		{"even pad", "ab", 6, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := center(tt.text, tt.width)
			dw := displayWidth(got)
			textDW := displayWidth(tt.text)
			if dw < textDW {
				t.Errorf("center(%q, %d) = %q, display width (%d) < text width (%d)", tt.text, tt.width, got, dw, textDW)
			}
			leftPad := strings.Repeat(" ", tt.wantLeft)
			if !strings.HasPrefix(got, leftPad) {
				t.Errorf("center(%q, %d) = %q, want prefix %q", tt.text, tt.width, got, leftPad)
			}
		})
	}
}

func TestPrintLogo(t *testing.T) {
	// Should not panic
	PrintLogo()
}

func TestPrintHelp(t *testing.T) {
	PrintHelp()
}
