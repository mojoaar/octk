package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSkillFile(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantName string
		wantDesc string
		wantErr  bool
	}{
		{
			name: "simple",
			content: `---
name: agent-docs
description: Audit AGENTS.md
---
`,
			wantName: "agent-docs",
			wantDesc: "Audit AGENTS.md",
		},
		{
			name: "quoted description",
			content: `---
name: "PocketBase SDK"
description: "Use when calling PocketBase"
---
`,
			wantName: `"PocketBase SDK"`,
			wantDesc: `"Use when calling PocketBase"`,
		},
		{
			name: "multiline pipe",
			content: `---
name: learn
description: |
  Search, install, update, and rate AI agent skills from agentskill.sh
  (100,000+ skills). Use when the user asks to find skills.
---
`,
			wantName: "learn",
			wantDesc: "Search, install, update, and rate AI agent skills from agentskill.sh\n(100,000+ skills). Use when the user asks to find skills.",
		},
		{
			name: "multiline pipe with blank lines",
			content: `---
name: multi
description: |
  line one

  line two
---
`,
			wantName: "multi",
			wantDesc: "line one\n\nline two",
		},
		{
			name: "multiline gt",
			content: `---
name: fold
description: >
  line one
  line two
  line three
---
`,
			wantName: "fold",
			wantDesc: "line one line two line three",
		},
		{
			name: "multiline gt with dash",
			content: `---
name: fold2
description: >-
  line one
  line two
---
`,
			wantName: "fold2",
			wantDesc: "line one line two",
		},
		{
			name: "no frontmatter",
			content: `# Just a comment
no yaml here`,
			wantErr: true,
		},
		{
			name: "missing name",
			content: `---
description: no name here
---
`,
			wantErr: true,
		},
		{
			name: "no closing frontmatter",
			content: `---
name: broken
`,
			wantErr: true,
		},
		{
			name: "multiline with YAML metadata block",
			content: `---
name: learn
description: >
  Search, install, update, and rate AI agent skills from agentskill.sh
  (100,000+ skills).
metadata:
  author: agentskill-sh
  version: "3.1"
---
`,
			wantName: "learn",
			wantDesc: "Search, install, update, and rate AI agent skills from agentskill.sh (100,000+ skills).",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := filepath.Join(t.TempDir(), "SKILL.md")
			if err := os.WriteFile(tmp, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			entry, err := parseSkillFile(tmp)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if entry.Name != tt.wantName {
				t.Errorf("name = %q, want %q", entry.Name, tt.wantName)
			}
			if entry.Description != tt.wantDesc {
				t.Errorf("description = %q, want %q", entry.Description, tt.wantDesc)
			}
		})
	}
}

func TestSetField(t *testing.T) {
	var name, description string

	setField("name", "test-name", &name, &description)
	if name != "test-name" {
		t.Errorf("name = %q, want %q", name, "test-name")
	}

	setField("description", "test-desc", &name, &description)
	if description != "test-desc" {
		t.Errorf("description = %q, want %q", description, "test-desc")
	}

	setField("unknown", "ignored", &name, &description)
	if name != "test-name" || description != "test-desc" {
		t.Error("unknown field should not overwrite existing values")
	}
}

func TestJoinMultiline(t *testing.T) {
	tests := []struct {
		name  string
		style string
		lines []string
		want  string
	}{
		{
			name:  "pipe",
			style: "|",
			lines: []string{"line one", "line two"},
			want:  "line one\nline two",
		},
		{
			name:  "pipe dash",
			style: "|-",
			lines: []string{"line one", "line two"},
			want:  "line one\nline two",
		},
		{
			name:  "pipe plus",
			style: "|+",
			lines: []string{"line one", "line two"},
			want:  "line one\nline two",
		},
		{
			name:  "gt",
			style: ">",
			lines: []string{"line one", "line two"},
			want:  "line one line two",
		},
		{
			name:  "gt dash",
			style: ">-",
			lines: []string{"line one", "line two"},
			want:  "line one line two",
		},
		{
			name:  "gt plus",
			style: ">+",
			lines: []string{"line one", "line two"},
			want:  "line one line two",
		},
		{
			name:  "empty",
			style: "|",
			lines: []string{},
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := joinMultiline(tt.style, tt.lines)
			if got != tt.want {
				t.Errorf("joinMultiline(%q, %v) = %q, want %q", tt.style, tt.lines, got, tt.want)
			}
		})
	}
}

func TestFindSkillDirs(t *testing.T) {
	base := t.TempDir()
	skillsDir := filepath.Join(base, "skills")

	skillDir := filepath.Join(skillsDir, "test-skill")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}

	content := `---
name: test-skill
description: A test skill
---
`
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	entries, err := findSkillDirs(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Name != "test-skill" {
		t.Errorf("name = %q, want %q", entries[0].Name, "test-skill")
	}
	if entries[0].Description != "A test skill" {
		t.Errorf("description = %q, want %q", entries[0].Description, "A test skill")
	}
}
