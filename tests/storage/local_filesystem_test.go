package storage_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/N95Ryan/leaf/internal/storage"
)

func TestNewLocalFileSystem(t *testing.T) {
	// Create an instance
	fs, err := storage.NewLocalFileSystem()
	if err != nil {
		t.Fatalf("NewLocalFileSystem() failed: %v", err)
	}

	if fs == nil {
		t.Fatal("NewLocalFileSystem() returned nil")
	}

	// Verify that the directory exists
	if _, err := os.Stat(fs.NotesDir()); os.IsNotExist(err) {
		t.Fatalf("notes directory was not created: %s", fs.NotesDir())
	}

	t.Logf("Notes directory created: %s", fs.NotesDir())
}

func TestSaveAndGetNote(t *testing.T) {
	fs, err := storage.NewLocalFileSystem()
	if err != nil {
		t.Fatalf("NewLocalFileSystem() failed: %v", err)
	}

	ctx := context.Background()

	// Create a note
	note := storage.NewNote("Test Note", "This is the test note content")
	originalID := note.ID

	// Save the note
	err = fs.SaveNote(ctx, note)
	if err != nil {
		t.Fatalf("SaveNote() failed: %v", err)
	}

	// Retrieve the note
	retrievedNote, err := fs.GetNote(ctx, originalID)
	if err != nil {
		t.Fatalf("GetNote() failed: %v", err)
	}

	// Verify the data
	if retrievedNote.Title != note.Title {
		t.Errorf("title mismatch: expected %q, got %q", note.Title, retrievedNote.Title)
	}

	if retrievedNote.Content != note.Content {
		t.Errorf("content mismatch: expected %q, got %q", note.Content, retrievedNote.Content)
	}

	if retrievedNote.ID != originalID {
		t.Errorf("ID mismatch: expected %q, got %q", originalID, retrievedNote.ID)
	}

	// Cleanup
	os.Remove(note.FilePath)
}

func TestListNotes(t *testing.T) {
	fs, err := storage.NewLocalFileSystem()
	if err != nil {
		t.Fatalf("NewLocalFileSystem() failed: %v", err)
	}

	ctx := context.Background()

	// Clean up the directory first
	os.RemoveAll(fs.NotesDir())
	os.MkdirAll(fs.NotesDir(), 0755)

	// Create multiple notes
	notes := []*storage.Note{
		storage.NewNote("Note 1", "Content 1"),
		storage.NewNote("Note 2", "Content 2"),
		storage.NewNote("Note 3", "Content 3"),
	}

	for _, note := range notes {
		if err := fs.SaveNote(ctx, note); err != nil {
			t.Fatalf("SaveNote() failed: %v", err)
		}
		time.Sleep(10 * time.Millisecond) // Small delay to differentiate UpdatedAt
	}

	// List the notes
	listedNotes, err := fs.ListNotes(ctx)
	if err != nil {
		t.Fatalf("ListNotes() failed: %v", err)
	}

	if len(listedNotes) != len(notes) {
		t.Errorf("note count mismatch: expected %d, got %d", len(notes), len(listedNotes))
	}

	// Verify order (descending by UpdatedAt)
	for i := 0; i < len(listedNotes)-1; i++ {
		if !listedNotes[i].UpdatedAt.After(listedNotes[i+1].UpdatedAt) {
			t.Error("notes are not sorted by UpdatedAt in descending order")
		}
	}

	// Cleanup
	for _, note := range notes {
		os.Remove(note.FilePath)
	}
}

func TestSearchNotes(t *testing.T) {
	fs, err := storage.NewLocalFileSystem()
	if err != nil {
		t.Fatalf("NewLocalFileSystem() failed: %v", err)
	}

	ctx := context.Background()

	// Cleanup
	os.RemoveAll(fs.NotesDir())
	os.MkdirAll(fs.NotesDir(), 0755)

	// Create notes with different titles and contents
	note1 := storage.NewNote("Go Tutorial", "Learn Go and concurrency")
	note2 := storage.NewNote("Python Tips", "Useful tips for Python")
	note3 := storage.NewNote("JavaScript Guide", "Complete Go guide for beginners")

	for _, note := range []*storage.Note{note1, note2, note3} {
		if err := fs.SaveNote(ctx, note); err != nil {
			t.Fatalf("SaveNote() failed: %v", err)
		}
	}

	// Search for "Go"
	results, err := fs.SearchNotes(ctx, "Go")
	if err != nil {
		t.Fatalf("SearchNotes() failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("search 'Go' should return 2 results, got %d", len(results))
	}

	// Search for "Python"
	results, err = fs.SearchNotes(ctx, "Python")
	if err != nil {
		t.Fatalf("SearchNotes() failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("search 'Python' should return 1 result, got %d", len(results))
	}

	// Cleanup
	for _, note := range []*storage.Note{note1, note2, note3} {
		os.Remove(note.FilePath)
	}
}

func TestDeleteNote(t *testing.T) {
	fs, err := storage.NewLocalFileSystem()
	if err != nil {
		t.Fatalf("NewLocalFileSystem() failed: %v", err)
	}

	ctx := context.Background()

	// Create and save a note
	note := storage.NewNote("Note to delete", "Temporary content")
	if err := fs.SaveNote(ctx, note); err != nil {
		t.Fatalf("SaveNote() failed: %v", err)
	}

	// Verify that the file exists
	if _, err := os.Stat(note.FilePath); os.IsNotExist(err) {
		t.Fatal("file was not created")
	}

	// Delete the note
	if err := fs.DeleteNote(ctx, note.ID); err != nil {
		t.Fatalf("DeleteNote() failed: %v", err)
	}

	// Verify that the file was deleted
	if _, err := os.Stat(note.FilePath); !os.IsNotExist(err) {
		t.Fatal("file was not deleted")
	}

	// Try to delete a non-existent note
	err = fs.DeleteNote(ctx, "non-existent-note")
	if err == nil {
		t.Error("DeleteNote() should return an error for a non-existent note")
	}
}
