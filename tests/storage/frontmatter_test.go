package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/N95Ryan/leaf/internal/storage"
)

func TestParseFrontMatterRoundTrip(t *testing.T) {
	fs := newTestFS(t)

	note := storage.NewNote("Metadata Note", "Body")
	created := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	note.CreatedAt = created

	if err := fs.SaveNote(context.Background(), note); err != nil {
		t.Fatalf("SaveNote() failed: %v", err)
	}

	loaded, err := fs.GetNote(context.Background(), note.ID)
	if err != nil {
		t.Fatalf("GetNote() failed: %v", err)
	}

	if !loaded.CreatedAt.Equal(created) {
		t.Fatalf("expected CreatedAt %v, got %v", created, loaded.CreatedAt)
	}
	if loaded.UpdatedAt.IsZero() {
		t.Fatal("expected UpdatedAt to be set")
	}
}
