package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateBackup(t *testing.T) {
	skills := []SkillEntry{
		{Name: "test-skill", Description: "desc", Source: "/src", Path: "/src"},
	}

	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "test.tar.gz")

	// CreateBackup will try to read SKILL.md from skill.Path, which won't exist.
	// That's fine — it skips missing files.
	err := CreateBackup(skills, target)
	if err != nil {
		t.Fatalf("CreateBackup: %v", err)
	}

	if _, err := os.Stat(target); os.IsNotExist(err) {
		t.Error("archive file was not created")
	}
}

func TestRestoreBackup(t *testing.T) {
	tmpDir := t.TempDir()

	skillDir := filepath.Join(tmpDir, "src-skill")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}

	content := `---
name: restore-skill
description: A skill to restore
---
`
	skillFile := filepath.Join(skillDir, "SKILL.md")
	if err := os.WriteFile(skillFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	skills := []SkillEntry{
		{Name: "restore-skill", Description: "A skill to restore", Source: "/src", Path: skillDir},
	}

	archive := filepath.Join(tmpDir, "test.tar.gz")
	if err := CreateBackup(skills, archive); err != nil {
		t.Fatal(err)
	}

	dest := filepath.Join(tmpDir, "restored")
	restored, err := RestoreBackup(archive, dest)
	if err != nil {
		t.Fatalf("RestoreBackup: %v", err)
	}

	if len(restored) != 1 {
		t.Fatalf("expected 1 restored skill, got %d", len(restored))
	}
	if restored[0] != "restore-skill" {
		t.Errorf("restored skill name = %q, want %q", restored[0], "restore-skill")
	}

	restoredFile := filepath.Join(dest, "restore-skill", "SKILL.md")
	if _, err := os.Stat(restoredFile); os.IsNotExist(err) {
		t.Error("restored SKILL.md does not exist")
	}

	data, err := os.ReadFile(restoredFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != content {
		t.Errorf("restored content mismatch:\ngot: %q\nwant: %q", string(data), content)
	}
}

func TestRestoreBackupErrors(t *testing.T) {
	_, err := RestoreBackup("/nonexistent/path.tar.gz", t.TempDir())
	if err == nil {
		t.Error("expected error for nonexistent archive")
	}
}
