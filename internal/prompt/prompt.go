// Package prompt compone il prompt LLM (system + user) per una capability LabNexus (FR-5).
package prompt

import (
	"errors"
	"fmt"
	"sort"
	"strings"
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
			blocks = append(blocks, fmt.Sprintf("--- FILE: %s ---\n\n%s", name, inputFiles[name]))
		}
		b.WriteString(strings.Join(blocks, "\n\n"))
	}

	return &Composed{System: system, User: b.String()}, nil
}
