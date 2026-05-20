// Package prompt compone il prompt LLM (system + user) per una capability LabNexus (FR-5).
package prompt

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// Composed è il prompt risultante.
type Composed struct {
	System string
	User   string
}

// Compose costruisce system e user message secondo FR-5:
//
//	system = concat(kbFileTexts, "\n\n---\n\n")
//	user   = triggerPrompt + "\n\n## File di input\n\n" + concat(inputFiles, "\n\n--- FILE: <name> ---\n\n<content>")
//
// inputFiles è un mapping nome→testo estratto. L'ordine è deterministico (nomi ordinati alfabeticamente).
func Compose(triggerPrompt string, kbFileTexts []string, inputFiles map[string]string) (*Composed, error) {
	if strings.TrimSpace(triggerPrompt) == "" {
		return nil, errors.New("prompt: trigger_prompt è vuoto")
	}

	system := strings.Join(kbFileTexts, "\n\n---\n\n")

	var b strings.Builder
	b.WriteString(triggerPrompt)
	b.WriteString("\n\n## File di input\n\n")

	if len(inputFiles) > 0 {
		names := make([]string, 0, len(inputFiles))
		for n := range inputFiles {
			names = append(names, n)
		}
		sort.Strings(names)

		blocks := make([]string, 0, len(names))
		for _, name := range names {
			blocks = append(blocks, fmt.Sprintf("--- FILE: %s ---\n\n%s", sanitizeFilename(name), inputFiles[name]))
		}
		b.WriteString(strings.Join(blocks, "\n\n"))
	}

	return &Composed{System: system, User: b.String()}, nil
}

// reChatRoleMarker matcha marker di turno conversazionale tipici dei chat
// template (case-insensitive). Bug #004 fix: questi marker in un filename
// possono "iniettare" un cambio di turno nel prompt e bypassare le
// istruzioni di sistema.
var reChatRoleMarker = regexp.MustCompile(`(?i)\b(system|user|assistant|developer):`)

const maxFilenameLen = 200

// sanitizeFilename normalizza un filename prima di interpolarlo nel marker
// "--- FILE: <name> ---" del prompt. Rimuove caratteri di controllo
// (newline, null byte, ecc.), neutralizza marker di chat role, trunca a
// maxFilenameLen. Bug #004 (.pipeline/bugs/prompt-injection-via-filename.md).
//
// Filename leciti italiani (accenti, spazi, punti, trattini, underscore)
// passano invariati.
func sanitizeFilename(name string) string {
	// 1. Rimuove tutti i caratteri non printable / di controllo (newline,
	//    null, tab, ecc.).
	var b strings.Builder
	b.Grow(len(name))
	for _, r := range name {
		if unicode.IsControl(r) {
			b.WriteRune('_')
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	// 2. Neutralizza marker di chat role ("System:", "User:", "Assistant:",
	//    "Developer:") sostituendoli con la versione senza i due punti
	//    (System_, User_, ...). Case-insensitive.
	out = reChatRoleMarker.ReplaceAllStringFunc(out, func(m string) string {
		// Mantiene il prefisso originale ma sostituisce ':' con '_'.
		return m[:len(m)-1] + "_"
	})
	// 3. Trunca se necessario (preserva l'estensione se possibile).
	if len(out) > maxFilenameLen {
		if dot := strings.LastIndexByte(out, '.'); dot > 0 && len(out)-dot < 20 {
			ext := out[dot:]
			out = out[:maxFilenameLen-len(ext)] + ext
		} else {
			out = out[:maxFilenameLen]
		}
	}
	return out
}
