package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mojoaar/octk/internal"
)

func List(jsonOutput, verbose bool) error {
	skills, err := internal.Discover()
	if err != nil {
		return fmt.Errorf("discover skills: %w", err)
	}

	if jsonOutput {
		return listJSON(skills)
	}

	return listTable(skills, verbose)
}

func listJSON(skills []internal.SkillEntry) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(skills)
}

func terminalWidth() int {
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if w, err := strconv.Atoi(cols); err == nil && w > 0 {
			return w
		}
	}
	return 80
}

func shortenPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}
	return path
}

func truncateRunes(s string, max int) string {
	if displayWidth(s) <= max {
		return s
	}
	if max < 3 {
		return s[:max]
	}
	runes := []rune(s)
	target := max - 3
	if target <= 0 {
		return "..."
	}
	return string(runes[:target]) + "..."
}

func padRight(s string, w int) string {
	dw := displayWidth(s)
	if dw >= w {
		return s
	}
	return s + strings.Repeat(" ", w-dw)
}

func padLeft(s string, w int) string {
	dw := displayWidth(s)
	if dw >= w {
		return s
	}
	return strings.Repeat(" ", w-dw) + s
}

func listTable(skills []internal.SkillEntry, verbose bool) error {
	tw := terminalWidth()

	numW := 4
	nameW := 30
	bodyW := tw - numW - 3 - nameW - 3 // │ separators and gaps
	if bodyW < 20 {
		bodyW = 20
	}
	srcW := bodyW * 3 / 10
	descW := bodyW - srcW

	if verbose {
		srcW = bodyW / 2
		descW = bodyW - srcW
	}

	fmt.Printf("Found %d skills:\n\n", len(skills))

	dash := func(w int) string { return strings.Repeat("─", w) }

	fmt.Printf("┌%s┬%s┬%s┬%s┐\n", dash(numW), dash(nameW), dash(srcW), dash(descW))
	fmt.Printf("│%s│%s│%s│%s│\n",
		padRight("#", numW),
		padRight("Name", nameW),
		padRight("Source", srcW),
		padRight("Description", descW))
	fmt.Printf("├%s┼%s┼%s┼%s┤\n", dash(numW), dash(nameW), dash(srcW), dash(descW))

	for i, s := range skills {
		name := truncateRunes(s.Name, nameW)
		desc := truncateRunes(s.Description, descW)
		src := shortenPath(s.Source)
		if verbose {
			name = s.Name
			desc = s.Description
		} else {
			src = truncateRunes(src, srcW)
		}
		fmt.Printf("│%s│%s│%s│%s│\n",
			padLeft(fmt.Sprintf("%d", i+1), numW),
			padRight(name, nameW),
			padRight(src, srcW),
			padRight(desc, descW))
	}

	fmt.Printf("└%s┴%s┴%s┴%s┘\n", dash(numW), dash(nameW), dash(srcW), dash(descW))

	return nil
}
