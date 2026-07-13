package storage_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/N95Ryan/leaf/internal/storage"
)

func newTestFS(t *testing.T) *storage.LocalFileSystem {
	t.Helper()
	fs, err := storage.NewLocalFileSystemWithDir(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocalFileSystemWithDir() failed: %v", err)
	}
	return fs
}

func TestNewLocalFileSystemWithDir(t *testing.T) {
	fs := newTestFS(t)

	if fs == nil {
		t.Fatal("NewLocalFileSystemWithDir() returned nil")
	}

	if _, err := os.Stat(fs.NotesDir()); os.IsNotExist(err) {
		t.Fatalf("notes directory was not created: %s", fs.NotesDir())
	}
}

func TestSaveAndGetNote(t *testing.T) {
	fs := newTestFS(t)
	ctx := context.Background()

	note := storage.NewNote("Test Note", "This is the test note content")
	originalID := note.ID
	originalCreatedAt := note.CreatedAt

	if err := fs.SaveNote(ctx, note); err != nil {
		t.Fatalf("SaveNote() failed: %v", err)
	}

	retrievedNote, err := fs.GetNote(ctx, originalID)
	if err != nil {
		t.Fatalf("GetNote() failed: %v", err)
	}

	if retrievedNote.Title != note.Title {
		t.Errorf("title mismatch: expected %q, got %q", note.Title, retrievedNote.Title)
	}
	if retrievedNote.Content != note.Content {
		t.Errorf("content mismatch: expected %q, got %q", note.Content, retrievedNote.Content)
	}
	if retrievedNote.ID != originalID {
		t.Errorf("ID mismatch: expected %q, got %q", originalID, retrievedNote.ID)
	}
	if !retrievedNote.CreatedAt.Equal(originalCreatedAt.Truncate(time.Second)) {
		t.Errorf("CreatedAt not persisted: expected %v, got %v", originalCreatedAt, retrievedNote.CreatedAt)
	}
}

func TestListNotes(t *testing.T) {
	fs := newTestFS(t)
	ctx := context.Background()

	notes := []*storage.Note{
		storage.NewNote("Note 1", "Content 1"),
		storage.NewNote("Note 2", "Content 2"),
		storage.NewNote("Note 3", "Content 3"),
	}

	for i, note := range notes {
		if err := fs.SaveNote(ctx, note); err != nil {
			t.Fatalf("SaveNote() failed: %v", err)
		}
		if i < len(notes)-1 {
			time.Sleep(1100 * time.Millisecond)
		}
	}

	listedNotes, err := fs.ListNotes(ctx)
	if err != nil {
		t.Fatalf("ListNotes() failed: %v", err)
	}

	if len(listedNotes) != len(notes) {
		t.Errorf("note count mismatch: expected %d, got %d", len(notes), len(listedNotes))
	}

	for i := 0; i < len(listedNotes)-1; i++ {
		if !listedNotes[i].UpdatedAt.After(listedNotes[i+1].UpdatedAt) {
			t.Error("notes are not sorted by UpdatedAt in descending order")
		}
	}
}

func TestSearchNotes(t *testing.T) {
	fs := newTestFS(t)
	ctx := context.Background()

	note1 := storage.NewNote("Go Tutorial", "Learn Go and concurrency")
	note2 := storage.NewNote("Python Tips", "Useful tips for Python")
	note3 := storage.NewNote("JavaScript Guide", "Complete Go guide for beginners")

	for _, note := range []*storage.Note{note1, note2, note3} {
		if err := fs.SaveNote(ctx, note); err != nil {
			t.Fatalf("SaveNote() failed: %v", err)
		}
	}

	results, err := fs.SearchNotes(ctx, "Go")
	if err != nil {
		t.Fatalf("SearchNotes() failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("search 'Go' should return 2 results, got %d", len(results))
	}

	results, err = fs.SearchNotes(ctx, "Python")
	if err != nil {
		t.Fatalf("SearchNotes() failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("search 'Python' should return 1 result, got %d", len(results))
	}
}

func TestDeleteNote(t *testing.T) {
	fs := newTestFS(t)
	ctx := context.Background()

	note := storage.NewNote("Note to delete", "Temporary content")
	if err := fs.SaveNote(ctx, note); err != nil {
		t.Fatalf("SaveNote() failed: %v", err)
	}

	if _, err := os.Stat(note.FilePath); os.IsNotExist(err) {
		t.Fatal("file was not created")
	}

	if err := fs.DeleteNote(ctx, note.ID); err != nil {
		t.Fatalf("DeleteNote() failed: %v", err)
	}

	if _, err := os.Stat(note.FilePath); !os.IsNotExist(err) {
		t.Fatal("file was not deleted")
	}

	err := fs.DeleteNote(ctx, "non-existent-note")
	if err == nil {
	 t.Error("DeleteNote() should return an error for a non-existent note")
	}
}

func TestParseNoteLegacyFormat(t *testing.T) {
	fs := newTestFS(t)
	ctx := context.Background()

	legacyPath := fs.NotesDir() + string(os.PathSeparator) + "legacy-id.md"
	legacyContent := "# Legacy Note\n\nOld content without front matter"
	if err := os.WriteFile(legacyPath, []byte(legacyContent), 0644); err != nil {
		t.Fatalf("failed to write legacy note: %v", err)
	}

	note, err := fs.GetNote(ctx, "legacy-id")
	if err != nil {
		t.Fatalf("GetNote() failed: %v", err)
	}
	if note.Title != "Legacy Note" {
		t.Errorf("expected title Legacy Note, got %q", note.Title)
	}
	if note.Content != "Old content without front matter" {
		t.Errorf("unexpected content: %q", note.Content)
	}
}
