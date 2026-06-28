package internal

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var errNoName = errors.New("skill missing name in frontmatter")

func Discover() ([]SkillEntry, error) {
	seen := make(map[string]SkillEntry)

	cwd, err := os.Getwd()
	if err == nil {
		for _, base := range walkUp(cwd) {
			projectDirs, err := findSkillDirs(filepath.Join(base, ".agents"))
			if err != nil {
				continue
			}
			for _, e := range projectDirs {
				key := e.Name
				if _, ok := seen[key]; !ok {
					seen[key] = e
				}
			}
		}
	}

	home, err := os.UserHomeDir()
	if err == nil {
		globalDirs, err := findSkillDirs(filepath.Join(home, ".agents"))
		if err == nil {
			for _, e := range globalDirs {
				key := e.Name
				if _, ok := seen[key]; !ok {
					seen[key] = e
				}
			}
		}

		superDirs, err := findSuperpowersSkills(home)
		if err == nil {
			for _, e := range superDirs {
				key := e.Name
				if _, ok := seen[key]; !ok {
					seen[key] = e
				}
			}
		}
	}

	items := make([]SkillEntry, 0, len(seen))
	for _, e := range seen {
		items = append(items, e)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	return items, nil
}

func walkUp(start string) []string {
	var dirs []string
	dir := start
	for {
		dirs = append(dirs, dir)
		info, err := os.Stat(filepath.Join(dir, ".git"))
		if err == nil && info.IsDir() {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return dirs
}

func findSkillDirs(base string) ([]SkillEntry, error) {
	skillsDir := filepath.Join(base, "skills")
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, err
	}

	var results []SkillEntry
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillDir := filepath.Join(skillsDir, entry.Name())
		skillFile := filepath.Join(skillDir, "SKILL.md")
		absPath, _ := filepath.Abs(skillDir)

		entry, err := parseSkillFile(skillFile)
		if err != nil {
			continue
		}
		entry.Path = absPath
		entry.Source = absPath
		results = append(results, entry)
	}

	return results, nil
}

func findSuperpowersSkills(home string) ([]SkillEntry, error) {
	base := filepath.Join(home, ".cache", "opencode", "packages")
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}

	var results []SkillEntry
	for _, pkg := range entries {
		if !pkg.IsDir() || !strings.HasPrefix(pkg.Name(), "superpowers") {
			continue
		}

		err := filepath.WalkDir(filepath.Join(base, pkg.Name()), func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.Name() == "SKILL.md" && strings.Contains(path, "node_modules/superpowers/skills/") {
				parent := filepath.Dir(path)
				entry, parseErr := parseSkillFile(path)
				if parseErr != nil {
					return nil
				}
				entry.Path = parent
				entry.Source = "superpowers"
				results = append(results, entry)
			}
			return nil
		})
		if err != nil {
			continue
		}
	}

	return results, nil
}

func parseSkillFile(path string) (SkillEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return SkillEntry{}, err
	}

	content := string(data)
	if !strings.HasPrefix(content, "---") {
		return SkillEntry{}, err
	}

	end := strings.Index(content[3:], "---")
	if end == -1 {
		return SkillEntry{}, err
	}

	fm := content[3 : end+3]

	var name, description string
	scanner := bufio.NewScanner(strings.NewReader(fm))
	var multilineKey string
	var multilineStyle string
	var multilineLines []string
	indent := 0

	for scanner.Scan() {
		line := scanner.Text()

		if multilineKey != "" {
			trimmed := strings.TrimLeft(line, " ")
			if len(trimmed) == 0 && len(multilineLines) > 0 {
				multilineLines = append(multilineLines, "")
				continue
			}
			curIndent := len(line) - len(trimmed)
			if curIndent <= indent || trimmed == "" {
				value := joinMultiline(multilineStyle, multilineLines)
				setField(multilineKey, value, &name, &description)
				multilineKey = ""
				multilineLines = nil

				if trimmed == "" {
					continue
				}
			} else {
				multilineLines = append(multilineLines, trimmed)
				continue
			}
		}

		idx := strings.Index(line, ":")
		if idx == -1 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		rawValue := strings.TrimSpace(line[idx+1:])

		if rawValue == "|" || rawValue == "|-" || rawValue == "|+" {
			multilineKey = key
			multilineStyle = rawValue
			multilineLines = nil
			indent = -1
			continue
		}
		if rawValue == ">" || rawValue == ">-" || rawValue == ">+" {
			multilineKey = key
			multilineStyle = rawValue
			multilineLines = nil
			indent = -1
			continue
		}

		setField(key, rawValue, &name, &description)
	}

	if multilineKey != "" && len(multilineLines) > 0 {
		value := joinMultiline(multilineStyle, multilineLines)
		setField(multilineKey, value, &name, &description)
	}

		if name == "" {
			return SkillEntry{}, errNoName
		}

	return SkillEntry{
		Name:        name,
		Description: description,
	}, nil
}

func setField(key, value string, name, description *string) {
	switch key {
	case "name":
		*name = value
	case "description":
		*description = value
	}
}

func joinMultiline(style string, lines []string) string {
	sep := "\n"
	if style == ">" || style == ">-" || style == ">+" {
		sep = " "
	}
	return strings.TrimSpace(strings.Join(lines, sep))
}
