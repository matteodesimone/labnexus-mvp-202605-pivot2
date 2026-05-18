# Materiali di dominio del cliente — KB, visione, golden file dei test precedenti

Questa cartella contiene il materiale del cliente. **Va trattato in modo differente a seconda del tipo**: alcuni file sono contesto di expertise da iniettare come system prompt; altri sono visione strategica di lungo periodo; altri ancora sono golden file di confronto per i primi due profili. Vedi sotto per gruppo.

---

## 1. `KB-ispettore/` — Expertise ispettiva di Denis (system context per i profili)

**Cosa è**: la KnowledgeBase scritta da Denis Brazzo (esperto SGQ + ispettore Accredia). È **expertise di dominio**, scritta da lui in forma di prompt per LLM. Copre ISO/IEC 17025 in modo strutturato.

**Come usarla**: i profili dell'eseguibile LabNexus la **iniettano come system context** prima di chiamare Qwen. Ogni file `profili/<nome>.yml` specifica nel campo `kb_files:` quali file di questa cartella concatenare. I file sono trattati come prompt normali (non slash command), nell'ordine indicato.

**Struttura**:

- `CLAUDE.md` — definisce identità, tono, regole base dell'agente. Va sempre incluso per primo nei profili
- `come-pensa-un-ispettore.md` — framework di lettura ispettiva. Centrale per quasi tutti i profili
- `MAPPA_ESPLOSA_REQUISITI_ISO17025.md` — mappa dei requisiti normativi, sezione per sezione
- `sezioni-ISO/` — file per sezione ISO con le 3 domande ispettive per requisito
- `NC-patterns/` — pattern di gestione delle non-conformità, errori fatali da evitare
- `INPUT_esempio/` — **cartella di scaffolding originale di Denis**. Non è la struttura reale del SGQ del cliente: era una cartella temporanea che Denis ha usato per passare i file di test. **I file rilevanti sono già stati copiati in `test-precedenti/`** (vedi sotto). Questa cartella si può ignorare nel flusso di sviluppo, è preservata solo per integrità della KB originale
- `bozze/` — output del Test 1 di Denis. **Già copiato in `test-precedenti/test-1-output-approvato/`**

**Importante**: la KB è scritta da Denis e va trattata come fonte di verità per quanto riguarda il tono ispettivo e i requisiti normativi. Il CTO può **selezionare** quali file iniettare per profilo, ma non riscriverli. Le iterazioni sulla KB le fa Denis sulla sua versione locale; questa è una copia controllata al momento dell'avvio dello sprint.

**Per il CTO**: a fine sprint andrà allineata con la versione che Denis avrà nel frattempo modificato. Per lo sviluppo dei profili, **questa versione è la fonte canonica**.

---

## 2. `LabNexus_CoWork_KB.md` e `LabNexus_WikiLLM_Orchestrator.md` — Visione del prodotto finale

**Cosa sono**: due documenti scritti da Denis che descrivono **il sistema AICertus completo**, l'orizzonte di 12-18 mesi. Contengono:

- Pattern di interazione *osserva → propone → approva → registra*
- Lista di 11 agenti specializzati (DocAgent, NormAgent, CAPAAgent, ReviewAgent, AuditAgent, EquipmentAgent, CompetenceAgent, SupplierAgent, PTQCAgent, DataIntegrityAgent, RMCRM_Agent)
- 8 stress test del sistema
- Comandi slash `/sgq-*` (ingest, update, audit, draft, score, impact, capa, equipment, pt, data-integrity, review, watch, risk-update, 17034-rm)
- Struttura cartella `_LabNexus_Output/`
- Template per CAPA Pack, Equipment Alert, Management Review Pack, ecc.

**Come usarli**:

- **Per lo Sprint 1: non sono spec da implementare**, sono visione di lungo periodo. Lo Sprint 1 implementa solo i task LLM-critical del sistema (le 7 capability), non l'orchestrazione automatica né la pipeline deterministica.
- **Per la scrittura dei trigger_prompt dei profili**: sono **fonte autorevole per il tono e il template di output** di ciascuna capability. Es. quando si scrive il `trigger_prompt` del profilo `review-pack`, si fa riferimento al `TEMPLATE_Management_Review_Pack` (sezione 15.4 del CoWork KB) come specifica del formato di output atteso.
- **Per il prompt meta** (`feat-meta`): sono fonte autorevole sulla tassonomia degli agenti specializzati di Denis, utile quando Claude propone profili nuovi.
- **Per Sprint 2 e successivi**: diventano la roadmap del sistema completo. Conservare come riferimento.

**Non confondere**: la maggior parte di ciò che è descritto in questi due documenti (`/sgq-watch`, `/sgq-index`, pattern CoWork, watcher di cartella) è **roba di Sprint 2 e successivi**, non di questo sprint.

---

## 3. `test-precedenti/` — Golden file per gli shakedown del motore

**Cosa è**: il materiale dei Test 1 e Test 2 condotti manualmente in chat Ollama dal CTO **prima dell'avvio dello sviluppo del motore**. Entrambi i test hanno ricevuto **OK formale dal cliente**.

**Come usarlo**: i profili `revisione` e `rilievi` dell'eseguibile devono **replicare questi due test** quando eseguiti sull'eseguibile. Gli output prodotti dall'eseguibile vanno confrontati con i golden file qui contenuti:

- Se l'output dell'eseguibile è qualitativamente paragonabile al golden → motore validato per quella capability
- Se l'output è significativamente peggiore → c'è un problema in pipeline software (prompt assembly, context injection, streaming). Da debuggare prima di proseguire con i profili successivi

**Cosa c'è dentro**:

- `test-1-input/` — i 3 file di input del Test 1: `PG_RISK_LAB_Rev_00` (documento da revisionare, in 2 formati .docx e .txt), `DE0779_RT_08rev03.pdf` (norma vecchia), `RT-08-rev.05-_EC-10-02-2022.pdf` (norma nuova)
- `test-1-output-approvato/` — `PG_RISK_LAB_Rev_01_BOZZA-qwen3.6-locale.docx` (output di Qwen approvato da Denis) + versione markdown leggibile
- `test-2-input/` — `ACIAA A1 rilievi.csv` (lista rilievi reali)
- `test-2-output/` — **vuota**: il CTO non ha conservato il file di output del Test 2 in formato condivisibile. La validazione del profilo `rilievi` userà criteri qualitativi descritti nel brief tecnico (sezione `feat-003`), confrontati con il template `CAPA_Pack` del CoWork KB sezione 15.2

**Importante**: gli shakedown (profili 1 e 2) servono a validare il **motore**, non il modello. Le Capability A e B sono già state validate sul modello dai Test 1 e 2 originali. Se l'eseguibile riproduce i risultati, il software è corretto e si procede con i profili nuovi.
