---
date: 2026-05-22
status: backlog
source: compound learning Sprint 1.5.D
tags: [workflow, denis, docx, pdf, l2-validation]
priority: medium
effort: medium
---

# DOCX → PDF post-L2-validation

## Idea

Workflow completo Denis post-bozza:

1. labnexus produce `.md`, `.pdf`, `.docx` accoppiati (Sprint 1.5.D)
2. Denis **edita** il `.docx` in Word: track changes, comments, accept/reject, modifiche manuali
3. Denis salva il `.docx` validato come `<basename>_validato.docx`
4. **Nuovo comando** `labnexus finalize <docx-path>` converte il DOCX validato in PDF finale per consegna formale al cliente / archiviazione ISO 17025

## Motivazione

Lo Sprint 1.5.D produce 4 output per run: `.md`, `.log`, `.pdf` (consegna), `.docx` (L2 edit). Il flusso "naturale" è: Denis lavora in Word (DOCX) → PDF finale per il cliente. Ma il PDF generato è quello della bozza, non dell'edit di Denis. Mancanza:
- Il `.pdf` accoppiato è la bozza pre-validation
- Manca il PDF "validato" = quello che Denis effettivamente consegna

## Implementazione proposta

Nuovo subcommand `labnexus finalize <input.docx> [--output <out.pdf>]`:
- Pipeline: pandoc input.docx → typst markup → typst compile → PDF
- Usa stesso `template.typ` LabNexus
- Frontmatter del DOCX preservato (date, validatore, versione)
- Aggiunge nel footer: "Validato da {validatore}" (Denis edita un campo Word custom o usa Word document properties)
- Marker visivo "DEFINITIVO" o "BOZZA VALIDATA L2" nell'header

## Sfide

- DOCX → Typst markup via pandoc: la qualità della conversione dipende da quanto Word ha "sporcato" lo styling durante l'edit (font specifici, colori inline, ecc.)
- Track changes residue: se Denis non ha accepted/rejected tutti i changes, il DOCX contiene revisioni. Va validato/strippato prima della conversione
- Identificazione "validatore": Word file properties `lastModifiedBy` o campo custom?

## Effort stimato

~3-5h per implementazione baseline + test su DOCX reali di Denis.

## Trigger

Quando Denis avrà completato il primo ciclo L2 reale e ci saranno i `.docx` validati su cui iterare. Per ora backlog.
