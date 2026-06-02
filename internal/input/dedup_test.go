package input_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/labnexus/labnexus/internal/input"
)

// De-dup: un file designato come trigger (trigger_prompt_file) ed esistente
// anche nella cartella di input NON deve essere incluso tra i file di input,
// altrimenti il prompt verrebbe inviato due volte (una come trigger, una come
// dato). ParseDir accetta path relativi POSIX da escludere.
func TestParseDir_ExcludesTriggerFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "Prompt_INPUT.txt"), []byte("ISTRUZIONI DEL TASK"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "dato.txt"), []byte("DATI REALI"), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := input.ParseDir(dir, "Prompt_INPUT.txt")
	if err != nil {
		t.Fatalf("ParseDir: %v", err)
	}
	for _, f := range res.Files {
		if f.Name == "Prompt_INPUT.txt" {
			t.Errorf("il file-trigger %q non deve comparire tra i file di input (de-dup)", f.Name)
		}
	}
	if len(res.Files) != 1 || res.Files[0].Name != "dato.txt" {
		t.Errorf("atteso solo dato.txt tra gli input, got %d file: %+v", len(res.Files), res.Files)
	}
}

// I file infrastrutturali (Esegui.command, _labnexus.toml, .DS_Store) NON
// devono finire tra i dati inviati al modello. Nota: LEGGIMI.txt NON è
// infrastrutturale (può contenere dati veri) — si esclude via campo `exclude`
// del _labnexus.toml, non hard-coded.
func TestParseDir_SkipsInfraFiles(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"Esegui.command", "dato.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x contenuto"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	res, err := input.ParseDir(dir)
	if err != nil {
		t.Fatalf("ParseDir: %v", err)
	}
	for _, f := range res.Files {
		if f.Name == "Esegui.command" {
			t.Errorf("file infrastrutturale %q non deve essere mandato al modello", f.Name)
		}
	}
	if len(res.Files) != 1 || res.Files[0].Name != "dato.txt" {
		t.Errorf("atteso solo dato.txt, got %+v", res.Files)
	}
}

// L'esclusione deve normalizzare il path (filepath.Clean): trigger_prompt_file
// può arrivare in forma non canonica (es. "./Prompt.txt") ma il walk produce
// rel-path puliti. Senza Clean la de-dup salterebbe silenziosamente.
func TestParseDir_ExcludeNormalizesPath(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "Prompt_INPUT.txt"), []byte("TRIGGER"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "dato.txt"), []byte("DATI"), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := input.ParseDir(dir, "./Prompt_INPUT.txt")
	if err != nil {
		t.Fatalf("ParseDir: %v", err)
	}
	for _, f := range res.Files {
		if f.Name == "Prompt_INPUT.txt" {
			t.Errorf("path non canonico './Prompt_INPUT.txt' deve comunque escludere il file (Clean)")
		}
	}
}

// Esclusione per path relativo (non solo basename): un file omonimo in una
// sottocartella NON deve essere escluso per sbaglio.
func TestParseDir_ExcludeMatchesRelPathNotBasename(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "Prompt_INPUT.txt"), []byte("TRIGGER"), 0o644); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "Prompt_INPUT.txt"), []byte("OMONIMO MA DATO"), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := input.ParseDir(dir, "Prompt_INPUT.txt")
	if err != nil {
		t.Fatalf("ParseDir: %v", err)
	}
	var foundSub bool
	for _, f := range res.Files {
		if f.Name == "Prompt_INPUT.txt" {
			t.Errorf("il trigger root %q non deve comparire", f.Name)
		}
		if f.Name == "sub/Prompt_INPUT.txt" {
			foundSub = true
		}
	}
	if !foundSub {
		t.Error("l'omonimo in sub/ deve restare incluso (esclusione per rel-path, non basename)")
	}
}
