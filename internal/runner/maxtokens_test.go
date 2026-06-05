package runner

import (
	"strings"
	"testing"
)

func TestEffectiveMaxTokens(t *testing.T) {
	cases := []struct {
		name                              string
		configMax, prompt, context, want int
	}{
		{"auto: usa tutto il disponibile", 0, 70000, 262000, 262000 - 70000 - 13100},
		{"config sotto il disponibile: rispettato", 65536, 70000, 262000, 65536},
		{"config troppo alto per il context: ristretto", 131072, 200000, 262000, 262000 - 200000 - 13100},
		{"context sconosciuto: usa config", 8192, 70000, 0, 8192},
		{"input enorme: available negativo, NIENTE floor", 0, 261000, 262000, 262000 - 261000 - 13100},
	}
	for _, c := range cases {
		if got := effectiveMaxTokens(c.configMax, c.prompt, c.context); got != c.want {
			t.Errorf("%s: effectiveMaxTokens(%d,%d,%d)=%d, want %d", c.name, c.configMax, c.prompt, c.context, got, c.want)
		}
	}
}

// La guardia: se il prompt riempie il context (>~93%) il budget di output
// scende sotto il minimo utile → blocco con errore chiaro (non output troncato
// né sforamento del context).
func TestResolveOutputBudget_BlocksWhenTooLow(t *testing.T) {
	// prompt al ~95% → budget ~0 → blocco
	if _, err := resolveOutputBudget(0, 248900, 262000); err == nil {
		t.Error("atteso blocco quando il budget di output < minimo")
	} else if !strings.Contains(err.Error(), "input troppo grande") {
		t.Errorf("messaggio poco chiaro: %v", err)
	}
	// prompt al ~99,7% → available negativo → blocco (NON deve sforare il context)
	if _, err := resolveOutputBudget(0, 261216, 262000); err == nil {
		t.Error("atteso blocco con prompt ~99.7%")
	}
}

func TestResolveOutputBudget_OkWhenEnoughRoom(t *testing.T) {
	got, err := resolveOutputBudget(0, 70000, 262000)
	if err != nil {
		t.Fatalf("non doveva bloccare: %v", err)
	}
	if got != 262000-70000-13100 {
		t.Errorf("budget = %d, want %d", got, 262000-70000-13100)
	}
	// prompt all'88% → ~18900 budget → ok
	if _, err := resolveOutputBudget(0, 230000, 262000); err != nil {
		t.Errorf("prompt 88%% non deve bloccare: %v", err)
	}
}
