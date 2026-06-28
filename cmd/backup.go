package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mojoaar/octk/internal"
)

func Backup() error {
	skills, err := internal.Discover()
	if err != nil {
		return fmt.Errorf("discover skills: %w", err)
	}

	if len(skills) == 0 {
		fmt.Println("No skills found to back up.")
		return nil
	}

	name := fmt.Sprintf("octk-skills-%s.tar.gz", time.Now().Format("2006-01-02"))
	target, err := filepath.Abs(name)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}

	fmt.Printf("Found %d skills.\n", len(skills))
	fmt.Printf("Creating backup: %s\n", target)

	if err := internal.CreateBackup(skills, target); err != nil {
		return fmt.Errorf("create backup: %w", err)
	}

	fi, err := os.Stat(target)
	if err != nil {
		return fmt.Errorf("stat backup: %w", err)
	}

	fmt.Printf("Backup complete: %s (%d bytes)\n", target, fi.Size())
	return nil
}
