package provider

import (
	"encoding/json"
	"strings"
	"testing"
)

// Regressione: max_tokens e temperature DEVONO finire nel JSON inviato a
// eurouter (prima erano ignorati → il modello usava i propri default, causa
// probabile delle risposte vuote su modelli thinking).
func TestBuildEurouterRequest_SendsParams(t *testing.T) {
	req := buildEurouterRequest("sys", "usr", Options{Modello: "kimi-k2.6", MaxTokens: 8192, Temperature: 0.9})
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, want := range []string{`"max_tokens":8192`, `"temperature":0.9`, `"model":"kimi-k2.6"`} {
		if !strings.Contains(s, want) {
			t.Errorf("la richiesta deve contenere %s\n--- json ---\n%s", want, s)
		}
	}
}

// max_tokens=0 → omesso (lascia il default del modello), così l'utente può
// disattivare l'invio se un modello strict lo rifiuta.
func TestBuildEurouterRequest_MaxTokensZeroOmitted(t *testing.T) {
	b, _ := json.Marshal(buildEurouterRequest("s", "u", Options{Modello: "m", MaxTokens: 0, Temperature: 0.9}))
	if strings.Contains(string(b), "max_tokens") {
		t.Errorf("max_tokens=0 deve essere omesso, got %s", b)
	}
}
