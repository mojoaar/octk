package internal

import "time"

type SkillEntry struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Source      string `json:"source"`
	SourceURL   string `json:"source_url,omitempty"`
	Path        string `json:"path"`
}

type Manifest struct {
	Version   string       `json:"version"`
	CreatedAt time.Time    `json:"created_at"`
	Skills    []SkillEntry `json:"skills"`
	Total     int          `json:"total"`
}

const ToolVersion = "0.1.0"
