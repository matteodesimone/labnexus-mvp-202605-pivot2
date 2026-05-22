package runner

import "strings"

// stripFenceWrapping rimuove un fence wrapping spurio attorno all'intero body
// generato dall'LLM. Caso visto in produzione 2026-05-22T11:17 con Qwen3.5-122B
// sul profilo "revisione": il modello ha emesso l'intero output (frontmatter
// YAML + markdown body + callout) racchiuso in un singolo ```yaml ... ```.
// Effetto collaterale: il renderer PDF (e qualsiasi viewer markdown) tratta
// tutto come code block monospace invece di renderizzare heading/callout/tabelle.
//
// Heuristic conservativa:
//   - se la prima riga "significativa" (dopo trim whitespace) inizia con "```"
//   - E l'ultima riga "significativa" è esattamente "```"
//   - allora strippa la prima e ultima linea, ritorna il contenuto interno
//
// Negli altri casi (no fence, fence solo all'inizio senza chiusura, fence
// interno come esempio di codice nel mezzo del documento) ritorna il body
// invariato. Niente strip se il pattern non è chiaro — better safe than rovinare
// output legittimi.
func stripFenceWrapping(body string) string {
	trimmed := strings.TrimSpace(body)
	if !strings.HasPrefix(trimmed, "```") {
		return body
	}
	// Trova la fine della prima riga (separatore prima del contenuto)
	nl := strings.IndexByte(trimmed, '\n')
	if nl < 0 {
		return body // body è una singola riga "```..." malformata
	}
	// Verifica che termini con "```"
	if !strings.HasSuffix(trimmed, "```") {
		return body
	}
	// Estrai contenuto interno: dopo prima newline, prima dell'ultimo "```"
	inner := trimmed[nl+1:]
	inner = strings.TrimSuffix(inner, "```")
	return strings.TrimRight(inner, "\n ")
}
