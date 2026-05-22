package runner

import "testing"

// TestStripFenceWrapping_DetectsAndStripsYamlWrap: caso visto in produzione
// 2026-05-22T11:17 — il modello Qwen3.5-122B ha racchiuso l'intero output
// (frontmatter + markdown body) dentro un singolo ```yaml ... ``` fence.
// Risultato: il PDF mostrava code block monospace invece di markdown reso.
// Il sanitizer deve rilevare il wrap e restituire il contenuto raw.
func TestStripFenceWrapping_DetectsAndStripsYamlWrap(t *testing.T) {
	wrapped := "```yaml\n---\nprofilo: revisione\n---\n# Titolo\n\nCorpo.\n```"
	want := "---\nprofilo: revisione\n---\n# Titolo\n\nCorpo."
	got := stripFenceWrapping(wrapped)
	if got != want {
		t.Errorf("stripFenceWrapping:\n got: %q\nwant: %q", got, want)
	}
}

// TestStripFenceWrapping_DetectsAndStripsMarkdownWrap: variante con ```markdown.
func TestStripFenceWrapping_DetectsAndStripsMarkdownWrap(t *testing.T) {
	wrapped := "```markdown\n# Titolo\n\nCorpo.\n```"
	want := "# Titolo\n\nCorpo."
	got := stripFenceWrapping(wrapped)
	if got != want {
		t.Errorf("stripFenceWrapping markdown:\n got: %q\nwant: %q", got, want)
	}
}

// TestStripFenceWrapping_DetectsAndStripsPlainFence: ``` senza linguaggio.
func TestStripFenceWrapping_DetectsAndStripsPlainFence(t *testing.T) {
	wrapped := "```\n# Titolo\n```"
	want := "# Titolo"
	got := stripFenceWrapping(wrapped)
	if got != want {
		t.Errorf("stripFenceWrapping plain:\n got: %q\nwant: %q", got, want)
	}
}

// TestStripFenceWrapping_NoOpWhenNoFence: body senza fence wrapping resta invariato.
func TestStripFenceWrapping_NoOpWhenNoFence(t *testing.T) {
	body := "# Titolo\n\nCorpo del documento.\n"
	got := stripFenceWrapping(body)
	if got != body {
		t.Errorf("body senza fence deve restare invariato, got: %q", got)
	}
}

// TestStripFenceWrapping_NoOpWhenFenceIsInternalCode: se il fence è all'interno
// (es. esempio di codice nel mezzo del documento), NON strippare. La prima riga
// deve essere ``` per attivare lo strip.
func TestStripFenceWrapping_NoOpWhenFenceIsInternalCode(t *testing.T) {
	body := "# Titolo\n\nEsempio:\n\n```go\nfmt.Println()\n```\n\nFine."
	got := stripFenceWrapping(body)
	if got != body {
		t.Errorf("fence interno NON deve attivare strip, got: %q", got)
	}
}

// TestStripFenceWrapping_NoOpWhenOnlyOpeningFence: edge case difensivo —
// se manca il fence di chiusura, non strippare (body malformato, lascia decidere).
func TestStripFenceWrapping_NoOpWhenOnlyOpeningFence(t *testing.T) {
	body := "```yaml\n---\nprofilo: x\n---\n# Titolo\n"
	got := stripFenceWrapping(body)
	if got != body {
		t.Errorf("missing closing fence: lascia invariato, got: %q", got)
	}
}

// TestStripFenceWrapping_HandlesLeadingTrailingWhitespace: il modello a volte
// aggiunge righe vuote prima/dopo il fence. Strippiamo whitespace pre-detection.
func TestStripFenceWrapping_HandlesLeadingTrailingWhitespace(t *testing.T) {
	wrapped := "\n\n```yaml\n# Titolo\n```\n\n"
	want := "# Titolo"
	got := stripFenceWrapping(wrapped)
	if got != want {
		t.Errorf("whitespace stripping:\n got: %q\nwant: %q", got, want)
	}
}
