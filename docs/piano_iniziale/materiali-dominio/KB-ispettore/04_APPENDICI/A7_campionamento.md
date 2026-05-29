---
codice: A7_campionamento
tipo: appendice_tecnica
livello: misto
titolo: "Appendice A7 — Campionamento accreditato e incertezza da campionamento"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.3" }
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.4" }
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.6.1" }
  - { doc: "RT-08", rev: "05", paragrafo: "7.3" }
  - { doc: "RT-23", rev: "05", note: "definizione campo accreditamento (campionamento)" }
fonti_secondarie:
  - { doc: "EURACHEM-Sampling", rev: "2nd ed. 2019" }
  - { doc: "ISTISAN-22/39", rev: "2023", note: "trad. it. Eurachem Sampling" }
  - { doc: "Nordtest-TR-537", rev: "ed. 4, 2017", note: "incertezza ambientale, include campionamento" }
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[00_GLOSSARIO]]"
  - "[[03_PROTOCOLLO_HANDOFF_QM]]"
  - "[[02_TRE_LIVELLI_DI_CONTROLLO]]"
  - "[[sez_7_3_campionamento]]"
  - "[[sez_7_4_manipolazione_oggetti]]"
  - "[[sez_6_2_personale]]"
  - "[[A1_validazione_metodi]]"
  - "[[A2_incertezza_misura]]"
tags: [campionamento, sampling_uncertainty, piano_campionamento, scope_campionamento, RT-23]
escalation_qm: true
---

> **Marcatore di livello: 🟢🟡🟣** — Appendice trasversale. Contenuti normativi (livello 1+2), prassi ispettiva (livello 3), best practice tecnica (livello 4).

# Appendice A7 — Campionamento accreditato

> **Riassunto operativo.** Il campionamento è una **attività accreditata distinta** dalle prove: chi è accreditato per "prove" non è automaticamente accreditato per "prove + campionamento". La distinzione cambia il dossier richiesto, il personale autorizzato, l'incertezza, il rapporto emesso. Questa appendice serve al QM/RT per impostare/verificare campionamento accreditato, incertezza di campionamento e collegamento con il § 7.3 di norma. L'agente la carica quando il QM parla di "siamo accreditati anche per il campionamento?", "piano di campionamento", "incertezza di campionamento", "tracciabilità del campione", "campionatore esterno".

---

## 1. Principio normativo (🟢 LIVELLO 1+2)

### 1.1 ISO 17025:2018 § 7.3 — Campionamento

> Testo a pagamento, non riprodotto. Sintesi operativa:

- **§ 7.3.1** — quando il laboratorio esegue campionamento di sostanze, materiali o prodotti per la successiva prova o taratura, deve avere un **piano di campionamento** e un **metodo di campionamento**. Il metodo deve indirizzare i fattori che devono essere controllati per assicurare la validità dei successivi risultati. Il piano e il metodo devono essere disponibili nella località di campionamento. Il piano deve essere basato, ove ragionevole, su metodi statistici appropriati.
- **§ 7.3.2** — il metodo di campionamento deve descrivere:
  - (a) la selezione dei campioni o dei siti,
  - (b) il piano di campionamento,
  - (c) la preparazione e il trattamento del/i campione/i di una sostanza, materiale o prodotto per produrre l'item richiesto per la prova o taratura successiva.
- **§ 7.3.3** — il laboratorio deve conservare le registrazioni del campionamento, che includono come applicabile:
  - riferimento al metodo di campionamento usato,
  - data e ora del campionamento,
  - dati per identificare e descrivere il campione (es. numero, quantità, nome),
  - identificazione del personale che esegue il campionamento,
  - identificazione dell'equipaggiamento usato,
  - condizioni ambientali o di trasporto,
  - diagrammi o altri mezzi equivalenti per identificare la posizione del campionamento, quando appropriato,
  - eventuali scostamenti, aggiunte o esclusioni dal metodo di campionamento e dal piano.

### 1.2 ISO 17025:2018 § 7.4 e § 7.6.1 — collegamenti

- **§ 7.4** — Manipolazione degli oggetti da provare o tarare: identificazione univoca, condizioni di conservazione, trasporto, ricezione. Coerente con la tracciabilità del campione dal punto di prelievo al banco di prova.
- **§ 7.6.1** — La valutazione dell'incertezza deve considerare **tutti i contributi significativi**, inclusi quelli derivanti dal **campionamento** quando rilevante.

### 1.3 RT-08 rev.05 § 7.3

> **Fonte**: `[RT-08 rev.05 § 7.3 — EC 10-02-2022]`

RT-08 ricorda che il campionamento è attività accreditabile come tale. Punti chiave:
- Quando il laboratorio è accreditato per il campionamento, il **campo di accreditamento (CAdA)** dichiarato secondo `RT-23` indica esplicitamente "**prove + campionamento**" per le voci interessate. Diversamente, l'accreditamento copre solo "**prove**" su campione fornito dal cliente.
- Il personale che esegue il campionamento accreditato deve essere **autorizzato esplicitamente** in matrice § 6.2.6, distintamente dall'autorizzazione a eseguire prove.
- Il **rapporto di prova** (§ 7.8.5) di prova eseguita su campionamento accreditato deve identificare le informazioni di campionamento (data, sito, metodo, personale) e dichiarare se l'attività di campionamento ricade nello scope di accreditamento.

### 1.4 RT-23 rev.05 — Definizione campo di accreditamento

> **Fonte**: `[RT-23 rev.05 (Accredia, free) — in vigore 01-11-2026; fino al 31-10-2026 vige rev.04]`

`RT-23` definisce le regole per dichiarare il CAdA. Per il campionamento:
- la voce di scope esplicita la combinazione **matrice × parametro × metodo di campionamento** (in modo analogo alla combinazione matrice × analita × metodo di prova),
- la voce "campionamento + prove" è formalmente diversa da "sole prove": l'estensione dell'una all'altra richiede modifica del CAdA secondo `RG-02`.

---

## 2. Come un ispettore lo verifica (🟡 LIVELLO 3 — prassi ispettiva)

> Pattern: meccanismo → estensione → efficacia.

### 2.1 Domanda di scoping (preliminare alle tre domande)

**"Siete accreditati per il campionamento o solo per le prove? Mostratemi il CAdA dichiarato secondo RT-23: quali voci hanno 'campionamento + prove' esplicito?"**

> Questa domanda di scoping è critica: l'intera istruttoria sul § 7.3 cambia in base alla risposta. Se il lab dice "siamo accreditati anche per il campionamento" ma in CAdA risulta "sole prove", siamo già in NC potenziale grave (uscita formale che dichiara campionamento accreditato non avendolo).

### 2.2 Domanda 1 — Meccanismo

**"Per ogni voce 'campionamento + prove' del CAdA, mostratemi: il piano di campionamento, il metodo di campionamento (riferimento normativo o procedura interna), la procedura di compilazione del modulo di campionamento al sito, l'autorizzazione del personale § 6.2.6 specifica per il campionamento, l'equipaggiamento di campionamento e la sua taratura/verifica."**

Analisi ispettiva:
- Cerco un **dossier per voce di CAdA**, non un documento unico generico.
- Verifico l'**equipaggiamento di campionamento** (es. campionatori automatici, contenitori dedicati, preservanti, ghiaccio secco): se è strumentazione che influenza la rappresentatività, deve essere tarato/verificato e tracciato (§ 6.4).
- Controllo le **autorizzazioni del personale**: campionamento è una voce separata in matrice § 6.2.6. Frequente: il tecnico è autorizzato alla prova ma non al campionamento, eppure firma il modulo di campionamento.

### 2.3 Domanda 2 — Estensione

**"Per il campionamento accreditato, esiste copertura della stagionalità/variabilità della matrice (es. acqua reflua, aria, suolo)? Esiste un piano statistico per la rappresentatività (n. punti, n. aliquote, distribuzione spaziale/temporale)? Le condizioni di trasporto e conservazione fino al laboratorio sono documentate?"**

Analisi ispettiva:
- Per matrici a forte variabilità (acque, aria, suoli, alimenti freschi, biota), un piano di campionamento "spot" (1 prelievo, 1 punto, 1 momento) raramente è scientificamente difendibile. Cerco evidenza che il piano consideri **eterogeneità** e **rappresentatività**.
- Verifico la **catena fredda** o le condizioni di conservazione richieste dal metodo: registrazione di T° all'arrivo, registrazione tempi di trasporto.
- Verifico la **tracciabilità campione → punto di prelievo**: ogni campione deve essere ricostruibile fino al sito, all'operatore, all'orario.

### 2.4 Domanda 3 — Efficacia

**"Negli ultimi 24 mesi, è successo che un campione sia stato rifiutato all'accettazione per problemi di campionamento (temperatura, tempo, contenitore, identificazione)? Cosa è stato fatto? Avete valutato l'incertezza da campionamento per i metodi dove è significativa, e come l'avete inserita nell'incertezza totale dichiarata?"**

Analisi ispettiva:
- Un sistema di campionamento maturo intercetta gli errori al check-in. Se "non è mai successo niente" in volumi significativi, o il controllo è cieco o gli errori passano.
- Verifico se l'**incertezza da campionamento** è stata stimata e dichiarata: l'incertezza totale dei risultati può essere dominata dal campionamento più che dall'analitica. Vedi `EURACHEM-Sampling 2019` / `ISTISAN-22/39`.

### 2.5 Pattern di errori comuni

1. **"Campionamento accreditato" dichiarato dal lab ma non in CAdA**: errore fatale di scope. NC su § 7.3 + RT-23 + § 7.8 (rapporto emesso fuori scope).
2. **Campionatore esterno** (terzista) che esegue campionamento "accreditato" senza essere autorizzato in matrice § 6.2.6 del lab e senza essere sub-fornitore qualificato ex § 6.6: doppia NC.
3. **Modulo di campionamento al sito** assente o incompleto (manca firma operatore, manca timestamp, manca identificativo punto): NC su § 7.3.3.
4. **Catena di custodia / trasporto** non tracciata: NC tecnica + impatto sull'integrità del campione (§ 7.4).
5. **Equipaggiamento di campionamento** non gestito come strumento (no taratura/verifica/registro): NC su § 6.4.
6. **Incertezza da campionamento ignorata** quando rilevante: NC potenziale su § 7.6.1.
7. **Piano di campionamento "uguale per tutti i clienti"** senza adattamento alla matrice o all'obiettivo: NC su § 7.3.1.

---

## 3. Best practice tecnica (🟣 LIVELLO 4 — best practice)

### 3.1 Riferimento d'elezione: Eurachem Sampling 2019 / ISTISAN 22/39

`EURACHEM-Sampling 2nd ed. 2019` e la sua traduzione italiana `ISTISAN-22/39:2023` sono il riferimento per la stima dell'**incertezza dovuta al campionamento**. Concetti chiave:

- L'incertezza totale di un risultato di prova su campione **prelevato dal laboratorio** è:
  - `u_totale² = u_campionamento² + u_analitica²`
- L'incertezza da campionamento è composta da contributi:
  - eterogeneità della matrice in situ (between-target),
  - eterogeneità del campione raccolto (within-target),
  - effetti del metodo di campionamento (es. perdite, contaminazione),
  - effetti del trasporto e della conservazione.
- Approcci principali:
  - **Approccio empirico (duplicati)**: si prelevano campioni in duplicato in modo casuale, si analizzano in duplicato. Schema ANOVA fornisce stime di u_campionamento e u_analitica.
  - **Approccio modellistico**: si costruisce un modello di tutti i contributi (analogo al GUM applicato al campionamento).
  - **Approccio interlaboratorio sul campionamento** (CTS — Collaborative Trial in Sampling): più campionatori operano allo stesso target.

### 3.2 Nordtest TR-537 (per laboratori ambientali)

`Nordtest-TR-537 ed. 4 2017` integra Eurachem Sampling con un approccio pragmatico orientato al settore ambientale (acque, suoli, aria). Fornisce esempi numerici e schemi di calcolo basati su dati che il laboratorio già genera (PT, repliche di campionamento).

### 3.3 Piano di campionamento — buone pratiche

- **Definizione dell'obiettivo del campionamento** (decisionale, conoscitivo, di conformità): determina la strategia.
- **Rappresentatività**: numero di sotto-campioni, distribuzione spaziale/temporale, criteri di selezione siti.
- **Tipologia campionamento**: random, stratificato, sistematico, judgmental — ciascuno con applicabilità diversa.
- **Equipaggiamento**: contenitori, preservanti, condizioni di trasporto, materiali a basso rilascio (vedere norme specifiche del settore: acque ISO 5667-x, aria ISO 16000-x, suoli ISO 18400-x, alimenti regolamenti EU, ecc.).
- **Tempistica**: tempo massimo dal prelievo all'analisi per analiti instabili.
- **Catena di custodia**: tracciamento documentato dal sito al banco di prova.

### 3.4 Settori cogenti — riferimenti

| Settore | Riferimenti tipici (oltre ISO 17025) |
|---|---|
| Acque (potabili, reflue, superficiali) | Famiglia ISO 5667; normative nazionali/EU; APAT/IRSA-CNR per Italia |
| Aria (indoor, outdoor, emissioni) | Famiglia ISO 16000; UNI EN per emissioni in atmosfera |
| Suoli | Famiglia ISO 18400 |
| Alimenti / mangimi | Reg. UE 2017/625 (controlli ufficiali); Reg. UE 2021/808 (residui farmacologici); regolamenti settoriali |
| Costruzioni | Reg. UE 305/2011, Reg. UE 2024/3110, D.Lgs. 106/2017 |

> L'agente cita questi riferimenti **come applicabili**, sempre con il caveat che è il QM a verificare se sono pertinenti allo scope del laboratorio.

### 3.5 Quando il campionamento NON è del laboratorio

Se il laboratorio non è accreditato per il campionamento e riceve campioni dal cliente, deve:
- al check-in, **verificare l'idoneità** del campione (contenitore, T°, integrità, identificazione, quantità): § 7.4,
- documentare le **anomalie** e decidere se accettare, accettare con riserva, rifiutare,
- nel rapporto di prova, **dichiarare esplicitamente** che il campionamento non è coperto dall'accreditamento del laboratorio e che il risultato si riferisce al campione "come ricevuto" (cfr. § 7.8.2),
- non garantire la rappresentatività del campione rispetto al lotto/sito da cui proviene.

---

## 4. Template operativo per il QM

### 4.1 Decision tree "Campionamento accreditato o no?"

```
La voce del CAdA per quella matrice/parametro include "campionamento + prove"?
    Sì → Campionamento accreditato. Si applicano tutti i § 7.3 e RT-23.
       → Personale autorizzato in matrice § 6.2.6 per campionamento
       → Equipaggiamento campionamento sotto § 6.4
       → Incertezza da campionamento valutata e dichiarata (§ 7.6.1)
       → Rapporto di prova menziona l'attività di campionamento come accreditata
    No → Solo prove. Campione fornito dal cliente.
       → Procedura di accettazione campione (§ 7.4)
       → Rapporto dichiara "campione come ricevuto"; risultato non estensibile al lotto/sito senza disclaimer
```

### 4.2 Modulo di campionamento al sito — campi minimi

| Campo | Note |
|---|---|
| Codice univoco campione | Generato secondo numerazione SGQ |
| Riferimento metodo di campionamento | Norma o procedura interna + revisione |
| Riferimento metodo di prova previsto | Per coerenza tra campionamento e analisi |
| Data e ora prelievo | Timestamp preciso |
| Sito / punto di prelievo | Coordinate, descrizione, eventuale schema/foto |
| Condizioni ambientali | T°, meteo, qualsiasi condizione rilevante |
| Operatore (cognome + matricola) | Firma in originale |
| Equipaggiamento usato | Codice strumenti, contenitori (tipo + lotto), preservanti |
| Quantità prelevata + n. aliquote | Per gli analiti che richiedono aliquote multiple |
| Condizioni di trasporto | Catena fredda, T° all'arrivo |
| Eventuali deviazioni dal piano | Con motivazione e firma |
| Catena di custodia | Sequenza di passaggi fino al laboratorio |

### 4.3 Checklist "Campionamento accreditato in audit interno"

- [ ] CAdA include esplicitamente la voce "campionamento + prove" per le attività ispezionate
- [ ] Personale campionatore autorizzato in matrice § 6.2.6 specifica per campionamento
- [ ] Piano e metodo di campionamento documentati e disponibili al sito
- [ ] Modulo di campionamento compilato in tutti i campi minimi
- [ ] Equipaggiamento di campionamento gestito come strumento (taratura/verifica, registro)
- [ ] Tracciabilità campione → punto/sito ricostruibile
- [ ] Condizioni di trasporto/conservazione documentate e rispettate
- [ ] Incertezza da campionamento valutata (con riferimento a Eurachem Sampling / Nordtest TR-537)
- [ ] Incertezza totale dei risultati dichiarata coerentemente
- [ ] Rapporto di prova identifica l'attività di campionamento come accreditata (con riferimento al metodo)
- [ ] Quando applicabile: comunicazione al cliente delle deviazioni dal piano

### 4.4 Note minime sul rapporto di prova (§ 7.8.5)

Per prove con campionamento accreditato, il rapporto include:
- riferimento al metodo di campionamento applicato (norma o procedura),
- data, ora, sito di campionamento,
- identificazione del personale campionatore,
- eventuali deviazioni dal metodo/piano,
- dichiarazione che il campionamento ricade nel CAdA accreditato (logo Accredia + numero accreditamento),
- contributo dell'incertezza da campionamento all'incertezza totale (se rilevante).

---

## 5. Quando l'agente attiva questa appendice

1. Il QM chiede *"siamo accreditati anche per il campionamento?"*, *"come gestiamo il campionatore esterno?"*, *"l'incertezza da campionamento è obbligatoria?"*.
2. Un rapporto di prova sta per essere emesso con menzione di campionamento accreditato: serve verifica della coerenza con CAdA.
3. È in lavorazione un'estensione di scope da "sole prove" a "campionamento + prove": serve istruttoria su personale, equipaggiamento, incertezza.
4. Una NC tocca § 7.3 (rilevata da Accredia o internamente).
5. Si analizza l'incertezza totale di un risultato dove il campionamento può essere dominante (matrici eterogenee, ambientale, agroalimentare).
6. Si valuta un sub-fornitore di campionamento (§ 6.6): serve qualificazione + autorizzazione.

---

## 6. Cosa l'agente NON può chiudere su questo tema (handoff QM obbligatorio)

Vedi `[[03_PROTOCOLLO_HANDOFF_QM]] § 4`. In materia di campionamento, l'agente **non** può:

1. **Dichiarare che il campionamento è "accreditato"**. La fonte è il CAdA dichiarato secondo RT-23; modifica del CAdA è atto formale verso Accredia (RG-02).
2. **Autorizzare un campionatore** in matrice § 6.2.6. Atto del RT/QM.
3. **Validare un piano di campionamento** per uno specifico cliente/sito. L'agente propone bozza, il RT valida tecnicamente.
4. **Decidere accettazione/rifiuto di un campione** all'arrivo. Atto del personale autorizzato (§ 7.4).
5. **Decidere se un sub-fornitore di campionamento è qualificato** (§ 6.6). Atto del QM dopo qualifica formale.
6. **Valutare l'impatto sui risultati** quando un'anomalia di campionamento è scoperta a posteriori. Resta giudizio tecnico del RT (vedi `[[03_RISPOSTA_NC]]`).
7. **Dichiarare chiusa una NC** su § 7.3. Chiusura del QM.
