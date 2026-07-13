package markdown_test

import (
	"strings"
	"testing"

	"github.com/N95Ryan/leaf/internal/markdown"
)

func TestRenderEmpty(t *testing.T) {
	out, err := markdown.Render("", 80)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}
	if out != "" {
		t.Fatalf("expected empty output, got %q", out)
	}
}

func TestRenderHeadingAndParagraph(t *testing.T) {
	content := "# Hello\n\nThis is a paragraph."
	out, err := markdown.Render(content, 80)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}
	if !strings.Contains(out, "Hello") {
		t.Fatalf("expected heading text in output, got %q", out)
	}
	if !strings.Contains(out, "This is a paragraph.") {
		t.Fatalf("expected paragraph in output, got %q", out)
	}
}

func TestRenderCodeBlock(t *testing.T) {
	content := "```go\nfunc main() {}\n```"
	out, err := markdown.Render(content, 80)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}
	if !strings.Contains(out, "func") || !strings.Contains(out, "main") {
		t.Fatalf("expected code block content, got %q", out)
	}
}

func TestRenderList(t *testing.T) {
	content := "- first\n- second"
	out, err := markdown.Render(content, 80)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}
	if !strings.Contains(out, "- first") || !strings.Contains(out, "- second") {
		t.Fatalf("expected list items, got %q", out)
	}
}

func TestRenderWrapsLongLines(t *testing.T) {
	content := "word " + strings.Repeat("long ", 30)
	out, err := markdown.Render(content, 20)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}
	if !strings.Contains(out, "\n") {
		t.Fatalf("expected wrapped output, got %q", out)
	}
}
