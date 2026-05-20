package features

import (
	"strings"
	"testing"
)

// Test unitari delle assertion strutturali di Fetta 2 — H4 fix dalla review
// loop 1. Coprono i casi negativi (body malformato) che gli scenari BDD
// "happy path" non possono esercitare perché il canned content è controllato
// da noi.

func TestCountNumberedSections(t *testing.T) {
	cases := []struct {
		name string
		body string
		want int
	}{
		{"vuoto", "", 0},
		{"nessuna numerazione", "## Sintesi\n## Conclusione", 0},
		{"3 sezioni", "## 1. A\nblabla\n## 2. B\nblabla\n## 3. C\n", 3},
		{"13 sezioni", strings.Repeat("## %d. X\n", 1) + buildNumberedSections(13), 13},
		{"5 sezioni equipment-alert canned", buildNumberedSections(5), 5},
		{"non match se manca il punto", "## 1 senza punto\n## 2 senza punto", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := countNumberedSections(tc.body); got != tc.want {
				t.Fatalf("countNumberedSections: got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestHasExactNumberedSectionSequence(t *testing.T) {
	// Happy path 1..5.
	ok5 := buildNumberedSections(5)
	if err := hasExactNumberedSectionSequence(ok5, 5); err != nil {
		t.Fatalf("happy 5: %v", err)
	}
	// Conteggio sbagliato.
	if err := hasExactNumberedSectionSequence(ok5, 13); err == nil {
		t.Fatalf("5 sezioni vs atteso 13: atteso errore")
	}
	// Sequenza con gap (manca 3).
	withGap := "## 1. A\n## 2. B\n## 4. D\n## 5. E"
	if err := hasExactNumberedSectionSequence(withGap, 5); err == nil {
		t.Fatalf("sequenza con gap (1,2,4,5): atteso errore")
	}
	// Duplicato (## 2 due volte).
	withDup := "## 1. A\n## 2. B\n## 2. B bis\n## 3. C\n## 4. D\n## 5. E"
	if err := hasExactNumberedSectionSequence(withDup, 5); err == nil {
		t.Fatalf("sequenza con duplicato (## 2 due volte): atteso errore")
	}
	// Out of order (## 3 prima di ## 2).
	outOfOrder := "## 1. A\n## 3. C\n## 2. B\n## 4. D\n## 5. E"
	if err := hasExactNumberedSectionSequence(outOfOrder, 5); err == nil {
		t.Fatalf("sequenza out-of-order: atteso errore")
	}
	// Sequence 1..13.
	ok13 := buildNumberedSections(13)
	if err := hasExactNumberedSectionSequence(ok13, 13); err != nil {
		t.Fatalf("happy 13: %v", err)
	}
}

func TestHasSintesiInizialeECoerenza(t *testing.T) {
	good := "## Sintesi iniziale\nIntro coerente con la 13.\n## 13. Decisioni\nDeliberato X."
	missingSintesi := "## Premessa\nNo intro.\n## 13. Decisioni\nDeliberato X."
	missingDecisioni := "## Sintesi iniziale\nIntro.\n## 13. Conclusioni\nFine."
	missingSezione13 := "## Sintesi iniziale\nIntro.\n## Sezione X\nDecisioni varie."

	if err := hasSintesiInizialeECoerenza(good); err != nil {
		t.Fatalf("good caso fallisce: %v", err)
	}
	if err := hasSintesiInizialeECoerenza(missingSintesi); err == nil {
		t.Fatalf("missingSintesi: atteso errore, ottenuto nil")
	}
	if err := hasSintesiInizialeECoerenza(missingDecisioni); err == nil {
		t.Fatalf("missingDecisioni: atteso errore, ottenuto nil")
	}
	if err := hasSintesiInizialeECoerenza(missingSezione13); err == nil {
		t.Fatalf("missingSezione13: atteso errore, ottenuto nil")
	}
}

func TestCountNumberedQuestions(t *testing.T) {
	body := "1. Prima domanda?\n2. Seconda?\nnon numerato\n3. Terza?\n10. Decima?"
	if got := countNumberedQuestions(body); got != 4 {
		t.Fatalf("got %d, want 4", got)
	}
	if got := countNumberedQuestions(""); got != 0 {
		t.Fatalf("vuoto: got %d, want 0", got)
	}
}

func TestCountAreaHeadings(t *testing.T) {
	body := "### Area 1\ntext\n### Area 2: subtitle\ntext\n### Not Area\n### Area 3\n"
	if got := countAreaHeadings(body); got != 3 {
		t.Fatalf("got %d, want 3", got)
	}
}

func TestSplitAreaBlocksAndPerAreaChecks(t *testing.T) {
	body := `## Checklist

### Area 1
Campioni documentali: certificati. Livello di rischio: Alto.
1. Domanda uno?

### Area 2
Campioni documentali: piani. Livello di rischio: Medio.
2. Domanda due?

### Area 3 — incompleta
1. Domanda tre senza campioni né rischio?
`
	blocks := splitAreaBlocks(body)
	if len(blocks) != 3 {
		t.Fatalf("split: got %d blocchi, want 3", len(blocks))
	}
	// Area 1 e 2 hanno campioni+rischio; Area 3 no.
	if err := areaBlockHasCampioniERischio(blocks[0]); err != nil {
		t.Fatalf("area1 doveva passare: %v", err)
	}
	if err := areaBlockHasCampioniERischio(blocks[1]); err != nil {
		t.Fatalf("area2 doveva passare: %v", err)
	}
	if err := areaBlockHasCampioniERischio(blocks[2]); err == nil {
		t.Fatalf("area3 incompleta doveva fallire ma è passata")
	}
}

func TestHasFiveEquipmentAlertSections(t *testing.T) {
	good := `## 1. Stato apparecchiatura
text
## 2. Rischio tecnico
text
## 3. Azioni proposte
text
## 4. Bozza email al fornitore
text
## 5. Checklist al rientro
text`
	if err := hasFiveEquipmentAlertSections(good); err != nil {
		t.Fatalf("happy path fallisce: %v", err)
	}

	tooFew := `## 1. Stato
## 2. Rischio
## 3. Azioni
## 4. Email`
	if err := hasFiveEquipmentAlertSections(tooFew); err == nil {
		t.Fatalf("4 sezioni: atteso errore")
	}

	tooMany := `## 1. A
## 2. B
## 3. C
## 4. D email
## 5. E checklist
## 6. F extra`
	if err := hasFiveEquipmentAlertSections(tooMany); err == nil {
		t.Fatalf("6 sezioni: atteso errore")
	}

	// 5 sezioni ma sezione 4 senza parola chiave email/messaggio/...
	missingSec4 := `## 1. Stato
## 2. Rischio
## 3. Azioni
## 4. Riepilogo
## 5. Checklist`
	if err := hasFiveEquipmentAlertSections(missingSec4); err == nil {
		t.Fatalf("sez 4 senza email/sinonimi: atteso errore")
	}

	// L2 fix: sec4 con sinonimi accettati.
	syn1 := strings.Replace(good, "Bozza email al fornitore", "Bozza messaggio al fornitore", 1)
	if err := hasFiveEquipmentAlertSections(syn1); err != nil {
		t.Fatalf("sinonimo 'messaggio' atteso passare: %v", err)
	}
	syn2 := strings.Replace(good, "Bozza email al fornitore", "Comunicazione al fornitore", 1)
	if err := hasFiveEquipmentAlertSections(syn2); err != nil {
		t.Fatalf("sinonimo 'comunicaz' atteso passare: %v", err)
	}
}

func TestExtractPCMCodes(t *testing.T) {
	text := "Usiamo PCM-01, poi PCM-04 e infine PCM-09. PCM-01 ripetuto. Non confondere con SCM-99."
	got := extractPCMCodes(text)
	want := map[string]bool{"PCM-01": true, "PCM-04": true, "PCM-09": true}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for k := range want {
		if !got[k] {
			t.Fatalf("mancante: %s", k)
		}
	}
}

func TestCheckMetodiOnlyFromSet(t *testing.T) {
	allowed := map[string]bool{"PCM-01": true, "PCM-04": true, "PCM-09": true}
	okBody := "Cita PCM-01 e PCM-04 — solo metodi del set."
	if err := checkMetodiOnlyFromSet(okBody, allowed); err != nil {
		t.Fatalf("body conforme fallisce: %v", err)
	}
	noPCM := "Nessun metodo nominato."
	if err := checkMetodiOnlyFromSet(noPCM, allowed); err == nil {
		t.Fatalf("body senza PCM: atteso errore")
	}
	hallucinated := "Cita PCM-01 e PCM-99 (non nel set)."
	if err := checkMetodiOnlyFromSet(hallucinated, allowed); err == nil {
		t.Fatalf("metodo allucinato PCM-99: atteso errore")
	}
}

func TestHasCausalConnective(t *testing.T) {
	cases := []struct {
		body string
		want bool
	}{
		{"poiché succede X, allora Y", true},
		{"X, quindi Y", true},
		{"perché vale Z", true},
		{"se A allora B", true},
		{"questo implica Z", true},
		{"frase neutra senza connettivi", false},
		{"qua nessun motivo è esposto", false},
	}
	for _, tc := range cases {
		if got := hasCausalConnective(tc.body); got != tc.want {
			t.Fatalf("body %q: got %v, want %v", tc.body, got, tc.want)
		}
	}
}

func TestExtractFrontmatter(t *testing.T) {
	with := "---\nkey: value\nother: x\n---\nbody"
	// extractFrontmatter ritorna il contenuto fino (non incluso) all'`\n---`
	// finale — convenzione semplice: niente newline finale incluso.
	if got := extractFrontmatter(with); got != "key: value\nother: x" {
		t.Fatalf("unexpected frontmatter: %q", got)
	}
	noFm := "no leading delim\nstuff"
	if got := extractFrontmatter(noFm); got != "" {
		t.Fatalf("non-frontmatter atteso vuoto, got %q", got)
	}
	missingClose := "---\nkey: value\nno closing delim"
	if got := extractFrontmatter(missingClose); got != "" {
		t.Fatalf("frontmatter aperto ma non chiuso: atteso vuoto, got %q", got)
	}
}

func TestFrontmatterContainsKV(t *testing.T) {
	fm := "tipo: equipment_alert\nstato: bozza_da_validare_qm\nprofilo_labnexus: equipment-alert\nlocale: true\nextra: cosa"
	expected := map[string]string{
		"tipo":             "equipment_alert",
		"profilo_labnexus": "equipment-alert",
		"locale":           "true",
	}
	if err := frontmatterContainsKV(fm, expected); err != nil {
		t.Fatalf("happy path fallisce: %v", err)
	}

	missing := map[string]string{"chiave_assente": "x"}
	if err := frontmatterContainsKV(fm, missing); err == nil {
		t.Fatalf("chiave assente: atteso errore")
	}

	wrong := map[string]string{"tipo": "altro_tipo"}
	if err := frontmatterContainsKV(fm, wrong); err == nil {
		t.Fatalf("valore sbagliato: atteso errore")
	}
}

// buildNumberedSections costruisce un body con N sezioni numerate ## 1. → ## N.
// per i test di counting.
func buildNumberedSections(n int) string {
	var b strings.Builder
	for i := 1; i <= n; i++ {
		b.WriteString("## ")
		b.WriteString(itoa(i))
		b.WriteString(". Section\nbody\n")
	}
	return b.String()
}

func itoa(i int) string {
	if i < 10 {
		return string(rune('0' + i))
	}
	// supporto 10–99 sufficiente per i test
	return string(rune('0'+i/10)) + string(rune('0'+i%10))
}
