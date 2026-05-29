---
codice: A8_gestione_NC_da_accreditatore
tipo: appendice_tecnica
livello: misto
titolo: "Appendice A8 — Gestione operativa delle NC ricevute da Accredia"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "8.7" }
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.10" }
  - { doc: "RT-08", rev: "05", paragrafo: "8.7" }
  - { doc: "RG-02", rev: "08", note: "regolamento accreditamento — gestione esiti e ricorsi" }
fonti_secondarie:
  - { doc: "RT-23", rev: "05", note: "scope, applicabile a NC su variazioni di scope" }
  - { doc: "Pol-RemAss", rev: "00", note: "verifiche da remoto" }
  - { doc: "Racc-CdIG", rev: "00" }
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[00_GLOSSARIO]]"
  - "[[02_PROTOCOLLO_ASK_BEFORE_ANSWER]]"
  - "[[03_PROTOCOLLO_HANDOFF_QM]]"
  - "[[01_COME_PENSA_ISPETTORE]]"
  - "[[02_TRE_LIVELLI_DI_CONTROLLO]]"
  - "[[03_RISPOSTA_NC]]"
  - "[[04_ERRORI_FATALI]]"
  - "[[sez_7_10_attivita_non_conformi]]"
  - "[[sez_8_7_azioni_correttive]]"
  - "[[A6_audit_interni_riesame_direzione]]"
tags: [NC_accredia, risposta_NC, ricorso, tempistiche, efficacia, attuazione]
escalation_qm: true
---

> **Marcatore di livello: 🟢🟡🟣** — Appendice trasversale. Contenuti normativi (livello 1+2), prassi ispettiva (livello 3), best practice tecnica (livello 4). Per la **struttura della risposta tecnica** ad una NC vedi `[[03_RISPOSTA_NC]]`. Questa appendice copre il **versante gestionale-procedurale**: come si organizza il laboratorio dal ricevimento del rilievo alla chiusura formale presso Accredia.

# Appendice A8 — Gestione operativa delle NC ricevute da Accredia

> **Riassunto operativo.** Quando arriva un rilievo da Accredia, due cose si attivano in parallelo: il lavoro **tecnico** (capire la causa, valutare estensione e impatto, definire azione correttiva, verificare efficacia — vedi `[[03_RISPOSTA_NC]]`) e il lavoro **gestionale-procedurale** (rispettare tempistiche, formato di risposta, allegati richiesti, gestione dell'eventuale ricorso). Questa appendice è dedicata al secondo: cosa fa il QM dal "rilievo arrivato" alla "NC chiusa da Accredia". L'agente la carica appena il QM dice "abbiamo ricevuto NC", "verbale di visita", "ricorso", "risposta scaduta", "sospensione".

---

## 1. Principio normativo (🟢 LIVELLO 1+2)

### 1.1 ISO 17025:2018 § 8.7 e § 7.10

> Testo a pagamento, non riprodotto. Sintesi operativa:

- **§ 8.7 Azioni correttive.**
  - 8.7.1 — quando si verifica una NC, il laboratorio deve: (a) reagire alla NC e prendere azioni per controllarla e correggerla, gestendone le conseguenze; (b) valutare la necessità di azioni per eliminare la causa, affinché non si ripeta o si verifichi altrove (review NC, determinazione cause, identificazione di NC simili o potenziali).
  - 8.7.2 — riesaminare l'efficacia delle azioni correttive.
  - 8.7.3 — aggiornare i rischi e le opportunità (§ 8.5) se necessario.
  - 8.7.4 — apportare modifiche al SGQ se necessario.
  - 8.7.5 — conservare registrazioni della NC, delle azioni, dei loro esiti.
- **§ 7.10 Attività non conformi.** Procedure per identificare attività di laboratorio non conformi a requisiti propri o del cliente, includere: responsabilità e autorità per la gestione, azioni basate sui rischi, valutazione del significato (eventuale impatto sui risultati emessi), decisione sull'accettabilità, decisione su comunicazione cliente, responsabilità per autorizzare la ripresa.

> Una NC ricevuta da Accredia è insieme una **NC del SGQ** (su cui si applica § 8.7) e tipicamente anche un'**attività non conforme** o un suo segnale (su cui può applicarsi § 7.10 quando il rilievo tocca un risultato/processo tecnico).

### 1.2 RT-08 rev.05 § 8.7

> **Fonte**: `[RT-08 rev.05 § 8.7 — EC 10-02-2022]`

RT-08 si applica integralmente con enfasi su:
- l'azione correttiva deve indirizzare la **causa**, non l'effetto;
- l'**efficacia** va verificata, non solo dichiarata;
- la registrazione di NC, azioni correttive, esiti e verifica di efficacia è obbligatoria.

### 1.3 RG-02 rev.08 — Regolamento per l'accreditamento

> **Fonte**: `[RG-02 rev.08 (Accredia, free)]`

`RG-02` è il documento contrattuale che lega il laboratorio ad Accredia. Su gestione esiti delle visite e dei rilievi:
- Accredia comunica i rilievi al laboratorio al termine della verifica e formalmente nel rapporto/verbale di visita.
- Il laboratorio dispone di un **periodo definito da RG-02** per trasmettere ad Accredia la risposta (analisi causa + piano azioni + evidenza di attuazione + verifica di efficacia, secondo lo stato richiesto per la tipologia di rilievo).
- Accredia riesamina la risposta e dichiara la NC **chiusa** se la risposta è adeguata; **non chiusa** se incompleta o non convincente, con eventuali richieste di integrazione entro nuovo termine.
- In caso di **mancato rispetto delle tempistiche** o di **risposte non adeguate ripetutamente**, Accredia può attivare provvedimenti che vanno fino alla **sospensione** o **revoca** dell'accreditamento (in tutto o per categorie specifiche).
- È prevista una procedura di **ricorso** (appeal) contro le decisioni di Accredia, secondo termini e modalità di RG-02.

> Le tempistiche puntuali, la classificazione dei rilievi (osservazioni / non conformità / rilievi gravi / situazioni di immediato impatto sull'accreditamento) e la procedura di ricorso vanno **lette nella versione vigente di RG-02** sul sito Accredia. L'agente cita il quadro, il QM verifica i termini esatti.

### 1.4 Documenti correlati

- **`Pol-RemAss`** — verifiche da remoto: alcune NC possono derivare da verifiche in remoto; la risposta segue gli stessi requisiti formali.
- **`Racc-CdIG`** — raccomandazioni per omogeneità ispettiva: utili per anticipare il modo in cui un riesaminatore Accredia leggerà la risposta.
- **`RT-23`** — quando la NC riguarda il CAdA (es. attività dichiarata accreditata fuori scope), la risposta tocca anche la modifica formale dello scope.

---

## 2. Come un ispettore lo verifica e come si imposta la risposta (🟡 LIVELLO 3 — prassi ispettiva)

> Questa sezione è specifica della prassi ispettiva sulla **gestione** delle NC da Accredia. Per la **struttura della risposta tecnica** (causa, estensione, impatto, correzione, efficacia) si rimanda integralmente a `[[03_RISPOSTA_NC]]`.

### 2.1 Tipologie di rilievi tipiche

In una visita Accredia si distinguono comunemente:

| Tipologia | Significato operativo | Azione tipica |
|---|---|---|
| **Osservazione / Raccomandazione** | Segnalazione di un punto di miglioramento, non una violazione del requisito | Trattamento interno (registro miglioramento, valutazione opportunità). Non richiede risposta formale "chiudibile" come una NC, ma è bene rispondere e tracciare per coerenza |
| **Non Conformità (NC)** | Violazione di un requisito (ISO/RT-08/RG/regolamento) | Risposta formale entro tempi RG-02: causa, azione correttiva, piano attuazione, verifica efficacia |
| **NC grave / rilievo grave** | Violazione che incide sull'affidabilità del dato, sulla riferibilità, sulla competenza o sull'integrità del SGQ | Risposta rapida (tempi più stretti); azione immediata di contenimento (es. sospensione interna dell'attività interessata); verifica di efficacia particolarmente solida |
| **Situazione di immediato impatto sull'accreditamento** | Es. mancata copertura di prove dichiarate accreditate, perdita di personale chiave, perdita di sede, scadenza non gestita | Trattata da Accredia secondo RG-02 in tempi e modalità specifiche; può portare a sospensione parziale/totale |

> **Importante**: la classificazione finale del rilievo è di Accredia. L'agente e il QM possono interpretare, ma la categoria definitiva è quella indicata nel verbale di visita o nella comunicazione formale.

### 2.2 Domanda 1 — Meccanismo (come il lab gestisce la NC ricevuta)

**"Mostratemi la procedura interna di gestione delle NC ricevute da Accredia: chi riceve la comunicazione, chi apre la pratica, chi è responsabile della risposta, chi coinvolge RT/Direzione, come si gestiscono le tempistiche, come si protocolla la risposta verso Accredia?"**

Analisi ispettiva:
- Cerco una **procedura SGQ** dedicata, separata da quella delle NC interne ma raccordata (le NC Accredia diventano comunque NC del SGQ ex § 8.7).
- Verifico che la **responsabilità** sia attribuita (QM coordina, RT contribuisce sul tecnico, Direzione firma). Se "la fa il QM da solo" su una NC grave, c'è un problema di governance.
- Verifico le **tempistiche interne**: il lab deve avere un piano che rispetti le scadenze RG-02 con margine, non con sovrapposizione all'ultimo giorno.

### 2.3 Domanda 2 — Estensione (storico NC e analisi di pattern)

**"Mostratemi lo storico delle NC ricevute da Accredia negli ultimi due cicli di accreditamento. Esistono NC ricorrenti su stesso requisito? La causa è stata davvero diversa ogni volta, o c'è una causa di sistema che continua a manifestarsi?"**

Analisi ispettiva:
- La **NC ricorrente** è uno dei segnali più forti di azione correttiva debole. Se ho NC sul § X.Y in tre visite consecutive, le azioni precedenti hanno aggiustato il sintomo, non la causa.
- Verifico se l'**analisi di pattern** è stata fatta — di norma in sede di riesame di direzione (vedi `[[A6_audit_interni_riesame_direzione]] § 1.2 input "valutazioni organismi esterni"`).

### 2.4 Domanda 3 — Efficacia (verifica post-chiusura della NC)

**"Per la NC chiusa più recente, mi mostrate la verifica di efficacia su casi reali dopo la chiusura formale? Quanto tempo dopo? Su che casistica? Quale audit interno o controllo ha confermato che la barriera nuova funziona?"**

Analisi ispettiva:
- La NC è chiusa quando Accredia accetta la risposta, **ma l'efficacia si vede dopo**, su casi reali. Cerco evidenza che il laboratorio abbia continuato a monitorare il punto debole anche dopo la chiusura formale.
- Verifico se gli **audit interni successivi** hanno toccato l'area della NC chiusa per confermarne la guarigione.

### 2.5 Pattern di errori comuni nella gestione

1. **Risposta in extremis**: inviata l'ultimo giorno utile, senza verifica interna. Spesso debole e incompleta. NC potenziale "non chiusa" alla prima lettura Accredia, con conseguente nuovo termine e perdita di tempo.
2. **Verifica di efficacia "concomitante"**: dichiarata efficace lo stesso giorno della formazione/aggiornamento procedura. Vedi `[[03_RISPOSTA_NC]] Errore 4`. Accredia di norma chiede tempo di osservazione su casistica reale.
3. **Causa = "errore umano"**: vedi `[[03_RISPOSTA_NC]] Errore 1`. Accredia spesso respinge.
4. **Mancata valutazione dell'estensione e dell'impatto sui risultati emessi** quando la NC è tecnica. Vedi `[[03_RISPOSTA_NC]] Errori 2-3`.
5. **Risposta "narrativa" senza evidenze**: testo lungo, allegati zero. La risposta deve essere supportata da documenti (procedure aggiornate, registrazioni nuove, output di audit di efficacia).
6. **Ricorso impulsivo**: il ricorso (appeal) ha senso quando esiste un fondato dubbio interpretativo o procedurale dell'ispettore. Usarlo come prima reazione difensiva è di norma controproducente. Vedi § 3.5.
7. **Mancata comunicazione interna**: il QM gestisce la NC senza informare RT/Direzione/personale interessato. Risultato: la NC ricorre perché il personale non sa che il punto era stato rilevato.
8. **Sottovalutazione della NC grave**: trattamento come fosse "una NC normale", senza azioni immediate di contenimento. Errore fatale.

---

## 3. Best practice tecnica e operativa (🟣 LIVELLO 4 — best practice)

### 3.1 Workflow standard "Ricevimento NC → Chiusura Accredia"

```
GIORNO 0 — Ricevimento verbale/comunicazione Accredia
    └─> QM protocolla la comunicazione
    └─> QM apre pratica nel registro NC Accredia
    └─> QM convoca riunione di analisi (entro 3-5 gg) con RT + Direzione
    └─> QM stima preliminare di tempistiche e responsabilità

GIORNO 1-7 — Analisi del rilievo
    └─> Lettura attenta del testo del rilievo (cosa Accredia ha osservato esattamente)
    └─> Identificazione dei § ISO/RT-08/RG coinvolti
    └─> Classificazione interna della tipologia (osservazione / NC / NC grave / impatto immediato)
    └─> Per NC grave o con impatto tecnico: azione immediata di contenimento (es. sospensione attività, isolamento risultati)

GIORNO 5-15 — Analisi di causa
    └─> Applicare metodo di analisi causa (5 Why, Ishikawa, FTA — quello che il SGQ adotta)
    └─> NON fermarsi a "errore umano"; cercare la barriera di sistema mancata
    └─> Coinvolgere personale del banco quando rilevante

GIORNO 10-25 — Definizione azione correttiva
    └─> Azione che colpisce la causa identificata, non l'effetto
    └─> Responsabile + tempi + risorse
    └─> Modifiche al SGQ (procedure, moduli, autorizzazioni) se necessarie
    └─> Aggiornamento registro rischi-opportunità se applicabile (§ 8.7.3)

GIORNO 20-30+ — Attuazione e verifica di efficacia
    └─> Attuazione delle azioni decise (con evidenza)
    └─> Verifica di efficacia su casi reali (NON concomitante)
    └─> Quando il tempo di osservazione richiede mesi, dichiarare nella risposta il piano di verifica di efficacia con timing realistico

ENTRO TERMINE RG-02 — Trasmissione risposta ad Accredia
    └─> Struttura: descrizione rilievo + analisi causa + estensione + impatto + correzione + azione correttiva + verifica di efficacia (vedi 03_RISPOSTA_NC)
    └─> Allegati: procedure aggiornate, registrazioni nuove, output verifica efficacia, comunicazioni cliente se applicabili
    └─> Firma Direzione + QM (+ RT se tecnica)
    └─> Protocollo in uscita

ATTESA RIESAME ACCREDIA
    └─> Eventuale richiesta integrazione → rispondere entro nuovo termine con evidenza aggiuntiva
    └─> Comunicazione di chiusura → archiviazione nel SGQ
    └─> Aggiornamento registro NC Accredia
    └─> Input al successivo riesame di direzione (§ 8.9.2 — valutazioni organismi esterni)

POST CHIUSURA
    └─> Audit interno mirato sull'area della NC entro 6-12 mesi per confermare efficacia nel tempo
    └─> Verifica che la NC non ricorra alla visita successiva
```

### 3.2 Registro NC Accredia — struttura minima

| ID interno | Riferimento Accredia (verbale, n.) | Data ricevimento | Tipologia (oss/NC/NC grave) | § ISO/RT/RG | Sintesi rilievo | Responsabile risposta | Termine RG-02 | Data invio risposta | Esito riesame Accredia | Data chiusura | Audit interno post-chiusura (data, esito) | Note |
|---|---|---|---|---|---|---|---|---|---|---|---|---|

### 3.3 Struttura standard della risposta scritta ad Accredia

```
INTESTAZIONE
- Riferimento Accredia (visita, numero rilievo)
- Tipologia rilievo (come da verbale)
- § ISO/RT/RG citati

1. DESCRIZIONE DEL RILIEVO
   Ricognizione neutra, senza minimizzare né negare.

2. ANALISI DI CAUSA
   Metodo applicato (5 Why / Ishikawa / FTA). Identificazione della barriera di sistema mancata.

3. ESTENSIONE
   Ricerca del problema in altri operatori, altri metodi, altri periodi, altre matrici. Evidenza.

4. IMPATTO SUI RISULTATI EMESSI (se applicabile)
   Identificazione rapporti potenzialmente impattati. Riesame tecnico. Decisione (difendibili / da rivalutare / da comunicare al cliente).

5. AZIONE IMMEDIATA (CORREZIONE)
   Cosa fatto subito per contenere. Eventuale sospensione attività.

6. AZIONE CORRETTIVA SISTEMICA
   Nuova barriera introdotta. Come intercetta il problema. Modifiche al SGQ.

7. VERIFICA DI EFFICACIA
   Su quali casi reali, con quale tempo di osservazione, con quale risultato. Se ancora in corso al momento della risposta, dichiarare il piano.

8. RESPONSABILE E TIMELINE
   Chi ha fatto cosa, quando.

9. ALLEGATI
   - Procedure aggiornate (con n. revisione, data)
   - Registrazioni nuove (esempi)
   - Output di audit di efficacia
   - Comunicazioni cliente (se applicabile)

FIRME: Direzione, QM, RT (su NC tecnica)
```

### 3.4 Note specifiche sulla verifica di efficacia

`RT-08` ed `RG-02` non prescrivono un tempo minimo universale per la verifica di efficacia, ma la prassi consolidata e i `Racc-CdIG` indicano:
- per NC su attività ripetuta (esecuzione metodo, compilazione moduli): osservazione su **almeno 30-45 giorni** di operatività su casistica reale,
- per NC su attività rara (audit interno, riesame di direzione): osservazione sul **ciclo successivo** dell'attività,
- per NC tecnica con impatto sui risultati: verifica con **PT/ILC dedicato** o **ri-controllo su campioni storici/spike** dove possibile.

Se la verifica di efficacia richiede tempi che eccedono il termine RG-02, nella risposta si **dichiara il piano** e si chiede ad Accredia di accettare il completamento successivo (Accredia di norma accetta se il piano è motivato e tracciato).

### 3.5 Ricorsi (appeal)

Il **ricorso** è un istituto previsto da RG-02 e si propone quando il laboratorio ritiene che:
- il rilievo non sia tecnicamente fondato,
- l'interpretazione del requisito da parte dell'ispettore sia errata,
- la classificazione (es. NC vs osservazione, NC vs NC grave) sia eccessiva.

Buone prassi:
1. **Prima leggere bene** il rilievo: spesso quello che sembra ingiusto è in realtà espresso in modo sintetico ma corretto.
2. **Coinvolgere consulente esperto o ispettore esterno** prima di formalizzare il ricorso.
3. **Argomentare tecnicamente**: il ricorso non è "non siamo d'accordo", è "ecco perché tecnicamente il rilievo non si applica al nostro caso, e i riferimenti normativi sono questi".
4. **Tempi RG-02**: il ricorso ha termini propri, da rispettare.
5. **Rapporto con Accredia**: il ricorso è uno strumento legittimo, ma usarlo in modo difensivo a fronte di rilievi sostanzialmente corretti danneggia la credibilità del laboratorio nel lungo periodo.

### 3.6 Quando la NC tocca lo scope (RT-23)

Se la NC riguarda attività dichiarate accreditate fuori dallo scope effettivo (es. matrice non coperta, parametro non in CAdA, campionamento dichiarato accreditato ma in CAdA "sole prove"), la risposta deve includere:
- **Sospensione immediata** dell'attività dichiarata fuori scope,
- **Comunicazione al cliente** se il rapporto è già stato emesso (decisione Direzione + RT),
- **Eventuale richiesta di modifica del CAdA** verso Accredia se l'attività è strutturalmente nello scope di interesse del laboratorio (RG-02 + RT-23),
- **Verifica** che non esistano altri rapporti emessi nella stessa condizione di errore.

### 3.7 Comunicazione interna

Una NC ricevuta va comunicata al personale interessato, di norma:
- **Riunione del settore** coinvolto entro 7-15 giorni dal ricevimento,
- **Coinvolgimento attivo** del personale del banco nell'analisi di causa,
- **Aggiornamento del personale** al termine, con la nuova procedura/barriera spiegata e formata,
- **Registrazione della comunicazione** nel SGQ.

Senza comunicazione interna la NC ricorre.

---

## 4. Template operativo per il QM

### 4.1 Checklist "Primo giorno dopo ricevimento NC Accredia"

- [ ] Protocollare la comunicazione in entrata
- [ ] Aprire pratica nel registro NC Accredia (ID univoco)
- [ ] Identificare § ISO/RT/RG citati e classificare la tipologia (osservazione / NC / NC grave / impatto immediato)
- [ ] Calcolare e annotare il termine RG-02 di risposta
- [ ] Convocare riunione di analisi (QM + RT + Direzione + personale chiave) entro 3-5 giorni
- [ ] Per NC grave o con potenziale impatto tecnico immediato: valutare azione di contenimento (es. sospensione interna attività)
- [ ] Definire responsabile della risposta (di norma QM coordina) e responsabili contributori
- [ ] Stilare bozza di cronoprogramma interno con margine sul termine RG-02

### 4.2 Checklist "Verifica di adeguatezza della risposta prima dell'invio"

- [ ] Descrizione del rilievo neutra e fedele al verbale
- [ ] Analisi di causa NON ridotta a "errore umano"; barriera di sistema esplicitata
- [ ] Estensione: ricerca del problema in altri operatori/strumenti/matrici/periodi
- [ ] Impatto sui risultati emessi valutato (se tecnica)
- [ ] Azione immediata documentata
- [ ] Azione correttiva sistemica che introduce nuova barriera
- [ ] Verifica di efficacia su casi reali (con tempo di osservazione) o piano motivato di verifica futura
- [ ] Allegati completi (procedure, registrazioni, output audit di efficacia, comunicazioni cliente)
- [ ] Firme corrette (Direzione + QM + RT se tecnica)
- [ ] Protocollo in uscita generato
- [ ] Invio entro termine RG-02 con margine

### 4.3 Format "Comunicazione interna NC Accredia ricevuta"

```
A: [personale interessato]
Da: [QM]
Oggetto: NC Accredia [ID] su § [...] — coinvolgimento personale del settore

CONTESTO
Nel verbale Accredia del [data] è stata sollevata la seguente NC: [testo].

PERCHÉ COINVOLGO IL SETTORE
La causa va indagata sul flusso operativo. Vi chiedo collaborazione per:
- Descrivermi onestamente come l'attività è eseguita oggi al banco
- Identificare insieme dove il sistema attuale può aver permesso il rilievo
- Costruire l'azione correttiva insieme, perché chi opera al banco la deve sentire propria

RIUNIONE PREVISTA: [data, ora]

NESSUNA CACCIA ALL'ERRORE
Il rilievo è del SISTEMA, non della persona. L'obiettivo è capire perché il sistema lo ha permesso, non chi ha sbagliato.
```

### 4.4 Action register post-chiusura NC Accredia

| ID NC | Azione di follow-up | Tempistica | Responsabile | Evidenza attesa | Stato |
|---|---|---|---|---|---|
| | Audit interno mirato sull'area | 6-12 mesi post chiusura | QM | Rapporto audit con conferma efficacia | |
| | Verifica al riesame di direzione successivo | Prossimo riesame | QM/Direzione | Verbale riesame con punto specifico | |
| | Inserimento del tema nel piano audit pluriennale | Prossimo piano | QM | Piano audit aggiornato | |

---

## 5. Quando l'agente attiva questa appendice

1. Il QM segnala *"abbiamo ricevuto NC da Accredia"* / *"è arrivato il verbale"* / *"dobbiamo rispondere entro…"*.
2. Si sta preparando la **bozza di risposta** ad un rilievo Accredia: l'agente assiste nell'istruttoria (combinando questa appendice con `[[03_RISPOSTA_NC]]`).
3. Si valuta l'opportunità di un **ricorso** (appeal).
4. Una NC Accredia chiusa ricorre alla visita successiva: serve riesame critico dell'efficacia precedente.
5. Si pianifica l'**audit interno mirato post-chiusura** (vedi `[[A6_audit_interni_riesame_direzione]]`).
6. Il riesame di direzione tratta le "valutazioni di organismi esterni" come input (§ 8.9.2).

---

## 6. Cosa l'agente NON può chiudere su questo tema (handoff QM obbligatorio)

Vedi `[[03_PROTOCOLLO_HANDOFF_QM]] § 4 e § 5`. In materia di NC ricevute da Accredia, l'agente **non** può:

1. **Firmare la risposta** ad Accredia. Firma Direzione + QM (+ RT se tecnica).
2. **Inviare la risposta** ad Accredia. È **comunicazione esterna formale**: atto del QM/Direzione (vedi `[[03_PROTOCOLLO_HANDOFF_QM]] § 7`).
3. **Decidere la classificazione interna** del rilievo (NC tecnica vs sistemica, gravità, urgenza interna): atto del QM con RT/Direzione.
4. **Decidere di proporre ricorso** (appeal) verso Accredia. Decisione di Direzione, di norma con consulente esperto.
5. **Decidere comunicazione al cliente** per rapporti potenzialmente impattati. Decisione Direzione + QM + RT.
6. **Dichiarare la NC "chiusa"** prima della comunicazione formale di chiusura da parte di Accredia. La chiusura è atto di Accredia; internamente la pratica resta aperta finché non arriva il riscontro.
7. **Dichiarare efficace** l'azione correttiva senza verifica su casi reali. La dichiarazione è del QM dopo evidenze.
8. **Sospendere o modificare il CAdA** in risposta a una NC che tocca lo scope. È atto formale verso Accredia (RG-02 + RT-23) di responsabilità della Direzione.

> **Nota di sistema**: la gestione di una NC Accredia è il momento in cui il SGQ si misura con il mondo reale. L'agente è uno strumento di istruttoria, di memoria documentale e di coerenza fra parti. La **risposta** è un atto del laboratorio. Tenere fermi questi confini è parte della credibilità della risposta stessa.
