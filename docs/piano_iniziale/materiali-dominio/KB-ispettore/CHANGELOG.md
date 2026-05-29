---
codice: CHANGELOG
tipo: cambiamenti
livello: 0
titolo: "Changelog KB-Ispettore LABNEXUS"
data_ultima_revisione: "2026-05-29"
---

# Changelog — KB-Ispettore LABNEXUS

## v2.1 — 2026-05-29 (audit e finalizzazione)

Pass di audit sui quattro assi (correttezza normativa, coerenza interna, ottimizzazione LLM, completezza) e prime modifiche approvate dal QM/Ispettore. Backup integrale pre-modifica conservato.

### Coerenza interna e wikilink
- **Normalizzati 277 wikilink** alla forma canonica `[[codice]]` (rimossi tutti i prefissi `../`, `./`, `sottocartella/`), come prescritto da `00_INDICE` § 5 e `04_PROTOCOLLO_CITAZIONI` § 9. Verifica programmatica: **zero link interni con path residui**.
- **Rettifica del claim v2.0**: la nota "zero refusi residui (~26 link verticale)" era inesatta — coesistevano 277 link interni in forma non canonica con path. Sono stati normalizzati ora; i ~25 placeholder della KB verticale (`TEMPLATE_*`, `MOD_*`, `Piano_PT_2026`, ecc.) restano corretti per definizione.

### Correttezza normativa
- **§ 5.1 / 5.2** in `00_MAPPA_ESPLOSA`: corretta l'inversione delle etichette (5.1 = entità giuridica; 5.2 = identificazione del management), riallineata al testo della norma e al file di dettaglio `sez_5_strutturali`. Affinate anche le etichette 5.3–5.7.
- **§ 4.3 inesistente**: in `A5_prove_valutative_PT_ILC` il rinvio "workflow ex § 4.3" è stato corretto in **§ 7.10 (attività non conformi)** — in ISO 17025:2018 il § 4 contiene solo 4.1 e 4.2.

### Contenuti — gold standard e profondità (focus: incertezza + validazione)
- **`sez_7_2`** portato a **gold-standard**: marcatori di confidenza standardizzati, nuova sezione **B-bis "Dal dato di validazione all'incertezza"** (ponte § 7.2.2.3 ↔ § 7.6, approccio top-down), nuovo errore fatale EF-7.2.2-E, sezione **G "Ruolo dell'agente AI"** (governance per-requisito).
- **Propagazione del ponte validazione↔incertezza** ai file collegati: `sez_7_6` (callout B-bis + sezione G + `ILAC-G17:01/2021` in fonti e § C), `A1_validazione_metodi` (callout "validazione → incertezza" + ILAC-G17), `A2_incertezza_misura` (callout "ponte con la validazione" + ILAC-G17, con nota su EA-4/16 ritirata). File portati a v2.1.
- **EA/ILAC**: inserito **`ILAC-G17:01/2021`** (Measurement Uncertainty in Testing) come riferimento vigente; aggiunto al registro `00_FONTI_NORMATIVE` § 7. Documentato che **`EA-4/16` è stata ritirata (marzo 2021)** e non va più citata (nota anti-allucinazione + rinvio a `EA-INF/01`).

### Recepimento RT-23 rev.05 (15-04-2026, in vigore 01-11-2026)
- Aggiornato il registro fonti: RT-23 **rev.05** (rev. generale; rev.04 vigente fino al 31-10-2026), con **promemoria di vigenza attivo**. Aggiunti al registro **`ISO-17011:2018`** (§ 7.8.3, struttura del campo di accreditamento) e **`ILAC-G18`** (formulazione dello scopo).
- `00_GLOSSARIO`: aggiornate le voci "Campo di accreditamento" e "Flessibile"; aggiunte **Atlante**, **DAonline**, **ILAC-G18**, **Macromatrice**, **Matrici assimilabili**; precisata "Metodo di prova" (cambio tecnica/rivelatore/terreno → metodo sviluppato).
- `sez_7_2`: nuova sezione **B-ter** "Modifiche del metodo e campo di accreditamento (RT-23 rev.05)" — matrici assimilabili/non assimilabili, parametri da calcolo, campo di misura, fasi/conferme, metodi superati ("Ritirato"), metodi sviluppati in Atlante, campo flessibile.
- Allineati a rev.05 i riferimenti in `sez_5`, `sez_7_3`, `A7`, `A8` (con nota di transizione).

### AI-governance
- Nuovo file **`01_SYSTEM/06_GOVERNANCE_AI.md`** (sempre in contesto): strato generico e solido di governance dell'AI in laboratorio (deployment locale, riservatezza § 4.2/§ 7.11, competenza all'uso § 6.2, integrità/ALCOA+, audit AI-driven § 8.8/§ 8.9, confine base↔verticale, mappa ai § ISO). Collegato dalla persona `aicertus`.

### Qualità redazionale
- Uniformati i **marcatori di livello** alla forma canonica `(🟢 LIVELLO 1 — fatto normativo)` / `(🟢 LIVELLO 2 — prescrizione Accredia)` / `(🟡 LIVELLO 3 — prassi ispettiva)` / `(🟣 LIVELLO 4 — best practice)`.
- Normalizzate le **vocali accentate** (piu'/puo'/perche'/cio' → più/può/perché/ciò) nei file di logica ispettiva.

### Conteggio file
- Riconciliato il conteggio: **52 file .md** (51 di KB + CHANGELOG). 9 file `sempre_in_contesto` (incluso il nuovo `06_GOVERNANCE_AI`).

### In sospeso (prossimi passi)
- **Promemoria vigenza RT-23**: dopo il 01-11-2026 rimuovere le note "fino al 31-10-2026 vige rev.04".
- Estensione (opzionale) delle sezioni G AI-governance ad altri file requisito dove la differenza operativa è rilevante (oggi presenti in `sez_7_2` e `sez_7_6`).

---

## v2.0 — 2026-05-28 (refactor integrale)

**Refactor completo della KB base** post-audit Fase 1. Backup integrale della v1.0 disponibile come `backup_KB_pre_refactor_2026-05-28.zip` nella cartella outputs di Cowork.

### Cosa è cambiato

**Architettura**
- Cartelle gerarchiche (`00_*`, `01_SYSTEM/`, `02_LOGICA_ISPETTIVA/`, `03_REQUISITI/`, `04_APPENDICI/`, `05_ORCHESTRATORE/`).
- 50 file in totale (rispetto agli 11 della v1.0), ognuno autocontenuto e ottimizzato per il caricamento on-demand da parte di un LLM locale (Qwen).
- Frontmatter YAML uniforme su tutti i file (`codice`, `tipo`, `livello`, `sezione_norma`, `fonti_primarie`, `fonti_secondarie`, `file_collegati`, `tags`).
- Sintassi wikilink `[[codice_file]]` standardizzata su tutta la KB.

**Vincoli operativi formalizzati**
- 5 protocolli espliciti in `01_SYSTEM/`: persona dell'agente, ASK-BEFORE-ANSWER, HANDOFF-QM, CITAZIONI, TRE LIVELLI DI CONFIDENZA. Sono il sistema operativo dell'agente, sempre in contesto.

**Tassonomia a 4 livelli di contenuto**
- 1 — fatto normativo (ISO 17025:2018)
- 2 — prescrizione Accredia (RT-08 rev.05, RG, DT, politiche, raccomandazioni CdIG)
- 3 — prassi ispettiva (Denis Brazzo, ispettore Accredia DL0871)
- 4 — best practice tecnica (Eurachem, CITAC, Nordtest, EUROLAB, ISTISAN, JCGM, JRC, EU regs)
- Ogni file marca esplicitamente il livello dei propri contenuti, l'agente non può presentare un livello 3/4 come fosse 1/2.

**Registro fonti chiuso**
- `00_FONTI_NORMATIVE.md` elenca ~50 documenti applicabili ai laboratori di prova, estratti integralmente dalla sezione 3 di **LS-04 rev.20 (in vigore 01-10-2025)**. Per ciascun documento: codice, revisione, data, accesso (free / a pagamento), scope di utilizzo nella KB. L'agente cita solo da qui; se un documento serve e non è nel registro, l'agente chiede al QM.

### Errori normativi della v1.0 corretti

1. **Sezione 8 ricostruita integralmente sull'indice ufficiale ISO/IEC 17025:2018**, confermato da RT-08 rev.05 (PDF Accredia fornito da Denis):
   - 8.1 Opzioni A/B
   - 8.2 Documentazione del SGQ
   - 8.3 Controllo dei documenti del SGQ
   - 8.4 Controllo delle registrazioni
   - 8.5 Azioni per affrontare rischi e opportunità *(prima era "Misurazioni, analisi e miglioramento" — titolo di ISO 9001, non esistente in 17025:2018)*
   - 8.6 Miglioramento *(prima era "Azioni correttive")*
   - 8.7 Azioni correttive *(prima era "Miglioramento")*
   - 8.8 Audit interni *(prima era "Reclami" — errore: in 17025 i Reclami sono § 7.9)*
   - 8.9 Riesame di direzione *(prima era "Gestione NC da accreditatore" — non è un § della norma)*

2. **Rimpatrio dei Reclami in § 7.9** (la posizione corretta secondo ISO 17025:2018), con nuovo file `sez_7_9_reclami.md` e disclaimer esplicito.

3. **Aggiunta di file dedicati per requisiti prima mancanti o incidentali**: § 8.1 Opzioni A/B (chiave per il SGQ), § 8.5 Rischi e opportunità (novità ISO 17025:2018 vs. 2005), § 8.6 Miglioramento, § 8.8 Audit interni (era trattato sotto la falsa 8.5), § 8.9 Riesame di direzione (idem).

4. **Granularità per sotto-sezione**: la sezione 6 (Risorse) è ora 6 file separati (6.1–6.6), la sezione 7 (Processo) è 11 file separati (7.1–7.11), la sezione 8 (SGQ) è 9 file separati (8.1–8.9). Ogni file è ~150-400 righe e si carica autonomamente.

### Contenuti nuovi (non presenti nella v1.0)

- **8 appendici trasversali** (`04_APPENDICI/`) con best practice ispettive sui temi più critici in audit: validazione metodi, incertezza di misura, riferibilità/taratura, CRM, PT/ILC, audit interni & riesame, campionamento, gestione NC da Accredia.

- **Glossario operativo** (`00_GLOSSARIO.md`) con definizioni ISO/Accredia ricorrenti.

- **00_MAPPA_ESPLOSA** in `03_REQUISITI/` come indice navigabile.

- **Pattern di escalation al QM** integrato in ogni Domanda 3 ("efficacia"): l'agente non si ferma alla domanda ispettiva, propone al QM una bozza di azione con riferimento normativo, urgenza suggerita e indicazione di cosa NON può chiudere da solo.

### Cosa è rimasto (porting con uplift)

I 3 file di logica ispettiva originali (Denis Brazzo, 2026-05-06) sono stati portati nella nuova cartella `02_LOGICA_ISPETTIVA/` mantenendo il contenuto integralmente:
- `01_COME_PENSA_ISPETTORE.md` (ex `come-pensa-un-ispettore.md`)
- `03_RISPOSTA_NC.md` (ex `come-rispondere-NC-ACCREDIA.md`)
- `04_ERRORI_FATALI.md` (ex `errori-fatali-da-evitare.md`)

A questi si è aggiunto un nuovo `02_TRE_LIVELLI_DI_CONTROLLO.md` che formalizza il pattern "meccanismo → estensione → efficacia" applicato in tutti i file di requisito.

Il file `LabNexus_CoWork_KB.md` (architettura agentica per modalità CoWork) è stato portato in `05_ORCHESTRATORE/` con wikilink aggiornati ai nuovi codici.

### Anti-allucinazione: verifiche eseguite

- Spot-check qualità su 6 file campione (`01_PERSONA_aicertus`, `sez_6_4_dotazioni`, `sez_7_9_reclami`, `sez_8_5_rischi_opportunita`, `sez_8_8_audit_interni`, `A2_incertezza_misura`): zero § di norma inventati, zero codici Accredia falsi, zero testo integrale ISO copiato, zero frasi cosmetiche vietate.
- Verifica programmatica dei wikilink interni: zero refusi residui (i ~26 link rimasti puntano a documenti della **KB verticale** del singolo laboratorio, che esistono solo dopo l'onboarding del lab — vedi nota in `00_INDICE.md` § 5).
- Verifica frontmatter YAML: presente in tutti i 50 file.
- Verifica numerazione sezioni: nessuna citazione di § ISO 17025 inventato; tutti i riferimenti a "§ 8.8 Reclami" o "§ 8.5 Misurazioni…" sono note storiche di disclaimer, non claim attivi.

### Sintesi dimensionale

- 50 file `.md`
- ~11.500 righe totali
- ~730 KB
- Costruita su 3 fonti autoritative ufficiali fornite da Denis Brazzo (PDF):
  - UNI CEI EN ISO/IEC 17025:2018
  - RT-08 rev.05 (Accredia, EC del 10-02-2022)
  - LS-04 rev.20 (Accredia, in vigore 01-10-2025)

---

## v1.0 — 2026-05-06

Versione iniziale costruita da Denis Brazzo (Ispettore Accredia DL0871). Conteneva 11 file: README, mappa esplosa, 4 file di logica ispettiva, 5 file di domande per sezione (4-8), aicertus, LabNexus_CoWork_KB. Sezione 8 con numerazione errata mutuata da ISO 9001. Reclami collocati erroneamente in § 8.8. Conservata integralmente nel backup `backup_KB_pre_refactor_2026-05-28.zip`.
