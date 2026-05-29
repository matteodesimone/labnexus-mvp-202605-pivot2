---
codice: sez_7_9_reclami
tipo: requisito_normativo
livello: misto
titolo: "§ 7.9 — Reclami"
sezione_norma: "7.9"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.9" }
  - { doc: "RT-08", rev: "05", paragrafo: "7.9" }
fonti_secondarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "8.8 (Audit interni, NON Reclami)" }
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[02_TRE_LIVELLI_DI_CONTROLLO]]"
  - "[[sez_7_10_attivita_non_conformi]]"
  - "[[sez_7_8_presentazione_risultati]]"
tags: [reclami, gestione_reclami, processo_documentato, indipendenza, comunicazione_cliente]
escalation_qm: true
---

# § 7.9 — Reclami

> **NOTA CRITICA — RICOLLOCAZIONE**: nella vecchia KB i Reclami erano erroneamente collocati in § 8.8. In **ISO/IEC 17025:2018 § 8.8 è "Audit interni"**, NON i Reclami. I Reclami sono **§ 7.9** della norma. Questo file è il riferimento corretto. Chiunque cerchi "reclami in 8.8" sta usando una vecchia KB o sta confondendo con altre norme (es. ISO 9001:2015 § 9.1.2 / 10.2; ISO 17020 § 7.5). Per gli **audit interni**, si veda il file dedicato `sez_8_8_audit_interni`.

> Il § 7.9 chiede al laboratorio un **processo documentato** per ricevere, valutare e decidere sui reclami. Punti chiave: indipendenza del valutatore dal personale coinvolto, conferma di ricezione al reclamante, comunicazione formale dell'esito, validità del processo (nessun decisore può essere stato direttamente coinvolto nell'attività oggetto del reclamo).

## A. Testo della norma (🟢 LIVELLO 1 — fatto normativo)

> **Fonte**: `[ISO-17025-2018 § 7.9]` — testo a pagamento, non riprodotto. Riassunto operativo:

- **7.9.1** Il laboratorio deve avere un **processo documentato** per ricevere, valutare e prendere decisioni su reclami.
- **7.9.2** Su richiesta di una qualsiasi parte interessata, deve essere disponibile una **descrizione del processo** di gestione dei reclami.
- **7.9.3** Alla ricezione del reclamo, il laboratorio deve **confermare** se il reclamo è relativo ad attività di sua responsabilità e, in tal caso, deve **gestirlo**. Il laboratorio è responsabile di tutte le decisioni a tutti i livelli del processo di gestione del reclamo.
- **7.9.4** Il processo deve includere almeno: **a)** descrizione di processo per ricezione, validazione, indagine e decisioni; **b)** tracciamento e registrazioni dei reclami, incluse azioni intraprese per risolverli; **c)** assicurazione che ogni azione appropriata venga intrapresa.
- **7.9.5** Il laboratorio che riceve il reclamo deve essere **responsabile di raccogliere e verificare** tutte le informazioni necessarie per validare il reclamo.
- **7.9.6** Ove possibile, il laboratorio deve **dare conferma di ricezione** del reclamo, fornire al reclamante report sull'avanzamento e l'esito.
- **7.9.7** Le decisioni comunicate al reclamante devono essere **prese o riesaminate e approvate da personale non coinvolto** nelle attività del laboratorio oggetto del reclamo. **Nota**: questo può essere effettuato da personale esterno. Se ciò non è possibile, possono essere implementate misure alternative per assicurare l'imparzialità.
- **7.9.8** (in alcune edizioni: 7.9.7 termina con la chiusura formale) Ovunque possibile, il laboratorio deve fornire la **notifica formale dell'esito** del processo di gestione del reclamo al reclamante.

## B. Prescrizioni Accredia (🟢 LIVELLO 2 — prescrizione Accredia)

> **Fonte**: `[RT-08 rev.05 § 7.9]`

I punti **7.9.1 → 7.9.7** sono dichiarati come "Si applica il requisito di norma" — RT-08 non aggiunge integrazioni specifiche al § 7.9. Tutto il peso ispettivo ricade sull'attuazione corretta del requisito di norma.

Nota operativa: nonostante l'assenza di prescrizioni Accredia aggiuntive, **i reclami sono uno dei principali input** che ACCREDIA esamina durante una visita per valutare:
- la **maturità del sistema di gestione**;
- la **trasparenza** del laboratorio verso il cliente;
- la **capacità di analisi causale** (collega § 7.10 — attività non conformi — e § 8.7 — azioni correttive).

Inoltre, i reclami pervenuti ad ACCREDIA direttamente (non solo al laboratorio) sono noti al team ispettivo: il laboratorio deve aver registrato gli stessi reclami e aver dimostrato gestione coerente.

## C. Documenti applicabili (LS-04 rev.20)

- `ISO-17025-2018` § 7.9.
- `RT-08 rev.05` § 7.9.
- (Nessun documento Accredia specifico aggiuntivo sui reclami al di fuori di RT-08.)
- Riferimento esterno utile (non vincolante per 17025): `ISO 10002` — Quality management — Customer satisfaction — Guidelines for complaints handling.

## D. Domande ispettive (🟡 LIVELLO 3 — prassi ispettiva)

> Pattern: meccanismo (processo documentato) → estensione (campione reclami reali) → efficacia (analisi cause, collegamento con NC e AC, indipendenza decisione).

### 7.9.1-7.9.4 — Processo documentato

#### Domanda 1 — Meccanismo
**"Mostratemi la procedura di gestione reclami. Voglio vedere: canali di ricezione (mail, telefono, portale, ACCREDIA che inoltra), criteri di **validazione** (è davvero un reclamo? è di mia responsabilità?), fasi (ricezione → registrazione → analisi → decisione → comunicazione → chiusura), responsabilità chiare in ciascuna fase, tempi target di risposta, supporto di registrazione (modulo, ticket nel LIMS, registro)."**

Analisi ispettiva: la procedura deve essere completa. Tipici buchi:
- canali non strutturati (es. il reclamo via telefono non viene mai registrato);
- assenza di criterio per distinguere "reclamo" da "richiesta di chiarimento";
- nessun tempo target di risposta.

#### Domanda 2 — Estensione
**"Disponibilità della descrizione del processo per parti interessate (§ 7.9.2): è pubblicata sul sito web? È fornita su richiesta? Come la mantenete aggiornata?"**

Analisi ispettiva: §7.9.2 esplicito. Se la procedura è "riservata interna" e mai disponibile a un cliente che la chiede → NC.

#### Domanda 3 — Efficacia
**"Mostratemi il **registro reclami** degli ultimi 24 mesi: quanti sono, tempi medi di chiusura, classificazione tipologica (es. errori RdP, ritardi, comportamento personale, fatturazione, contestazione tecnica del risultato)."**

Analisi ispettiva: il registro deve esistere. Se "abbiamo solo 1 reclamo in 2 anni" su grandi volumi → poco credibile, probabile sottoregistrazione. Esamino la tipologia: se sono tutti "fatturazione", potrebbe esserci sottoregistrazione delle contestazioni tecniche.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Procedura reclami completa + canali strutturati + criteri validazione + registro reclami + disponibilità descrizione processo a parti interessate. ISO § 7.9.1-7.9.4."*
- Urgenza: media — è igiene di sistema, ma quando c'è un caso ispettivo concreto sale ad alta.
- Cosa NON può chiudere l'agente: NON dichiara conforme la procedura senza vedere il registro e almeno 3 reclami reali gestiti.

### 7.9.5-7.9.6 — Raccolta info e comunicazione con reclamante

#### Domanda 1 — Meccanismo
**"Per ogni reclamo, come raccogliete e verificate le informazioni? Coinvolgete il reclamante per chiarimenti? Quale traccia tenete della corrispondenza?"**

Analisi ispettiva: la raccolta info deve essere documentata. Tipico problema: il reclamo arriva, viene "interpretato" dal capo lab senza interagire col reclamante, decisione presa unilateralmente.

#### Domanda 2 — Estensione
**"Mostratemi 3 reclami reali. Voglio vedere: data ricezione, conferma di ricezione al reclamante (§ 7.9.6), update intermedi, comunicazione esito formale."**

Analisi ispettiva: i 3 elementi (conferma ricezione, update, comunicazione esito) sono prescritti dalla norma "ove possibile". Tipico problema: conferma ricezione assente o solo verbale.

#### Domanda 3 — Efficacia
**"Per un reclamo tecnico contestato sul risultato di una prova: come avete validato la contestazione? Avete ricontattato il campione (se ancora disponibile)? Rifatto la prova? Verificato registrazioni tecniche (§ 7.5)?"**

Analisi ispettiva: efficacia tecnica. Collegamento esplicito con § 7.5 (ricostruzione) + § 7.10 (NC se confermato).

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Raccolta info reclamo + conferma ricezione + comunicazione esito + validazione tecnica per reclami su risultati. ISO § 7.9.5-7.9.6."*

### 7.9.7 — Indipendenza del decisore

#### Domanda 1 — Meccanismo
**"Chi prende la decisione finale sul reclamo? Come assicurate che il decisore (o chi riesamina/approva) non sia stato coinvolto nelle attività oggetto del reclamo (§ 7.9.7)?"**

Analisi ispettiva: questo è il punto più delicato. La norma è chiara: il decisore (o approvatore) non deve essere stato coinvolto. Tipico problema:
- Reclamo su un RdP firmato dal responsabile tecnico, decisione presa dallo stesso responsabile tecnico → NC esplicita.
- "Misure alternative per assicurare imparzialità" senza spiegare quali → NC.

#### Domanda 2 — Estensione
**"In un lab piccolo dove tutto il personale tecnico ha potenzialmente toccato la prova: quali misure alternative implementate per assicurare imparzialità (norma esplicita)?"**

Analisi ispettiva: opzioni accettabili: revisione da personale esterno (es. consulente, altro lab del gruppo), comitato interno con membri non coinvolti, scalata a direzione amministrativa con supporto tecnico documentato. Va dichiarato in procedura.

#### Domanda 3 — Efficacia
**"Per i 3 reclami visti prima, identificate il decisore e dimostratemi che non era stato coinvolto nelle attività contestate."**

Analisi ispettiva: prova del nove. Se il decisore è la stessa persona che ha firmato il RdP contestato → NC immediata.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Indipendenza decisione reclami: § 7.9.7 della norma. In lab piccolo, misure alternative documentate. Esempi reali: tracciabilità decisore vs personale coinvolto nell'attività contestata."*
- Urgenza: alta — è un requisito sostanziale di imparzialità.

### Collegamenti con § 7.10 e § 8.7

#### Domanda 1 — Meccanismo
**"Quando un reclamo viene validato (cioè confermato come fondato), come si collega con il sistema NC (§ 7.10) e con le azioni correttive (§ 8.7)?"**

Analisi ispettiva: ogni reclamo fondato dovrebbe generare almeno una NC. Se non c'è collegamento, il reclamo viene "chiuso" senza analisi cause e senza miglioramento. NC implicita su SGQ.

#### Domanda 2 — Estensione
**"Per i reclami fondati degli ultimi 12 mesi: mostrate le NC corrispondenti e le AC intraprese."**

Analisi ispettiva: matching reclami ↔ NC ↔ AC. Tipico problema: i reclami sono in un sistema, le NC in un altro, e non si parlano.

#### Domanda 3 — Efficacia
**"Per un'AC chiusa da reclamo: come avete verificato l'efficacia (§ 8.7.3)? Avete avuto reclami simili dopo?"**

Analisi ispettiva: prova di efficacia.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Collegamento Reclami ↔ NC ↔ AC: per reclami fondati, evidenza di NC aperte, AC chiuse, verifica efficacia. ISO § 7.9 + § 7.10 + § 8.7."*

### Reclami pervenuti ad ACCREDIA

#### Domanda 1 — Meccanismo
**"Avete avuto reclami pervenuti direttamente ad ACCREDIA che vi siano stati inoltrati? Come li trattate?"**

Analisi ispettiva: gestione standard, ma con doppia tracciabilità (registro lab + risposta ad ACCREDIA).

#### Domanda 2 — Estensione
**"Mostrate l'evidenza di gestione di un eventuale reclamo via ACCREDIA."**

#### Domanda 3 — Efficacia
**"Confermate che il vostro registro reclami è allineato con quanto ACCREDIA potrebbe avervi inoltrato?"**

Analisi ispettiva: importante perché ACCREDIA porta in visita la lista dei reclami pervenuti. Disallineamenti → NC.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Allineamento registro reclami interno ↔ reclami via ACCREDIA. Risposte tempestive e documentate."*

## E. Errori fatali correlati

- **EF-7.9-A**: assenza di procedura documentata di gestione reclami — NC esplicita § 7.9.1.
- **EF-7.9-B**: descrizione del processo non disponibile su richiesta di parti interessate — NC § 7.9.2.
- **EF-7.9-C**: reclami non tracciati su un registro / registrazione assente o incompleta — NC § 7.9.4.
- **EF-7.9-D**: decisore del reclamo coinvolto nelle attività contestate, senza misure alternative documentate — NC esplicita § 7.9.7.
- **EF-7.9-E**: reclamo fondato chiuso senza apertura di NC e senza azione correttiva — NC trasversale § 7.9 + § 7.10 + § 8.7.
- **EF-7.9-F**: reclamo pervenuto via ACCREDIA non registrato nel registro interno → disallineamento grave.
- **EF-7.9-G**: assenza di comunicazione di esito al reclamante.
- **EF-7.9-H**: sottoregistrazione cronica (es. solo 1-2 reclami in 2 anni su grandi volumi) — indizio di NC sostanziale di sistema.

Vedi anche [[04_ERRORI_FATALI]].

## F. Quando l'agente DEVE chiedere prima di rispondere

L'agente **non risponde** e **chiede al QM** se:

1. Il reclamo coinvolge contestazione tecnica di un risultato già emesso: profilo § 7.5 (ricostruzione) + § 7.10 (NC) + § 7.8.8 (correzione RdP).
2. Il decisore "naturale" sarebbe coinvolto e va identificata una misura alternativa di imparzialità.
3. Il reclamo è pervenuto via ACCREDIA: gestione e risposta vanno coordinate con direzione + Responsabile Qualità.
4. Il reclamo solleva sospetti di alterazione dati / cattiva fede del personale: profilo ben oltre il § 7.9 — può richiedere azioni di gestione del personale e comunicazione ACCREDIA.
5. Il reclamo riguarda potenziale violazione di riservatezza (§ 4.2): profilo specifico.
6. La risposta al reclamante richiede release di registrazioni tecniche non normalmente comunicate al cliente: profilo riservatezza vs trasparenza.
7. Il reclamo riguarda fatturazione/contrattuale ma incrocia aspetti tecnici (es. cliente lamenta che è stata fatturata attività non eseguita): può nascondere problema tecnico sottostante.

L'agente NON gestisce mai un reclamo in autonomia. Registra, attiva il QM/RT secondo procedura, e supporta la fase istruttoria — la decisione finale è del decisore indipendente come previsto da § 7.9.7.
