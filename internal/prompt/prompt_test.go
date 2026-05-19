package prompt_test

import (
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/prompt"
)

func TestCompose_SystemJoinsKbFilesWithSeparator(t *testing.T) {
	c, err := prompt.Compose("trigger", []string{"alpha", "beta"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(c.System, "alpha") || !strings.Contains(c.System, "beta") {
		t.Fatalf("system message should contain both kb texts, got: %q", c.System)
	}
	if !strings.Contains(c.System, "\n\n---\n\n") {
		t.Fatalf("system message should use \\n\\n---\\n\\n separator (FR-5), got: %q", c.System)
	}
}

func TestCompose_UserStartsWithTrigger(t *testing.T) {
	c, err := prompt.Compose("XYZ trigger", []string{"kb"}, map[string]string{"a.md": "alpha"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(c.User, "XYZ trigger") {
		t.Fatalf("user message must start with trigger_prompt, got: %q", c.User)
	}
}

func TestCompose_UserContainsFileBlocks(t *testing.T) {
	files := map[string]string{
		"doc1.md":   "alpha-content",
		"doc2.docx": "beta-content",
	}
	c, err := prompt.Compose("t", nil, files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(c.User, "## File di input") {
		t.Fatalf("user must contain '## File di input' header, got: %q", c.User)
	}
	for name := range files {
		marker := "--- FILE: " + name + " ---"
		if !strings.Contains(c.User, marker) {
			t.Fatalf("user must contain marker %q, got: %q", marker, c.User)
		}
	}
}
