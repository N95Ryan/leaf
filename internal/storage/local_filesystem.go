package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type LocalFileSystem struct {
	notesDir string
}

// NotesDir returns the path to the notes directory
func (fs *LocalFileSystem) NotesDir() string {
	return fs.notesDir
}

// NewLocalFileSystem creates an instance of the local storage system
// It determines the path ~/.leaf/notes/, creates the directory if it doesn't exist
func NewLocalFileSystem() (*LocalFileSystem, error) {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not determine home directory: %w", err)
	}

	// Build the path ~/.leaf/notes/ (cross-platform)
	notesDir := filepath.Join(homeDir, ".leaf", "notes")

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(notesDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create notes directory %s: %w", notesDir, err)
	}

	return &LocalFileSystem{
		notesDir: notesDir,
	}, nil
}

func (fs *LocalFileSystem) ListNotes(ctx context.Context) ([]*Note, error) {
	// Read the notes directory content
	entries, err := os.ReadDir(fs.notesDir)
	if err != nil {
		return nil, fmt.Errorf("could not read directory: %w", err)
	}

	var notes []*Note

	// Filter and parse .md files
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Check that the file is a .md file
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		// Build the full path
		filePath := filepath.Join(fs.notesDir, entry.Name())

		// Parse the file
		note, err := fs.parseNote(filePath)
		if err != nil {
			// Log the error but continue
			fmt.Fprintf(os.Stderr, "error parsing %s: %v\n", entry.Name(), err)
			continue
		}

		notes = append(notes, note)
	}

	// Sort by UpdatedAt (descending - most recent first)
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].UpdatedAt.After(notes[j].UpdatedAt)
	})

	return notes, nil
}

func (fs *LocalFileSystem) SaveNote(ctx context.Context, note *Note) error {
	// Build the path: notesDir/{id}.md
	filePath := filepath.Join(fs.notesDir, note.ID+".md")

	// Update UpdatedAt
	note.UpdatedAt = time.Now()
	note.FilePath = filePath

	// Write the file content
	// Format: # Title\n\nContent
	fileContent := fmt.Sprintf("# %s\n\n%s", note.Title, note.Content)

	if err := os.WriteFile(filePath, []byte(fileContent), 0644); err != nil {
		return fmt.Errorf("could not write note %s: %w", filePath, err)
	}

	return nil
}

func (fs *LocalFileSystem) GetNote(ctx context.Context, id string) (*Note, error) {
	// Build the path
	filePath := filepath.Join(fs.notesDir, id+".md")

	// Load and parse the note
	note, err := fs.parseNote(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not load note %s: %w", id, err)
	}

	return note, nil
}

func (fs *LocalFileSystem) DeleteNote(ctx context.Context, id string) error {
	// Build the path
	filePath := filepath.Join(fs.notesDir, id+".md")

	// Check that the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("note not found: %s", id)
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("could not delete note %s: %w", id, err)
	}

	return nil
}

func (fs *LocalFileSystem) SearchNotes(ctx context.Context, query string) ([]*Note, error) {
	// Get all notes
	notes, err := fs.ListNotes(ctx)
	if err != nil {
		return nil, err
	}

	var results []*Note
	query = strings.ToLower(query)

	// Filter by title and content
	for _, note := range notes {
		if strings.Contains(strings.ToLower(note.Title), query) ||
			strings.Contains(strings.ToLower(note.Content), query) {
			results = append(results, note)
		}
	}

	return results, nil
}

// parseNote parses a markdown file into a Note struct
func (fs *LocalFileSystem) parseNote(filePath string) (*Note, error) {
	// Read the file content
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	content := string(fileBytes)

	// Extract the title (first line if it starts with #)
	var title string
	lines := strings.Split(content, "\n")
	if len(lines) > 0 && strings.HasPrefix(lines[0], "# ") {
		title = strings.TrimPrefix(lines[0], "# ")
		// Remove the title from the content
		content = strings.Join(lines[1:], "\n")
		content = strings.TrimSpace(content)
	}

	// Extract the ID from the filename
	fileName := filepath.Base(filePath)
	id := strings.TrimSuffix(fileName, ".md")

	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return &Note{
		ID:        id,
		Title:     title,
		Content:   content,
		CreatedAt: fileInfo.ModTime(), // TODO: store CreatedAt in metadata
		UpdatedAt: fileInfo.ModTime(),
		FilePath:  filePath,
	}, nil
}
