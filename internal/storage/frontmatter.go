package storage

import (
	"fmt"
	"strings"
	"time"
)

type frontMatter struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func parseFrontMatter(content string) (frontMatter, string, bool) {
	if !strings.HasPrefix(content, "---\n") {
		return frontMatter{}, content, false
	}

	rest := content[4:]
	end := strings.Index(rest, "\n---")
	if end == -1 {
		return frontMatter{}, content, false
	}

	var meta frontMatter
	for _, line := range strings.Split(rest[:end], "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "id":
			meta.ID = val
		case "created_at":
			if t, err := time.Parse(time.RFC3339, val); err == nil {
				meta.CreatedAt = t
			}
		case "updated_at":
			if t, err := time.Parse(time.RFC3339, val); err == nil {
				meta.UpdatedAt = t
			}
		}
	}

	body := strings.TrimSpace(rest[end+len("\n---"):])
	return meta, body, true
}

func formatFrontMatter(note *Note) string {
	return fmt.Sprintf(
		"---\nid: %s\ncreated_at: %s\nupdated_at: %s\n---\n\n",
		note.ID,
		note.CreatedAt.UTC().Format(time.RFC3339),
		note.UpdatedAt.UTC().Format(time.RFC3339),
	)
}
