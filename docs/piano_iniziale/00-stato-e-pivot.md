# Stato del progetto LabNexus / AICertus Sprint 1 e cronologia dei pivot

**Versione**: 1.0
**Stato**: Avvio sviluppo motore + profili, post-pivot 2

---

## Dove siamo

Il pacchetto Sprint 1 è da 20 gettoni complessivi, 3 settimane di calendario. Il gettone è un'unità di complessità relativa, non un'ora; nessuna attività vale meno di un gettone.

Lo Sprint 1 ha avuto un assestamento iniziale e un pivot vero. Sono distinti per natura e vanno letti correttamente per non confondere "evoluzione naturale dell'esplorazione" con "cambio di domanda contrattuale".

### Assestamento iniziale (fine prima settimana)

Durante la prima settimana, il confronto con il cliente (Denis Brazzo, esperto SGQ + ispettore Accredia) e la lettura dei materiali tecnici di LabNexus hanno messo a fuoco una scelta di esecuzione che non era nei dettagli iniziali del lavoro: anziché un agente "puro" in cloud, il valore per Denis stava in uno strumento che lui potesse usare localmente sul suo computer. La domanda contrattuale approvata in offerta non è cambiata:

> *Un agente AI basato su LLM, alimentato dalla KnowledgeBase di LabNexus, produce relazioni conformi a una specifica certificazione che un esperto giudica utilizzabili?*

È rimasta la stella polare. Quello che è cambiato è il **come** ci si arriva: invece di costruire un agente conversazionale generico, costruiamo strumenti locali specializzati per i task ispettivi del QM.

Più che un pivot, è stato un assestamento dettato dall'esplorazione: nei primi giorni di lavoro abbiamo capito meglio cosa serve davvero, e abbiamo aggiustato la traiettoria di esecuzione. Costo dell'assestamento: 5 gettoni della prima fetta di sviluppo, che era allineata all'impostazione iniziale e che è stata archiviata quando l'obiettivo di esecuzione si è ricalibrato.

### Pivot vero (inizio seconda settimana)

Durante una riunione congiunta con il cliente (Stefano Fiorina + Denis Brazzo + CTO), è emerso esplicitamente un vincolo che fino a quel momento era implicito: **tutta l'elaborazione AI deve girare in locale, non solo l'interfaccia**. Questo è il pivot vero della domanda: passa da

> *...alimentato da un LLM in cloud...*

a

> *Un modello locale (Qwen 3 via Ollama, context 128k) produce output di qualità accettabile per Denis come QM sulle capability ispettive del modello AICertus?*

Questa è la domanda contrattuale che lo Sprint 1 deve chiudere entro fine sprint. Il pivot è arrivato dopo l'investimento di 5 gettoni nell'assestamento iniziale, ed è da qui che parte tutto il lavoro fatto dopo e quello che si sta facendo adesso.

## Bilancio gettoni

Il pacchetto è da 20 gettoni complessivi.

### Già speso

| Voce | Gettoni |
|---|---|
| Prima fetta sviluppata prima del pivot, archiviata | 5 |
| Setup ambiente Qwen + tuning context window | 1 |
| Test 1 manuale — revisione PG_RISK_LAB rev04→rev05 — **approvato da Denis** | 1 |
| Test 2 manuale — gestione rilievi Accredia — **approvato da Denis** | 1 |
| **Totale già speso** | **8** |

### Residuo dichiarato al cliente

12 gettoni. Il CTO ha dichiarato apertamente al cliente uno sforo possibile di 1 gettone, quindi il numero comunicato è "13 gettoni effettivi su un residuo formale di 12, con cautela".

### Bilancio interno reale (per il framework)

Il bilancio interno del lavoro tecnico residuo è il seguente:

| Voce | Gettoni |
|---|---|
| Motore eseguibile (incluso provider EUrouter) | 4 |
| 7 profili (1 ciascuno, complessità arrotondata a unità intere) | 7 |
| Prompt meta per generazione profili autonoma | 1 |
| Test indipendente di tenuta su context lungo (chat manuale Ollama, fuori motore) | 1 |
| Report finale + walkthrough conclusivo | 2 |
| **Totale interno** | **15** |

**Sforo interno**: 3 gettoni (15 vs 12 dichiarati). **Il CTO ha scelto di assorbire personalmente lo sforo** come perdita personale, mantenendo il messaggio al cliente coerente con "sforo dichiarato di 1 gettone".

Per il framework: il budget effettivo di lavoro tecnico è 15 gettoni. Trade-off, escalation, decisioni di trasferimento di lavoro a sprint successivi vanno fatti pensando a 15 come tetto operativo, non 12.

## Cosa porteremo a casa a fine sprint

Tre deliverable concreti.

### 1. Eseguibile LabNexus per Mac

Singolo binario Go nativo per macOS Apple Silicon. Esegue una delle 7 capability ispettive del modello AICertus alla volta, chiamando Qwen 3 in locale via Ollama. Tre modalità d'uso: doppio click (dialog macOS interattivo), drag & drop di cartella, riga di comando con parametri. Output markdown in cartella dedicata.

**Ognuna delle sette capability è un task del sistema AICertus completo** (non test artificiosi). Sono i pezzi LLM-critical del sistema futuro, sviluppati e validati uno per uno. Quando lo Sprint 2 costruirà l'orchestrazione e la pipeline deterministica, sarà sopra a questi sette pezzi, non al posto loro.

Le 7 capability:

| Sigla | Capability | Corrisponde a |
|---|---|---|
| A | Revisione documentale dopo cambio di norma | DocAgent + NormAgent |
| B | Gestione rilievi Accredia | CAPAAgent |
| C | Management Review Pack | ReviewAgent |
| D | Checklist d'audit da risk register | AuditAgent |
| E | Equipment alert per scadenze tarature | EquipmentAgent |
| F | Competence gap da matrice competenze | CompetenceAgent |
| G | Analisi z-score PT | PTQCAgent |

### 2. Prompt meta per generazione profili autonoma

Un prompt da incollare in una chat Claude. Guida Claude a produrre un file `<profilo>.yml` valido per LabNexus, a partire dalla descrizione testuale di un nuovo task ispettivo. Serve a dare al cliente autonomia tra Sprint 1 e sprint successivi: può proporre nuovi profili senza dover aspettare il CTO. Il CTO resta nel loop per validazione tecnica e messa in produzione.

### 3. Report di validazione di Qwen

Sintetico e onesto: i 7 esperimenti condotti, gli output prodotti, il giudizio formale del cliente per ciascuno, la conclusione tecnica su Qwen, le raccomandazioni per Sprint 2.

## Ordine di sviluppo dei profili — strategia di rischio crescente

Vincolo non negoziabile. Serve a separare i problemi del codice dai problemi del modello:

1. **Profilo `revisione`** (Capability A) — ricostruisce il Test 1 già approvato. Se l'eseguibile produce output paragonabile, il motore software funziona.
2. **Profilo `rilievi`** (Capability B) — ricostruisce il Test 2 già approvato. Secondo controllo del motore.
3. **Profilo `review-pack`** (Capability C) — prima capability nuova.
4. **Profilo `audit-checklist`** (Capability D) — capability nuova.
5. **Profilo `equipment-alert`** (Capability E) — capability nuova.
6. **Profilo `competence-gap`** (Capability F) — capability nuova.
7. **Profilo `pt-analysis`** (Capability G) — la più rischiosa per Qwen (ragionamento numerico, su cui gli LLM sono storicamente più deboli). Ultima apposta: a quel punto il motore ha sei cicli di shakedown alle spalle, eventuali problemi sono attribuibili al modello non al codice.

Il **prompt meta** (`feat-meta`) viene sviluppato **dopo tutti i 7 profili**, perché astrae il pattern che emerge dalla loro scrittura.

In parallelo a tutto questo: **test indipendente di tenuta su context lungo** (procedura grande 30-50 pagine del SGQ reale di Denis), in chat manuale Ollama, fuori dall'eseguibile. Serve a misurare il comportamento di Qwen vicino al limite di 128k. Si fa dopo i primi 2-3 profili.

## Cosa è dichiaratamente fuori da questo sprint

Tutto ciò che è nei documenti CoWork e WikiLLM di Denis oltre alle 7 capability sopra: orchestrazione automatica, watcher di cartella, trigger su filesystem, pipeline deterministica di indicizzazione del SGQ, multi-tenancy, interfaccia web HITL, database centralizzato. Non sparisce: è la roadmap ufficiale per Sprint 2 e successivi. Materiale prezioso, accettato come piano di crescita, non scartato.

ISO 17034 resta fuori dallo Sprint 1 per decisione del cliente. La KB-ispettore disponibile copre ISO 17025 in modo solido; per 17034 il materiale non è formalizzato. Tutti gli esperimenti si svolgono su perimetro ISO 17025.

## Avvertenza onesta da preservare nei messaggi al cliente

I test che facciamo validano il modello Qwen su input curato. In un sistema completo in produzione, il context viene preparato da una pipeline deterministica (filesystem walk, cross-reference, scadenze): è una capability del sistema, non di Qwen, e va validata separatamente nello Sprint 2.

In sintesi: *"Qwen funziona su queste sette capability con input ben strutturato"* è diverso da *"il sistema funzionerà in produzione"*. Le due cose si validano in due sprint distinti. Va detto apertamente per evitare fraintendimenti.

## Stabilità di scope

Il piano dei prossimi 10 giorni è denso. Il cliente ha accettato il principio della stabilità di scope sullo sprint in corso: nessun nuovo ramo di sviluppo aperto senza una motivazione esplicita basata su dati raccolti.

Se durante lo sviluppo emergono problemi seri sul modello (output non utilizzabili, allucinazioni gravi, capability che non reggono), il CTO ferma tempestivamente e chiama una riunione con il cliente per decidere come proseguire. Il piano può cambiare, ma per ragione esplicita su dati raccolti, non per aggiunta opportunistica di scope.
