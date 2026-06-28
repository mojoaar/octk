package cmd

import (
	"fmt"

	"github.com/mojoaar/octk/internal"
)

func List() error {
	skills, err := internal.Discover()
	if err != nil {
		return fmt.Errorf("discover skills: %w", err)
	}

	fmt.Printf("Found %d skills:\n\n", len(skills))

	for _, s := range skills {
		fmt.Printf("  %s", s.Name)
		if s.Description != "" {
			fmt.Printf(" — %s", s.Description)
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\nSources:\n")
	sources := collectSources(skills)
	for src, count := range sources {
		fmt.Printf("  %s: %d skills\n", src, count)
	}

	return nil
}

func collectSources(skills []internal.SkillEntry) map[string]int {
	out := make(map[string]int)
	for _, s := range skills {
		src := s.Source
		if src == "" {
			src = "(unknown)"
		}
		out[src]++
	}
	return out
}
