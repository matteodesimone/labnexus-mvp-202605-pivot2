package runner

import "testing"

func TestEffectiveMaxTokens(t *testing.T) {
	cases := []struct {
		name                              string
		configMax, prompt, context, want int
	}{
		{"auto: usa tutto il disponibile", 0, 70000, 262000, 262000 - 70000 - 13100},
		{"config sotto il disponibile: rispettato", 65536, 70000, 262000, 65536},
		{"config troppo alto per il context: ristretto", 131072, 200000, 262000, 262000 - 200000 - 13100},
		{"context sconosciuto: usa config", 8192, 70000, 0, 8192},
		{"input enorme: floor 1024", 0, 261000, 262000, 1024},
	}
	for _, c := range cases {
		if got := effectiveMaxTokens(c.configMax, c.prompt, c.context); got != c.want {
			t.Errorf("%s: effectiveMaxTokens(%d,%d,%d)=%d, want %d", c.name, c.configMax, c.prompt, c.context, got, c.want)
		}
	}
}
