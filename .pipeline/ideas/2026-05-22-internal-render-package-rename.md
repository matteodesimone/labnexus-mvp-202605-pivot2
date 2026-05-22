---
date: 2026-05-22
status: backlog
source: compound learning Sprint 1.5.D
tags: [refactor, naming, cleanup]
priority: low
effort: small
---

# Rinomina `internal/pdf/` → `internal/render/`

## Idea

Il package `internal/pdf/` Sprint 1.5.D contiene `Render()` (PDF) **E** `RenderDocx()` (DOCX). Il naming è fuorviante: la sottocartella suggerisce "solo PDF" mentre è il package di rendering documenti in generale.

## Quando

Quando aggiungiamo il terzo formato output (probabilmente HTML in Sprint 2 per visualizzazione web/preview), il rename diventa necessario. **NON farlo prima** — single-feature refactor non vale il blast radius isolato.

## Cosa cambia

- Rinomina directory `internal/pdf/` → `internal/render/`
- Aggiorna package declaration in tutti i file Go (`package pdf` → `package render`)
- Aggiorna import path nei consumer (`runner/pdf_writer.go`, eventuale altro)
- `pdf.Metadata` → `render.Metadata`, `pdf.BrandConfig` → `render.BrandConfig`, `pdf.Render` → `render.RenderPDF`, `pdf.RenderDocx` → `render.RenderDocx`
- Aggiorna nomi test e helper (`pdf_test.go` → `render_test.go`, `pdf_writer.go` → `output_writer.go` o simile)
- Aggiorna eventuali skill o solution che referenziano `internal/pdf/`

## Effort

~30-45 min di lavoro meccanico (Go ha tooling solido per il rename). Test devono restare verdi dopo.

## Trigger

Aprire questa idea quando il prossimo formato di output (HTML, EPUB, ecc.) viene proposto. Per ora archiviata in backlog.
