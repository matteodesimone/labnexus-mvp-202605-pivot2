---
codice: sez_7_10_attivita_non_conformi
tipo: requisito_normativo
livello: misto
titolo: "§ 7.10 — Attività non conformi"
sezione_norma: "7.10"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.10" }
  - { doc: "RT-08", rev: "05", paragrafo: "7.10" }
fonti_secondarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "8.7 (Azioni correttive)" }
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.8.8 (Correzione RdP)" }
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[02_TRE_LIVELLI_DI_CONTROLLO]]"
  - "[[03_RISPOSTA_NC]]"
  - "[[sez_7_5_registrazioni_tecniche]]"
  - "[[sez_7_7_assicurazione_validita]]"
  - "[[sez_7_8_presentazione_risultati]]"
  - "[[sez_7_9_reclami]]"
tags: [NC, attivita_non_conformi, autosospensione, ricaduta_RdP, gestione_rischio, ripristino_attivita]
escalation_qm: true
---

# § 7.10 — Attività non conformi

> Il § 7.10 disciplina cosa fare quando un'attività **non rispetta** le proprie procedure o i requisiti del cliente: gestione, valutazione di significatività, decisione su accettabilità del lavoro, eventuale fermo dell'attività, comunicazione al cliente, decisione di richiamo dei lavori. È il "interruttore di sicurezza" del processo. Accredia preme su un punto specifico: quando una NC può mettere in dubbio i risultati emessi (es. PT negativo, deriva strumentale rilevata), il laboratorio deve **autosospendere** verso ACCREDIA e **revisionare i RdP emessi**. Vedi anche la guida operativa [[03_RISPOSTA_NC]].

## A. Testo della norma (🟢 LIVELLO 1 — fatto normativo)

> **Fonte**: `[ISO-17025-2018 § 7.10.1 - 7.10.3]` — testo a pagamento, non riprodotto. Riassunto operativo:

- **7.10.1** Il laboratorio deve avere una **procedura** che si applica quando un aspetto dell'attività di laboratorio o i risultati di tale attività non rispettano le proprie procedure o i requisiti concordati col cliente. La procedura deve assicurare:
  - **a)** responsabilità e autorità per la gestione del lavoro non conforme definite;
  - **b)** azioni (incluse interruzione del lavoro o ripetizione del lavoro, e trattenuta dei rapporti di prova, secondo necessità) basate sui livelli di rischio stabiliti dal laboratorio;
  - **c)** valutazione del **significato** del lavoro non conforme, incluso un'**analisi di impatto sui risultati precedenti**;
  - **d)** decisione sull'**accettabilità** del lavoro non conforme;
  - **e)** ove necessario, **notifica al cliente** e richiamo del lavoro;
  - **f)** definizione di responsabilità per l'**autorizzazione alla ripresa** del lavoro.
- **7.10.2** Il laboratorio deve **conservare le registrazioni** del lavoro non conforme e delle azioni intraprese, come specificato in 7.10.1, da b) a f).
- **7.10.3** Quando la valutazione indica che il lavoro non conforme potrebbe **ripresentarsi** o c'è dubbio sulla conformità delle operazioni del laboratorio con il proprio sistema di gestione, il laboratorio deve implementare **azioni correttive** (collega § 8.7).

## B. Prescrizioni Accredia (🟢 LIVELLO 2 — prescrizione Accredia)

> **Fonte**: `[RT-08 rev.05 § 7.10]`

### 7.10.3 — Autosospensione dell'accreditamento

> **Citazione integrale RT-08** (CRITICA):
> *"Qualora si verifichino non conformità che potrebbero comportare la sospensione dell'attività di prova/campionamento o mettere in dubbio la validità dei risultati di prove accreditate (ad es. valutazioni negative della partecipazione a circuiti interlaboratorio, indisponibilità o deterioramento di risorse, ecc.) il Laboratorio dovrà comunicare ad ACCREDIA l'**autosospensione dell'accreditamento** per tali prove, fino alla avvenuta verifica positiva, da parte di ACCREDIA, delle azioni correttive intraprese."*

Quindi NC di certo tipo (PT negativo, indisponibilità/deterioramento risorse critiche) → obbligo di **autosospensione** + comunicazione ACCREDIA + ripresa solo dopo verifica positiva ACCREDIA delle AC.

### 7.10.3 — Ricaduta su RdP già emessi

> **Citazione integrale RT-08**:
> *"Nel caso di non conformità riscontrate sui rapporti di prova già emessi, che possano pregiudicare l'utilizzo dei risultati da parte dei clienti (es. non conformità di natura tecnica, oppure relative al riferimento all'accreditamento), il Laboratorio deve individuare i rapporti di prova affetti dalla medesima non conformità ed avvertire i clienti emettendo rapporti di prova sostitutivi."*

Quindi NC su RdP emessi → **individuare TUTTI i RdP affetti dalla stessa carenza** + comunicare clienti + emettere RdP sostitutivi (collega § 7.8.8).

**Punti 7.10.1 e 7.10.2**: "Si applica il requisito di norma" senza integrazioni Accredia esplicite.

## C. Documenti applicabili (LS-04 rev.20)

- `ISO-17025-2018` § 7.10.
- `RT-08 rev.05` § 7.10.
- Collegamenti diretti: § 8.7 (Azioni correttive), § 7.8.8 (Correzioni RdP), § 7.9 (Reclami), § 7.7.3 (Analisi monitoraggio), § 7.5 (Registrazioni).

## D. Domande ispettive (🟡 LIVELLO 3 — prassi ispettiva)

> Pattern: meccanismo → estensione → efficacia.

### 7.10.1 — Procedura NC

#### Domanda 1 — Meccanismo
**"Mostratemi la procedura di gestione attività non conformi. Voglio vedere: (a) responsabilità e autorità (chi rileva, chi classifica, chi decide su trattenimento del lavoro, chi autorizza ripresa); (b) criteri di valutazione del rischio per decidere se fermare o no l'attività; (c) come valutate il significato della NC (impatto su risultati precedenti); (d) decisione formale di accettabilità; (e) criteri di notifica al cliente; (f) chi autorizza la ripresa."**

Analisi ispettiva: la procedura deve coprire i 6 elementi del § 7.10.1. Tipici buchi:
- "Responsabilità diffuse": tutti possono aprire una NC ma nessuno è responsabile della decisione finale.
- Manca criterio per decidere fra "ripetiamo il lavoro" e "fermiamo l'attività".
- Analisi di impatto sui risultati precedenti assente.
- Notifica al cliente solo verbale, mai tracciata.

#### Domanda 2 — Estensione
**"Mostratemi il registro NC degli ultimi 12 mesi. Voglio capire: numero totale, classificazione tipologica (tecniche, gestionali, ambientali, personale, ecc.), tempi medi di chiusura, percentuale che ha portato a notifica cliente, percentuale che ha portato ad AC (collega § 8.7)."**

Analisi ispettiva: il registro deve esistere e essere vivo. Se "abbiamo 2 NC in 12 mesi" su grandi volumi, c'è sottoregistrazione. Anche distribuzione tipologica: se sono tutte "ambientali" e nessuna "tecnica" → poco credibile.

#### Domanda 3 — Efficacia
**"Scelgo 3 NC dal registro. Per ciascuna: come è stata aperta, chi l'ha classificata, chi ha valutato il significato, chi ha deciso accettabilità, chi ha autorizzato la ripresa, c'è stata notifica cliente, c'è stata AC?"**

Analisi ispettiva: prova del nove. Il registro è formalmente popolato ma le 6 fasi del § 7.10.1 sono effettivamente tracciate?

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Procedura NC: copertura 6 elementi § 7.10.1, registro vivo, distribuzione tipologica, tracciabilità 6 fasi su 3 NC reali. ISO § 7.10.1-7.10.2."*
- Urgenza: alta.
- Cosa NON può chiudere l'agente: NON dichiara la gestione NC adeguata senza vedere registro + esempi reali completi.

### 7.10.1.c — Valutazione significatività e impatto sui risultati precedenti

#### Domanda 1 — Meccanismo
**"Quando aprite una NC, come decidete se può aver impattato risultati GIÀ EMESSI? Quali criteri usate? È un'analisi sempre obbligatoria o solo per NC "gravi"?"**

Analisi ispettiva: § 7.10.1.c è **esplicito** — l'analisi di impatto sui risultati precedenti è parte della valutazione di significato. Tipico problema: l'analisi viene fatta solo se "qualcuno se ne preoccupa". Manca obbligatorietà.

#### Domanda 2 — Estensione
**"Per una NC tecnica recente (es. apparecchio fuori taratura, materiale di riferimento scaduto, deriva strumentale): mostrate l'analisi di impatto sui RdP emessi nel periodo affetto. Avete individuato e ricontattato i clienti?"**

Analisi ispettiva: per NC che impattano accuratezza/affidabilità, RT-08 § 7.10.3 (seconda parte) chiede esplicitamente di individuare i RdP affetti dalla **stessa carenza** + avvertire i clienti + emettere RdP sostitutivi. Tipico buco: il lab risolve la NC interna ma non valuta i RdP emessi.

#### Domanda 3 — Efficacia
**"Per la stessa NC: mostrate i RdP sostitutivi emessi e la comunicazione al cliente con identificazione di annullamento/sostituzione del RdP originale (collega § 7.8.8)."**

Analisi ispettiva: tracciabilità end-to-end NC → identificazione RdP → comunicazione → RdP sostitutivo. Se uno qualsiasi degli anelli manca, NC sull'NC.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Analisi impatto NC sui risultati precedenti + individuazione RdP affetti + comunicazione clienti + emissione RdP sostitutivi. ISO § 7.10.1.c-e + RT-08 § 7.10.3 (seconda parte) + § 7.8.8."*
- Urgenza: massima.
- Cosa NON può chiudere l'agente: NON dichiara conforme la gestione NC se manca l'analisi di impatto sui RdP emessi.

### 7.10.3 — Autosospensione verso ACCREDIA

#### Domanda 1 — Meccanismo
**"In quali casi prevedete l'autosospensione dell'accreditamento? La vostra procedura riprende l'elenco RT-08 § 7.10.3 (PT negativo, indisponibilità/deterioramento risorse, dubbio validità risultati)?"**

Analisi ispettiva: l'autosospensione è una procedura **non discrezionale** quando si applicano le condizioni RT-08. Se la procedura non la prevede esplicitamente o la rende discrezionale ("a giudizio del responsabile") → NC.

#### Domanda 2 — Estensione
**"Avete autosospeso negli ultimi anni? Per quale motivo? Quali sono state le AC? Quale verifica positiva ACCREDIA ha consentito la ripresa?"**

Analisi ispettiva: cerco esempi reali. Se "non abbiamo mai autosospeso" su lab grande con molti PT, fa pensare. Verifico anche su outlier PT (collega § 7.7).

#### Domanda 3 — Efficacia
**"Se domani un PT torna con z>3 su un metodo accreditato critico: cosa fate, oltre alla NC interna? Comunicate ad ACCREDIA? Su quale modulo / canale? In quanti giorni?"**

Analisi ispettiva: test di scenario. La risposta deve essere immediata e precisa.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Autosospensione ACCREDIA: trigger esplicitati (PT negativo, indisponibilità risorse), procedura non discrezionale, esempi reali. RT-08 § 7.10.3 (prima parte)."*
- Urgenza: massima.
- Cosa NON può chiudere l'agente: NON valuta in autonomia se autosospendere; il QM/RT decide con direzione e comunica ad ACCREDIA.

### 7.10.3 + § 8.7 — Collegamento con azioni correttive

#### Domanda 1 — Meccanismo
**"Quali NC generano AC? Avete criteri (es. ricorrenza, gravità, rischio sistemico)? Come si collega il flusso NC → AC?"**

Analisi ispettiva: § 7.10.3 è esplicito — se NC potrebbe ripresentarsi o c'è dubbio sulla conformità delle operazioni col SGQ → AC. Non tutte le NC richiedono AC, ma la decisione deve essere ragionata e tracciata.

#### Domanda 2 — Estensione
**"Mostratemi 3 NC che hanno generato AC e 3 che NON hanno generato AC. Per le seconde, motivazione tracciata?"**

Analisi ispettiva: bilanciamento. Se ogni NC genera AC → sistema gonfiato e burocratico. Se nessuna NC genera AC → sistema reattivo ma non preventivo. Il sano equilibrio passa dall'analisi cause (collega § 8.7).

#### Domanda 3 — Efficacia
**"Per le AC chiuse: come avete verificato l'efficacia (§ 8.7.3)?"**

Analisi ispettiva: efficacia AC è cruciale. Se chiusura solo "formale" senza follow-up reale → NC su SGQ.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Flusso NC → analisi cause → decisione AC sì/no → AC con verifica efficacia. ISO § 7.10.3 + § 8.7."*

## E. Errori fatali correlati

- **EF-7.10.1-A**: procedura NC priva di uno dei 6 elementi (responsabilità, criteri rischio, valutazione significato, decisione accettabilità, notifica cliente, autorizzazione ripresa).
- **EF-7.10.1-B**: analisi di impatto sui risultati precedenti **non eseguita** o non tracciata.
- **EF-7.10.1-C**: registro NC con sottoregistrazione cronica (es. < 5 NC/anno su lab grande).
- **EF-7.10.2-A**: registrazioni delle NC incomplete (mancano fasi b-f).
- **EF-7.10.3-A**: trigger di autosospensione RT-08 non esplicitati in procedura, o resi discrezionali — NC esplicita RT-08.
- **EF-7.10.3-B**: PT negativo / indisponibilità risorse critiche → NC aperta ma **autosospensione non comunicata ad ACCREDIA** — NC grave.
- **EF-7.10.3-C**: NC tecnica risolta internamente ma **RdP emessi non revisionati** — NC esplicita RT-08 § 7.10.3 + § 7.8.8.
- **EF-7.10.3-D**: identificati RdP affetti ma **clienti non avvertiti** o RdP sostitutivi non emessi.
- **EF-7.10-E**: NC ricorrente senza apertura AC § 8.7 — NC sul sistema di miglioramento.
- **EF-7.10-F**: NC su risultati pubblicati / inviati ad autorità (es. controllo ufficiale) gestita senza comunicazione formale all'autorità — possibili profili legali.

Vedi anche [[04_ERRORI_FATALI]] e [[03_RISPOSTA_NC]].

## F. Quando l'agente DEVE chiedere prima di rispondere

L'agente **non risponde** e **chiede al QM** se:

1. È stata identificata una potenziale NC e va decisa significatività + accettabilità + necessità di fermo: decisione esclusiva del QM/RT.
2. Una NC tecnica (es. PT negativo, taratura fuori, deriva strumentale) potrebbe impattare RdP già emessi: profilo § 7.10.3 + § 7.8.8 — analisi di impatto + comunicazione clienti.
3. Esiste un trigger di autosospensione RT-08 § 7.10.3: comunicazione formale ACCREDIA, decisione di direzione.
4. La NC riguarda dati inviati ad autorità (controllo ufficiale, ambientale, sanitario): valutazione di obblighi legali esterni.
5. La NC ricorrente non genera AC: l'agente segnala il rischio di NC sistemica ma non decide.
6. La NC è "borderline" — non chiaro se rientra nella procedura o è solo una "anomalia": classificazione decisa dal QM.
7. La NC coinvolge la riservatezza/integrità del personale (collusione, alterazione dati, falsificazione): profilo grave, va trattato fuori dalla normale procedura.
8. La NC mette a rischio l'imparzialità (§ 4.1): profilo strutturale.

L'agente NON dichiara mai una NC chiusa senza vedere: (a) analisi cause documentata, (b) impatto sui RdP precedenti valutato, (c) eventuali RdP sostitutivi emessi, (d) AC chiusa con verifica efficacia (se attivata), (e) comunicazione ACCREDIA se applicabile.
