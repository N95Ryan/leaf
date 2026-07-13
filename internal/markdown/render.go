package markdown

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

var parser = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
)

// Render converts markdown content to styled terminal output.
func Render(content string, width int) (string, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return "", nil
	}

	source := []byte(content)
	doc := parser.Parser().Parse(text.NewReader(source))

	var b strings.Builder
	if err := renderNode(&b, doc, source, width); err != nil {
		return "", err
	}

	return strings.TrimRight(b.String(), "\n"), nil
}

func renderNode(b *strings.Builder, n ast.Node, source []byte, width int) error {
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		switch node := child.(type) {
		case *ast.Heading:
			level := node.Level
			textValue := gatherText(node, source)
			prefix := strings.Repeat("#", level) + " "
			b.WriteString(bold(prefix + textValue))
			b.WriteString("\n\n")
		case *ast.Paragraph:
			textValue := gatherText(node, source)
			if textValue != "" {
				b.WriteString(wrapText(textValue, width))
				b.WriteString("\n\n")
			}
		case *ast.TextBlock:
			textValue := gatherText(node, source)
			if textValue != "" {
				b.WriteString(wrapText(textValue, width))
				b.WriteString("\n\n")
			}
		case *ast.ListItem:
			textValue := gatherText(node, source)
			b.WriteString("- ")
			b.WriteString(wrapText(textValue, width))
			b.WriteString("\n")
		case *ast.FencedCodeBlock:
			lang := string(node.Language(source))
			code := string(node.Text(source))
			highlighted, err := highlightCode(code, lang)
			if err != nil {
				b.WriteString(code)
			} else {
				b.WriteString(highlighted)
			}
			b.WriteString("\n\n")
		case *ast.CodeBlock:
			code := string(node.Text(source))
			b.WriteString(code)
			b.WriteString("\n\n")
		case *ast.Blockquote:
			textValue := gatherText(node, source)
			for _, line := range strings.Split(textValue, "\n") {
				if strings.TrimSpace(line) == "" {
					continue
				}
				b.WriteString("> ")
				b.WriteString(line)
				b.WriteString("\n")
			}
			b.WriteString("\n")
		default:
			if err := renderNode(b, child, source, width); err != nil {
				return err
			}
		}
	}
	return nil
}

func gatherText(n ast.Node, source []byte) string {
	var b strings.Builder
	_ = ast.Walk(n, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if t, ok := node.(*ast.Text); ok {
			b.Write(t.Segment.Value(source))
		}
		return ast.WalkContinue, nil
	})
	return strings.TrimSpace(b.String())
}

func highlightCode(code, lang string) (string, error) {
	if lang == "" {
		lang = "text"
	}
	var buf bytes.Buffer
	if err := quick.Highlight(&buf, code, lang, "terminal256", "monokai"); err != nil {
		return "", fmt.Errorf("highlight code block: %w", err)
	}
	return buf.String(), nil
}

func bold(s string) string {
	return "\x1b[1m" + s + "\x1b[0m"
}

func wrapText(s string, width int) string {
	if width <= 0 {
		return s
	}
	words := strings.Fields(s)
	if len(words) == 0 {
		return s
	}

	var lines []string
	line := words[0]
	for _, word := range words[1:] {
		if len(line)+1+len(word) > width {
			lines = append(lines, line)
			line = word
			continue
		}
		line += " " + word
	}
	lines = append(lines, line)
	return strings.Join(lines, "\n")
}
