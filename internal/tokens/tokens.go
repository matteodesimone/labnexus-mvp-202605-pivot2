// Package tokens stima il numero di token in input usando char_count/4 (FR-6).
// Warning a 70% del context_window, errore esplicito a 100%.
package tokens

import "errors"

const (
	WarnRatio  = 0.70
	BlockRatio = 1.00
)

// Estimate ritorna il numero stimato di token per textLen caratteri (textLen/4).
// È un'approssimazione sufficiente per il warning a 70%; un tokenizer vero è fuori scope Sprint 1.
func Estimate(textLen int) int {
	if textLen < 0 {
		return 0
	}
	return textLen / 4
}

// CheckResult è l'esito del check rispetto al context window del profilo.
type CheckResult struct {
	Tokens  int
	Ratio   float64
	Warning bool // tokens ≥ 70% del context
	Block   bool // tokens ≥ 100% del context
}

// Check verifica i token stimati contro contextWindow.
func Check(textLen, contextWindow int) (*CheckResult, error) {
	if contextWindow <= 0 {
		return nil, errors.New("tokens: context_window deve essere > 0")
	}
	t := Estimate(textLen)
	ratio := float64(t) / float64(contextWindow)
	return &CheckResult{
		Tokens:  t,
		Ratio:   ratio,
		Warning: ratio >= WarnRatio,
		Block:   ratio >= BlockRatio,
	}, nil
}
