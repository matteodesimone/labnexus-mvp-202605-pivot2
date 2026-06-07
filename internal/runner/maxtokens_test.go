package runner

import (
	"errors"
	"strings"
	"testing"
)

// conservativePromptTokens converte char/4 nella densità reale di kimi (~3,2
// char/tok): deve superare char/4 e avvicinarsi al reale (A: 733445 char =
// 229169 token; char/4 = 183361 → atteso ~229201).
func TestConservativePromptTokens(t *testing.T) {
	if got := conservativePromptTokens(183361); got <= 183361 {
		t.Errorf("la stima reale deve superare char/4: got %d", got)
	}
	if got := conservativePromptTokens(183361); got < 229169 {
		t.Errorf("la stima reale (%d) deve almeno raggiungere il conteggio reale 229169", got)
	}
}

// max_tokens=0 (auto) → budget = context − stima_reale(char/3,2) − buffer, così la
// prima chiamata entra al primo colpo. MAI omesso (senza, default provider ~1024).
func TestResolveOutputBudget_AutoComputesFromContext(t *testing.T) {
	got, err := resolveOutputBudget(0, 70000, 256000)
	if err != nil {
		t.Fatalf("non doveva bloccare: %v", err)
	}
	want := 256000 - conservativePromptTokens(70000) - outputBuffer
	if got != want {
		t.Errorf("auto = %d, want %d (context − stima_reale − buffer)", got, want)
	}
	// prompt più grande → budget più piccolo
	big, _ := resolveOutputBudget(0, 120000, 256000)
	if big >= got {
		t.Errorf("un prompt più grande deve lasciare meno output: big=%d, piccolo=%d", big, got)
	}
}

// Un valore esplicito > 0 nel config viene rispettato come tetto.
func TestResolveOutputBudget_HonorsExplicitCap(t *testing.T) {
	got, err := resolveOutputBudget(20000, 70000, 256000)
	if err != nil {
		t.Fatalf("non doveva bloccare: %v", err)
	}
	if got != 20000 {
		t.Errorf("budget = %d, want 20000 (tetto esplicito rispettato)", got)
	}
}

// Auto con prompt grande (sotto la guardia ottimistica char/4): la stima reale
// porta il budget sotto il minimo → floor al minimo; sarà il provider a decidere.
func TestResolveOutputBudget_AutoFloorsToMinimum(t *testing.T) {
	got, err := resolveOutputBudget(0, 200000, 256000)
	if err != nil {
		t.Fatalf("non doveva bloccare (sotto la guardia ottimistica): %v", err)
	}
	if got != minOutputBudget {
		t.Errorf("budget = %d, want floor %d", got, minOutputBudget)
	}
}

// Guardia LENIENTE: blocca solo se nemmeno la stima ottimistica (char/4) lascia
// spazio al minimo (input palesemente enorme).
func TestResolveOutputBudget_BlocksOnlyWhenOptimisticTooBig(t *testing.T) {
	if _, err := resolveOutputBudget(0, 250000, 256000); err == nil {
		t.Error("atteso blocco quando anche la stima ottimistica non lascia spazio")
	} else if !strings.Contains(err.Error(), "input troppo grande") {
		t.Errorf("messaggio poco chiaro: %v", err)
	}
}

// L'errore di context overflow va riconosciuto in TUTTI E 3 i formati: vLLM "at
// least N", vLLM esatto "N tokens from the input messages", gateway 400.
func TestParseContextOverflow(t *testing.T) {
	// Formato vLLM stima ("at least N" / value=N): lower bound.
	msgA := `eurouter: errore dal provider (evento SSE su HTTP 200): {"message":"400 This model's maximum context length is 262144 tokens. However, you requested 65539 output tokens and your prompt contains at least 196606 input tokens, for a total of at least 262145 tokens. (parameter=input_tokens, value=196606)","type":"stream_error","code":"BAD_REQUEST"}`
	if n, ok := parseContextOverflow(errors.New(msgA), 65539); !ok || n != 196606 {
		t.Errorf("vLLM stima = (%d,%v), want (196606,true)", n, ok)
	}
	// Formato vLLM esatto ("N tokens from the input messages"): reale, preferito.
	msgB := `eurouter: errore dal provider (evento SSE su HTTP 200): {"message":"400 Requested token count exceeds the model's maximum context length of 262144 tokens. You requested a total of 274919 tokens: 229169 tokens from the input messages and 45750 tokens for the completion. Please reduce...","type":"stream_error","code":"BAD_REQUEST"}`
	if n, ok := parseContextOverflow(errors.New(msgB), 45750); !ok || n != 229169 {
		t.Errorf("vLLM esatto = (%d,%v), want (229169,true)", n, ok)
	}
	// Formato GATEWAY 400: "Estimated total tokens (T) exceeds model context window".
	// prompt = T − max_tokens inviato = 259238 − 78139 = 181099.
	msgG := `eurouter: HTTP 400 POST https://api.eurouter.ai/v1/chat/completions (modello="kimi-k2.6"): {"error":{"code":400,"message":"Estimated total tokens (259238) exceeds model context window (256000). Reduce message length or max_tokens.","type":"invalid_request_error","param":null},"requestId":"x"}`
	if n, ok := parseContextOverflow(errors.New(msgG), 78139); !ok || n != 181099 {
		t.Errorf("gateway = (%d,%v), want (181099,true)", n, ok)
	}
	// errore non-overflow → non riconosciuto
	if _, ok := parseContextOverflow(errors.New("eurouter: HTTP 500 boom"), 0); ok {
		t.Error("un errore non-overflow non deve essere riconosciuto come tale")
	}
	if _, ok := parseContextOverflow(nil, 0); ok {
		t.Error("nil non deve essere overflow")
	}
}

// Il budget di calibrazione usa il conteggio del provider + buffer minimo (output
// massimo) e blocca se non resta spazio per l'output minimo.
func TestRetryOutputBudget(t *testing.T) {
	got, fits := retryOutputBudget(229169, 256000)
	if !fits {
		t.Fatal("doveva esserci spazio per l'output")
	}
	want := 256000 - 229169 - outputBuffer
	if got != want {
		t.Errorf("retryOutputBudget = %d, want %d (context − reale − buffer)", got, want)
	}
	// input troppo grande → niente spazio per il minimo → blocco
	if _, fits := retryOutputBudget(245000, 256000); fits {
		t.Error("input troppo grande: non deve esserci spazio per l'output minimo")
	}
}
