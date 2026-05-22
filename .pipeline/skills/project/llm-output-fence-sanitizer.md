---
name: llm-output-fence-sanitizer
date: 2026-05-22
source: Sprint 1.5.D — Qwen3.5-122B wrap output in ```yaml
tags: [llm, sanitize, output, markdown, rendering]
---

# LLM output fence wrapping sanitizer

## Problem

Alcuni modelli LLM (osservato Qwen3.5-122B su profile `revisione`, 2026-05-22T11:17) emettono l'intero output racchiuso in un singolo fenced code block:

```
```yaml
---
profilo: revisione
---
# Heading
> [!INFO] Box
> Body content
```
```

Effetto: **tutti i viewer markdown (Obsidian, GitHub, PDF renderer)** trattano correttamente il body come code block monospace, perdendo heading/callout/tabelle/formatting. Il PDF generato mostra il MD raw come blocco di Courier 9pt invece del documento impaginato.

## Detection

Pattern conservativo lato runner Go (no false-positive su body legittimo che ha code block in mezzo):

```go
func stripFenceWrapping(body string) string {
    trimmed := strings.TrimSpace(body)
    if !strings.HasPrefix(trimmed, "```") {
        return body
    }
    nl := strings.IndexByte(trimmed, '\n')
    if nl < 0 || !strings.HasSuffix(trimmed, "```") {
        return body
    }
    inner := trimmed[nl+1:]
    inner = strings.TrimSuffix(inner, "```")
    return strings.TrimRight(inner, "\n ")
}
```

**Regola**: prima riga `\`\`\`<lang>?` + ultima riga `\`\`\`` → strip. Altrimenti no-op.

Casi corretti (no-op):
- Body senza fence in cima
- Body con fence interno (code block legittimo nel mezzo)
- Body con solo opening fence (body troncato)

## Application

Sanitize una sola volta nel runner, **PRIMA** di scrivere qualsiasi output (`.md`, `.pdf`, `.docx`, `.log`). Pattern:

```go
body, stato, err := callProviderWithStreaming(...)
if err != nil { return nil, err }
if cleaned := stripFenceWrapping(body); cleaned != body {
    log.Info("body sanitization: rimosso fence wrapping spurio")
    body = cleaned
}
// poi writeOutput(.md), maybeWritePDF, maybeWriteDocx con `body` sanitized
```

## Test rosso pattern

```go
func TestStripFenceWrapping_DetectsAndStripsYamlWrap(t *testing.T) {
    wrapped := "```yaml\n# Titolo\n```"
    if got := stripFenceWrapping(wrapped); got != "# Titolo" { ... }
}
func TestStripFenceWrapping_NoOpWhenFenceIsInternalCode(t *testing.T) {
    body := "# Titolo\n\nEsempio:\n\n```go\nx := 1\n```\n\nFine."
    if got := stripFenceWrapping(body); got != body { ... }
}
```

## Quando NON applicare strip

- Output **legittimamente** pure-code (es. capability che genera JSON puro o codice TOML). Non è il nostro caso (tutti i 7 profili Sprint 1 producono MD con frontmatter). Se in futuro un profilo produce solo code → flag opt-out `[output] skip_sanitize = true` nel profilo TOML.
- Detection failed positive (mai osservato in 30+ run reali, pattern conservativo).

## Anti-pattern: prompt engineering "per evitare il wrap"

Aggiungere al system prompt "NON racchiudere l'output in fenced code block" è fragile:
- Modello può ignorare l'istruzione (specialmente su context lunghi)
- Aumenta tokens system prompt
- Non garantisce su future versioni del modello

Sanitize **lato code** è defense-in-depth indipendente dal modello. Combinare entrambe le strategie va bene, ma la sanitize Go è quella che dà la garanzia.
