package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/mojoaar/octk/internal"
)

func TestTruncateRunes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		max   int
		want  string
	}{
		{"shorter than max", "hello", 10, "hello"},
		{"exact", "hello", 5, "hello"},
		{"one over", "hello", 4, "h..."},
		{"much over", "hello world this is long", 10, "hello w..."},
		{"max too small", "hello", 2, "he"},
		{"empty", "", 5, ""},
		{"japanese", "こんにちは世界", 4, "こ..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateRunes(tt.input, tt.max)
			if got != tt.want {
				t.Errorf("truncateRunes(%q, %d) = %q, want %q", tt.input, tt.max, got, tt.want)
			}
		})
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		input string
		width int
		want  string
	}{
		{"a", 3, "a  "},
		{"ab", 2, "ab"},
		{"abc", 2, "abc"},
		{"", 4, "    "},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := padRight(tt.input, tt.width)
			if got != tt.want {
				t.Errorf("padRight(%q, %d) = %q, want %q", tt.input, tt.width, got, tt.want)
			}
		})
	}
}

func TestPadLeft(t *testing.T) {
	tests := []struct {
		input string
		width int
		want  string
	}{
		{"a", 3, "  a"},
		{"ab", 2, "ab"},
		{"abc", 2, "abc"},
		{"", 4, "    "},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := padLeft(tt.input, tt.width)
			if got != tt.want {
				t.Errorf("padLeft(%q, %d) = %q, want %q", tt.input, tt.width, got, tt.want)
			}
		})
	}
}

func TestShortenPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home dir")
	}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "home path",
			input: home + "/.agents/skills/foo",
			want:  "~/.agents/skills/foo",
		},
		{
			name:  "non-home path",
			input: "/tmp/something",
			want:  "/tmp/something",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shortenPath(tt.input)
			if got != tt.want {
				t.Errorf("shortenPath(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestTerminalWidth(t *testing.T) {
	orig := os.Getenv("COLUMNS")
	defer os.Setenv("COLUMNS", orig)

	os.Setenv("COLUMNS", "120")
	if got := terminalWidth(); got != 120 {
		t.Errorf("terminalWidth() = %d, want 120", got)
	}

	os.Setenv("COLUMNS", "invalid")
	if got := terminalWidth(); got != 80 {
		t.Errorf("terminalWidth() = %d, want 80 (fallback)", got)
	}

	if err := os.Unsetenv("COLUMNS"); err != nil {
		t.Fatal(err)
	}
	if got := terminalWidth(); got != 80 {
		t.Errorf("terminalWidth() = %d, want 80 (default)", got)
	}
}

func TestListJSON(t *testing.T) {
	skills := []internal.SkillEntry{
		{Name: "test", Description: "desc", Source: "/src", Path: "/src"},
	}

	err := listJSON(skills)
	if err != nil {
		t.Errorf("listJSON: unexpected error: %v", err)
	}
}

func TestListTable(t *testing.T) {
	skills := []internal.SkillEntry{
		{Name: "test-skill", Description: "a test skill description", Source: "/home/user/.agents/skills/test-skill", Path: "/home/user/.agents/skills/test-skill"},
		{Name: "another-skill", Description: "another description", Source: "/home/user/.agents/skills/another-skill", Path: "/home/user/.agents/skills/another-skill"},
	}

	t.Run("default", func(t *testing.T) {
		err := listTable(skills, false)
		if err != nil {
			t.Errorf("listTable: unexpected error: %v", err)
		}
	})

	t.Run("verbose", func(t *testing.T) {
		err := listTable(skills, true)
		if err != nil {
			t.Errorf("listTable verbose: unexpected error: %v", err)
		}
	})

	t.Run("empty", func(t *testing.T) {
		err := listTable(nil, false)
		if err != nil {
			t.Errorf("listTable empty: unexpected error: %v", err)
		}
	})
}

func TestList(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		err := List(true, false)
		if err != nil {
			t.Errorf("List json: unexpected error: %v", err)
		}
	})

	t.Run("table", func(t *testing.T) {
		err := List(false, false)
		if err != nil {
			t.Errorf("List table: unexpected error: %v", err)
		}
	})

	t.Run("verbose", func(t *testing.T) {
		err := List(false, true)
		if err != nil {
			t.Errorf("List verbose: unexpected error: %v", err)
		}
	})
}

var _ = strings.Repeat // ensure strings import is used
