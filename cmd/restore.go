package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mojoaar/octk/internal"
)

func Restore(archivePath string) error {
	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		return fmt.Errorf("archive not found: %s", archivePath)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("find home directory: %w", err)
	}

	targetDir := filepath.Join(home, ".agents", "skills")
	fmt.Printf("Restoring skills to: %s\n", targetDir)

	restored, err := internal.RestoreBackup(archivePath, targetDir)
	if err != nil {
		return fmt.Errorf("restore backup: %w", err)
	}

	for _, name := range restored {
		fmt.Printf("  ✓ %s\n", name)
	}

	fmt.Printf("\nRestored %d skills.\n", len(restored))
	return nil
}
