---
codice: 00_INDICE
tipo: indice
livello: 0
titolo: "Indice della KB-Ispettore LABNEXUS"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
sempre_in_contesto: true
note: "Questo file è caricato dall'agente all'avvio di ogni sessione. È la mappa di consultazione di tutta la KB."
---

# Indice della KB-Ispettore — LABNEXUS v2.1

> Questa è la **Knowledge Base base** del prodotto LABNEXUS per laboratori di prova accreditati UNI CEI EN ISO/IEC 17025:2018. È una KB **invariabile dal laboratorio**: contiene la norma, le prescrizioni Accredia e la prassi ispettiva. Le specificità del singolo laboratorio (metodi, matrici, parchi strumenti, procedure proprie) appartengono alla KB verticale, separata.

## 1. Cosa carica sempre l'agente

All'avvio l'agente carica **sempre** in contesto questi file:

- `00_INDICE.md` (questo file)
- `00_FONTI_NORMATIVE.md` — registro chiuso dei documenti citabili
- `00_GLOSSARIO.md` — termini ISO/Accredia con definizione e § sorgente
- tutta la cartella `01_SYSTEM/` (persona, protocolli ASK-BEFORE-ANSWER, HANDOFF-QM, CITAZIONI, TRE LIVELLI DI CONFIDENZA, GOVERNANCE-AI)

## 2. Cosa carica on-demand l'agente

In base alla richiesta del QM, l'agente carica selettivamente:

- file in `02_LOGICA_ISPETTIVA/` quando entra in modalità analisi/valutazione
- uno o più file in `03_REQUISITI/` corrispondenti al/ai § ISO/RT-08 in oggetto
- uno o più file in `04_APPENDICI/` quando il tema è trasversale (validazione, incertezza, riferibilità, CRM, PT, audit interni)
- `05_ORCHESTRATORE/LabNexus_CoWork_KB.md` quando opera in modalità agentica sul SGQ

## 3. Mappa file completa (52 file .md: 51 di KB + CHANGELOG)

### 3.1 Livello 0 — sempre in contesto

| File | Scopo |
|---|---|
| `00_INDICE.md` | Questo file |
| `00_FONTI_NORMATIVE.md` | Registro chiuso fonti citabili (norme + Accredia + best practice) |
| `00_GLOSSARIO.md` | Glossario tecnico ISO/Accredia |

### 3.2 `01_SYSTEM/` — comportamento dell'agente

| File | Scopo |
|---|---|
| `01_PERSONA_aicertus.md` | Identità, tono e ruolo dell'agente LabNexus |
| `02_PROTOCOLLO_ASK_BEFORE_ANSWER.md` | Quando l'agente DEVE chiedere prima di rispondere |
| `03_PROTOCOLLO_HANDOFF_QM.md` | Come ogni output dell'agente si chiude con handoff al QM |
| `04_PROTOCOLLO_CITAZIONI.md` | Regole formali di citazione norma/RT/RG/DT |
| `05_TRE_LIVELLI_CONFIDENZA.md` | Lettura solida / prudenziale / da confermare |
| `06_GOVERNANCE_AI.md` | Governance dell'AI nel laboratorio: deployment locale, riservatezza, prompting, audit AI-driven, mappa ai § ISO |

### 3.3 `02_LOGICA_ISPETTIVA/` — modalità analisi

| File | Scopo |
|---|---|
| `01_COME_PENSA_ISPETTORE.md` | Logica e mentalità di un ispettore Accredia |
| `02_TRE_LIVELLI_DI_CONTROLLO.md` | Pattern meccanismo → estensione → efficacia |
| `03_RISPOSTA_NC.md` | Come si risponde a una NC Accredia in modo solido |
| `04_ERRORI_FATALI.md` | Errori che generano NC grave o sospensione |

### 3.4 `03_REQUISITI/` — un file per sotto-sezione norma

| File | § ISO/IEC 17025:2018 |
|---|---|
| `00_MAPPA_ESPLOSA.md` | Indice puro che linka tutti i file requisito |
| `sez_4_imparzialita_riservatezza.md` | § 4 |
| `sez_5_strutturali.md` | § 5 |
| `sez_6_1_risorse_generali.md` | § 6.1 |
| `sez_6_2_personale.md` | § 6.2 |
| `sez_6_3_strutture_ambiente.md` | § 6.3 |
| `sez_6_4_dotazioni.md` | § 6.4 |
| `sez_6_5_riferibilita.md` | § 6.5 |
| `sez_6_6_prodotti_servizi_esterni.md` | § 6.6 |
| `sez_7_1_riesame_richieste.md` | § 7.1 |
| `sez_7_2_selezione_verifica_validazione_metodi.md` | § 7.2 |
| `sez_7_3_campionamento.md` | § 7.3 |
| `sez_7_4_manipolazione_oggetti.md` | § 7.4 |
| `sez_7_5_registrazioni_tecniche.md` | § 7.5 |
| `sez_7_6_incertezza.md` | § 7.6 |
| `sez_7_7_assicurazione_validita.md` | § 7.7 |
| `sez_7_8_presentazione_risultati.md` | § 7.8 |
| `sez_7_9_reclami.md` | § 7.9 |
| `sez_7_10_attivita_non_conformi.md` | § 7.10 |
| `sez_7_11_controllo_dati.md` | § 7.11 |
| `sez_8_1_opzioni_AB.md` | § 8.1 |
| `sez_8_2_documentazione_sgq.md` | § 8.2 |
| `sez_8_3_controllo_documenti.md` | § 8.3 |
| `sez_8_4_controllo_registrazioni.md` | § 8.4 |
| `sez_8_5_rischi_opportunita.md` | § 8.5 |
| `sez_8_6_miglioramento.md` | § 8.6 |
| `sez_8_7_azioni_correttive.md` | § 8.7 |
| `sez_8_8_audit_interni.md` | § 8.8 |
| `sez_8_9_riesame_direzione.md` | § 8.9 |

### 3.5 `04_APPENDICI/` — tematiche trasversali

| File | Tema |
|---|---|
| `A1_validazione_metodi.md` | Verifica e validazione metodi (best practice + § 7.2) |
| `A2_incertezza_misura.md` | Valutazione incertezza (GUM, Eurachem, Nordtest) |
| `A3_riferibilita_taratura.md` | Riferibilità metrologica, scelta laboratorio di taratura |
| `A4_materiali_riferimento_CRM.md` | CRM, materiali di riferimento, ISO 17034 |
| `A5_prove_valutative_PT_ILC.md` | PT e ILC, piano partecipazione, valutazione esiti |
| `A6_audit_interni_riesame_direzione.md` | Conduzione audit interni e riesame di direzione |
| `A7_campionamento.md` | Campionamento accreditato e incertezza da campionamento |
| `A8_gestione_NC_da_accreditatore.md` | Gestione operativa NC ricevute da Accredia |

### 3.6 `05_ORCHESTRATORE/`

| File | Scopo |
|---|---|
| `LabNexus_CoWork_KB.md` | Architettura agentica per SGQ vivo (trigger, output standard, governance) |

## 4. Tassonomia dei contenuti (livelli di citabilità)

Ogni contenuto della KB è marcato con il proprio **livello** nel frontmatter del file e/o nelle intestazioni di sezione:

| Livello | Significato | Citazione obbligatoria nell'output |
|---|---|---|
| **1 — Fatto normativo** | Testo / requisito di norma ISO/IEC 17025:2018 | doc + § + (se rinvio) "consulta la norma" |
| **2 — Prescrizione Accredia** | RT-08, RG-02, RG-09, RT-23, RT-26, RT-39, DT-XX, politiche, raccomandazioni CdIG | doc + rev + § + data |
| **3 — Prassi ispettiva** | Sapere esperienziale di Denis Brazzo, ispettore Accredia | dichiarata come tale; non requisito normativo |
| **4 — Best practice tecnica** | Eurachem, CITAC, Nordtest, EUROLAB, ISTISAN, EU regs, JCGM, JRC | citabile come raccomandazione, mai come obbligo |

L'agente, quando risponde, **non può presentare** un contenuto di livello 3 o 4 come se fosse di livello 1 o 2. È un vincolo non negoziabile della persona `aicertus` (vedi `01_SYSTEM/01_PERSONA_aicertus.md`).

## 5. Convenzioni di scrittura

- Tutti i file hanno frontmatter YAML uniforme (`codice`, `tipo`, `livello`, `sezione_norma`, `fonti_primarie`, `fonti_secondarie`, `file_collegati`, `tags`).
- I link interni alla KB usano la sintassi `[[codice_file]]` (es. `[[sez_6_4_dotazioni]]`, `[[A2_incertezza_misura]]`).
- I link a documenti della **KB verticale** del singolo laboratorio (es. `[[Piano_PT_2026]]`, `[[Matrice_Competenze]]`, `[[NC-2026-001]]`, `[[TEMPLATE_*]]`, `[[MOD_*]]`, `[[procedura-controllo-dati]]`) sono **placeholder** che esistono solo nel SGQ del laboratorio dopo l'onboarding. Quando l'agente li incontra in questa KB base, sa che sono riferimenti alla KB verticale, non a file di questa KB.
- I codici di rinvio normativo (`[[ISO17025-7.7]]`, `[[RT-08-7.7]]`) sono identificativi semantici di paragrafi, non file: l'agente li risolve consultando `00_FONTI_NORMATIVE.md`.
- Le citazioni di norma seguono il `04_PROTOCOLLO_CITAZIONI.md`.
- Le date sono in formato ISO `YYYY-MM-DD`.

## 6. Cosa NON è in questa KB

- **KB verticale del laboratorio** — file/cartelle separate `KB-laboratorio/` che il lab costruisce in onboarding e mantiene
- **ISO/IEC 17034** (RM producers) e **ISO/IEC 17020** (ispezione) — fuori scope (citabili come riferimento esterno se rilevante)
- **UNI EN ISO 15189** (lab medici) — fuori scope (LS-04 § 4)
- **ISO/IEC 17043** (PT providers) — fuori scope come accreditamento (citabile come riferimento per PT)
- **Procedura cliente "Gestione AI in laboratorio"** — co-scritta in onboarding, vive nella KB verticale. NB: i **principi generici** di governance AI (deployment locale, riservatezza, prompting, audit AI-driven, mappa ai § ISO) sono invece in questa KB base, in `01_SYSTEM/06_GOVERNANCE_AI.md`; la procedura verticale ne è la calata operativa nel singolo laboratorio.

## 7. Storico versioni

- **v2.1 — 2026-05-29** — Audit e finalizzazione. Normalizzati 277 wikilink alla forma `[[codice]]`. Corrette etichette § 5.1/5.2 (mappa esplosa) e un § inesistente (§ 4.3 → § 7.10 in A5). Nuovo file `01_SYSTEM/06_GOVERNANCE_AI.md` (strato AI-governance generico). `sez_7_2` portato a gold-standard con ponte validazione↔incertezza e ILAC-G17:01/2021. Aggiunto ILAC-G17 al registro fonti, con nota su EA-4/16 ritirata. Marcatori 🟢🟡🟣 e accenti uniformati. Vedi `CHANGELOG.md`.
- **v2.0 — 2026-05-28** — Refactor integrale post-audit. Sezione 8 ricostruita su indice ufficiale ISO 17025:2018 / RT-08 rev.05. Registro fonti completo (LS-04 rev.20). Tassonomia a 4 livelli. Pattern di file requisito uniforme. Reclami rimpatriati in § 7.9. Nuove appendici trasversali.
- **v1.0 — 2026-05-06** — Versione iniziale (vedi backup `backup_KB_pre_refactor_2026-05-28.zip`).
