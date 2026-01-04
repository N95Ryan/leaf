package storage

import (
	"fmt"
	"time"
)

// Note represents a markdown note with its metadata
type Note struct {
	ID        string
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	FilePath  string
}

// NewNote creates a new note with a generated ID
func NewNote(title, content string) *Note {
	now := time.Now()
	return &Note{
		ID:        generateID(),
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// generateID generates a unique ID for a note
// TODO: Implement ID generation (UUID or timestamp-based)
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
