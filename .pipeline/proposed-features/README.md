# Proposed Verification Scenarios — Sprint 1 LabNexus

Bozza di scenari Gherkin (godog) per gli FR e gli EC della spec. **Sono proposte**: vanno revisionati e approvati durante il gate `/v-spec`.

Convenzioni:
- Sintassi Gherkin in italiano (`# language: it`) — coerente con il dominio del cliente (Denis legge le scenario).
- Un file per gruppo coerente di FR. Edge case raggruppati in `edge-cases.feature`.
- Linguaggio "user-side": gli step descrivono *cosa fa l'utente* e *cosa osserva*, non *come è implementato*.
- Dati concreti: nomi file reali dei Test 1/2 dove possibile, sintetici plausibili altrove.
- Step implementations (per godog) verranno scritte in `/v-test-scaffold`. Qui c'è solo la specifica behaviorale.

## File

- `motore-cli.feature` — FR-1, FR-3, FR-4, FR-5, FR-6, FR-12 (CLI, schema profilo, parsing input, composizione prompt, gestione token, env var EUrouter).
- `motore-providers.feature` — FR-7, FR-8 (provider Ollama NDJSON e EUrouter SSE, progress bar, streaming).
- `motore-output.feature` — FR-9 (file markdown con frontmatter).
- `motore-modo-interattivo.feature` — FR-10, FR-11 (dialog osascript, bundle `.app`).
- `profili.feature` — FR-13…FR-19 (i 7 profili — uno scenario macro per ognuno).
- `prompt-meta.feature` — FR-20, FR-21, FR-22 (il deliverable + i 3 casi di validazione).
- `edge-cases.feature` — EC-1…EC-15 (selezione dei più rilevanti).

## Mapping FR → file

| FR | File |
|---|---|
| FR-1 motore CLI subcommands & flag richiesti | motore-cli.feature |
| FR-2 risoluzione path | motore-cli.feature |
| FR-3 schema profilo + `validate` | motore-cli.feature |
| FR-4 parsing input multi-formato | motore-cli.feature |
| FR-5 composizione prompt | motore-cli.feature |
| FR-6 stima token + warning/errore | motore-cli.feature |
| FR-7 LLMProvider interface + 2 impl | motore-providers.feature |
| FR-8 progress bar + log per step | motore-providers.feature |
| FR-9 output markdown + frontmatter | motore-output.feature |
| FR-10 dialog macOS interattivo | motore-modo-interattivo.feature |
| FR-11 bundle `.app` | motore-modo-interattivo.feature |
| FR-12 EUrouter env var check | motore-cli.feature |
| FR-13…FR-19 i 7 profili | profili.feature |
| FR-20 esistenza meta-prompt | prompt-meta.feature |
| FR-21 guida d'uso Denis | prompt-meta.feature |
| FR-22 3 casi validazione meta | prompt-meta.feature |
