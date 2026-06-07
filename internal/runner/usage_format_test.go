package runner

import (
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/provider"
)

// Bug 2026-06-07 (CAPABILITY C): eurouter ha riportato reasoning=14335 >
// completion=6341. Non è un sottoinsieme → non va dichiarato "di cui" (sarebbe
// contraddittorio) e il totale deve includere il reasoning separato, non
// sottostimare (prima logava totale=152457 escludendo il reasoning).
func TestFormatUsage_ReasoningSeparateWhenExceedsCompletion(t *testing.T) {
	u := &provider.Usage{PromptTokens: 146116, CompletionTokens: 6341, ReasoningTokens: 14335, TotalTokens: 152457}
	s := formatUsage(u)
	if strings.Contains(s, "di cui") {
		t.Errorf("reasoning>completion: non deve usare 'di cui' (contraddittorio): %q", s)
	}
	if !strings.Contains(s, "166792") { // 146116 + 6341 + 14335
		t.Errorf("il totale deve includere il reasoning separato (166792): %q", s)
	}
}

// Caso normale (reasoning <= completion = sottoinsieme, come B/D/E/G): resta la
// dicitura "di cui" e si usa il totale riportato dal provider.
func TestFormatUsage_NormalSubset(t *testing.T) {
	u := &provider.Usage{PromptTokens: 55987, CompletionTokens: 26927, ReasoningTokens: 11894, TotalTokens: 82914}
	s := formatUsage(u)
	if !strings.Contains(s, "di cui reasoning=11894") {
		t.Errorf("caso normale deve restare 'di cui reasoning=11894': %q", s)
	}
	if !strings.Contains(s, "totale=82914") {
		t.Errorf("caso normale deve usare il totale del provider (82914): %q", s)
	}
}
