# Skill: Post-Pivot Doc Audit Checklist

**Created**: 2026-05-21 in sub-fetta 1.5.A (compound)
**Trigger**: ogni volta che un pivot inverte una decisione architetturale o semantica già documentata (es. locale→cloud, opt-in→opt-out, framework-A→framework-B, file-format-X→file-format-Y).

## Il problema

`/v-implement` aggiorna i file di docs menzionati nel plan. Ma il plan tipicamente menziona solo il sottoinsieme di docs in cui il CTO ha pensato quando ha scritto il plan. Dopo grandi pivot semantici, il risultato è:

- Layer 3 standards (`02-project.md`) non aggiornato → ogni futura `/v-spec`/`/v-plan` legge info errate
- README intro contraddice la propria tabella sotto
- CONTEXT.md "Single source of truth per LLM agent" obsoleto
- USER-GUIDE.md fa false affermazioni di privacy che gli utenti leggono per decidere come usare lo strumento

Tutti questi sono violazioni del principio "la documentazione è codice" (Layer 0). Il /v-review li cattura come HIGH findings → loop 2 di fix necessario → costo tempo che si poteva risparmiare con un audit pre-emptive.

## Il pattern

Prima di concludere `/v-implement` (o durante il review come safety check):

1. **Enumera i pattern testuali obsoleti** che il pivot ha invalidato. Esempi:
   - Inversione provider locale→cloud: `"in locale"`, `"locale: true"`, `"non esce dalla macchina"`, `"nessun dato"` + nome del vecchio provider
   - Rimozione di un componente: nome del componente, suoi file specifici, suoi pattern dichiarativi
   - Cambio di un default: il vecchio default citato come tale in docs

2. **Grep su tutte le superfici dichiarative**:
   ```bash
   grep -rn "<pattern>" \
     .pipeline/standards/ \
     docs/ \
     README.md \
     CONTEXT.md \
     scripts/ \
     profili/      # o equivalente "schema by-product"
   ```

3. **Classifica ogni hit**:
   - **OUTDATED**: va aggiornato (es. README intro)
   - **HISTORY**: va riformulato come storia esplicita (es. "Sprint 1: ... [decommissionato in Sprint 1.5]")
   - **AUDIT-TRAIL BY-PRODUCT**: marker generato dal codice (frontmatter, log labels, headers HTTP) che era veritiero pre-pivot e ora è menzognero — va rimosso o reso runtime-dinamico
   - **REMAIN VALID**: il hit è in un context dove la frase originale resta corretta (es. test che valida che il vecchio pattern non sia più presente — quel test deve contenere la stringa per asserirla)

4. **Applica i fix prima del /v-review**. Costo: 5-10 min. Risparmio: 1 loop di review-fix con stesso budget.

## Anti-pattern

- **"Aggiorno solo i 5 docs menzionati nel plan"**: il plan è scritto prima di sapere tutto. I docs sono molte più superfici di quante tu pensi a plan-time. Grep > memoria.
- **"Il review esterno lo troverà comunque"**: vero, ma costa un loop intero. Inoltre quando i reviewer esterni sono degraded (quota, timeout), perdi anche quel safety net.
- **"`locale: true` nel frontmatter è un valore statico, non importa"**: importa. Se il frontmatter è un audit trail (chi-quando-dove ha generato), un valore statico che era veritiero pre-pivot diventa menzognero post-pivot. Privacy by design implica che gli audit trail siano veritieri.

## Precedente concreto

Sub-fetta 1.5.A (Sprint 1.5 LabNexus) — pivot 3 invertiva NFR-1 (locale → cloud EU-GDPR default). Il /v-implement aveva aggiornato 5 docs (README, USER-GUIDE, DEV-GUIDE, CONTEXT, Makefile). Il /v-review loop 1 ha trovato 5 HIGH findings di doc-drift:
- `.pipeline/standards/02-project.md` Layer 3 non toccato
- README riga 3 contraddice la sua stessa tabella
- CONTEXT righe 9 + 102 obsolete
- USER-GUIDE "Nessun dato esce dalla macchina" falsa
- 7 profili con `locale: true` (audit-trail by-product menzognero)

Tutti fixabili pre-emptively con un grep audit di 5 minuti. Costati invece ~45 min in loop 2 della review.

Vedi `.pipeline/solutions/2026-05-21-sub-fetta-1.5.A-pivot-3-cli-puro.md` per il caso completo.

## Estensione Sprint 1.5.B: pattern semantici, non solo sintattici

Lesson da sub-fetta 1.5.B compound (2026-05-21): la versione iniziale della skill catturava solo pattern SINTATTICI (`.yml`, `osascript`, `.app`). Ma il pivot semantico può lasciare residui nei docs come **affermazioni inverse** del nuovo stato, non solo riferimenti a nomi obsoleti.

Esempio 1.5.B: meta-prompt riga 284 diceva _"NON suggerire `provider: eurouter` per dati reali di Denis (deve passare per gate cloud + approvazione esplicita; **default ollama**)"_. Nessuna delle stringhe `.yml`/`ollama`/`gate` da sole è "obsoleta" sintatticamente — è l'**affermazione semantica complessiva** che è invertita post-pivot-3.

### Checklist semantica estesa

Per ogni pivot, oltre ai pattern sintattici, enumera anche:

1. **Affermazioni di stato precedente** che ora sono FALSE. Esempi pivot-3:
   - "non esce dalla macchina"
   - "default ollama"
   - "Ollama-only nel deliverable"
   - "in locale"
   - "no cloud"
   - "approvazione esplicita necessaria per eurouter"

2. **Affermazioni di direzione precedente** che ora sono INVERTITE. Esempi pivot-3:
   - "NON suggerire eurouter"
   - "evita cloud per dati reali"
   - "preferisci ollama"

3. **Marker by-product hardcoded** (frontmatter, header HTTP, log labels, env defaults) che riflettono il vecchio stato. Esempi 1.5.A:
   - `locale: true` nei frontmatter
   - `provider: ollama` come hardcoded default
   - Comment di funzioni che descrive precedence vecchia

### Grep pattern semantici (esempio)

```bash
# Per pivot locale→cloud:
grep -rni "no cloud\|non esce\|in locale\|default ollama\|preferisci.*ollama\|approvazione.*eurouter" \
  docs/ README.md CONTEXT.md .pipeline/standards/

# Per pivot YAML→TOML:
grep -rni "schema yaml\|formato yaml\|file yaml\|\.yml$" \
  docs/ README.md CONTEXT.md .pipeline/standards/
```

I pattern semantici sono PROGETTO-SPECIFIC: devi enumerarli a mano sulla base della spec amendment. La skill non può proporli automaticamente, ma deve ricordarti DI farli.

## Lesson 1.5.B integration nel /v-implement

Apply pre-emptive: durante step 8 del plan 1.5.B "Update docs", esegui BOTH grep sintattici E semantici PRIMA del /v-review. Il review loop 1 di 1.5.B ha trovato 5 HIGH doc-drift che si potevano catturare con grep semantici (es. "default ollama" in meta-prompt). Costo: 5 min in più al implement. Risparmio: ~30 min di loop 2 fix-and-verify.

## Quando NON applicare

- Bugfix di portata limitata (1-2 file, semantica invariata): doc audit grep è overkill
- Aggiunta di funzionalità isolata senza modifica di decisioni architetturali: i docs in scope sono enumerabili dal plan
- Refactor puramente strutturale (rinomina interna senza impatto su superfici dichiarative): la signature pubblica resta invariata, niente da auditare

## Generalizzabilità (verso framework promotion)

Questa skill è candidata a framework promotion (universal): il pattern doc-as-code post-pivot vale per qualsiasi stack/dominio. Vedi `.pipeline/proposed-updates/` se promossa.
