package storage

import (
	"context"
)

// FileSystem defines the interface for note storage operations
type FileSystem interface {
	// ListNotes returns the list of all notes
	ListNotes(ctx context.Context) ([]*Note, error)

	// GetNote retrieves a note by its ID
	GetNote(ctx context.Context, id string) (*Note, error)

	// SaveNote saves a note (create or update)
	SaveNote(ctx context.Context, note *Note) error

	// DeleteNote deletes a note by its ID
	DeleteNote(ctx context.Context, id string) error

	// SearchNotes searches notes by title or content
	SearchNotes(ctx context.Context, query string) ([]*Note, error)
}

// TODO: Implement concrete FileSystem using the filesystem
// Default directory will be ~/.leaf/notes/
