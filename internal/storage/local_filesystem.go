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

// NewLocalFileSystemWithDir creates a local storage system at the given directory.
func NewLocalFileSystemWithDir(notesDir string) (*LocalFileSystem, error) {
	if err := os.MkdirAll(notesDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create notes directory %s: %w", notesDir, err)
	}
	return &LocalFileSystem{notesDir: notesDir}, nil
}

// NewLocalFileSystem creates an instance of the local storage system
// at ~/.leaf/notes/.
func NewLocalFileSystem() (*LocalFileSystem, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not determine home directory: %w", err)
	}
	notesDir := filepath.Join(homeDir, ".leaf", "notes")
	return NewLocalFileSystemWithDir(notesDir)
}

func (fs *LocalFileSystem) ListNotes(ctx context.Context) ([]*Note, error) {
	entries, err := os.ReadDir(fs.notesDir)
	if err != nil {
		return nil, fmt.Errorf("could not read directory: %w", err)
	}

	var notes []*Note

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(fs.notesDir, entry.Name())
		note, err := fs.parseNote(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing %s: %v\n", entry.Name(), err)
			continue
		}
		notes = append(notes, note)
	}

	sort.Slice(notes, func(i, j int) bool {
		return notes[i].UpdatedAt.After(notes[j].UpdatedAt)
	})

	return notes, nil
}

func (fs *LocalFileSystem) SaveNote(ctx context.Context, note *Note) error {
	filePath := filepath.Join(fs.notesDir, note.ID+".md")

	if note.CreatedAt.IsZero() {
		note.CreatedAt = time.Now()
	}
	note.UpdatedAt = time.Now()
	note.FilePath = filePath

	fileContent := formatFrontMatter(note) +
		fmt.Sprintf("# %s\n\n%s", note.Title, note.Content)

	if err := os.WriteFile(filePath, []byte(fileContent), 0644); err != nil {
		return fmt.Errorf("could not write note %s: %w", filePath, err)
	}

	return nil
}

func (fs *LocalFileSystem) GetNote(ctx context.Context, id string) (*Note, error) {
	filePath := filepath.Join(fs.notesDir, id+".md")
	note, err := fs.parseNote(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not load note %s: %w", id, err)
	}
	return note, nil
}

func (fs *LocalFileSystem) DeleteNote(ctx context.Context, id string) error {
	filePath := filepath.Join(fs.notesDir, id+".md")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("note not found: %s", id)
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("could not delete note %s: %w", id, err)
	}

	return nil
}

func (fs *LocalFileSystem) SearchNotes(ctx context.Context, query string) ([]*Note, error) {
	notes, err := fs.ListNotes(ctx)
	if err != nil {
		return nil, err
	}

	var results []*Note
	query = strings.ToLower(query)

	for _, note := range notes {
		if strings.Contains(strings.ToLower(note.Title), query) ||
			strings.Contains(strings.ToLower(note.Content), query) {
			results = append(results, note)
		}
	}

	return results, nil
}

func (fs *LocalFileSystem) parseNote(filePath string) (*Note, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	raw := string(fileBytes)
	meta, body, hasMeta := parseFrontMatter(raw)

	var title string
	lines := strings.Split(body, "\n")
	contentStart := len(lines)
	for i, line := range lines {
		if strings.HasPrefix(line, "# ") {
			title = strings.TrimPrefix(line, "# ")
			contentStart = i + 1
			break
		}
	}
	content := strings.TrimSpace(strings.Join(lines[contentStart:], "\n"))

	fileName := filepath.Base(filePath)
	id := strings.TrimSuffix(fileName, ".md")
	if hasMeta && meta.ID != "" {
		id = meta.ID
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	modTime := fileInfo.ModTime()
	createdAt := modTime
	updatedAt := modTime
	if hasMeta {
		if !meta.CreatedAt.IsZero() {
			createdAt = meta.CreatedAt
		}
		if !meta.UpdatedAt.IsZero() {
			updatedAt = meta.UpdatedAt
		}
	}

	return &Note{
		ID:        id,
		Title:     title,
		Content:   content,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		FilePath:  filePath,
	}, nil
}
