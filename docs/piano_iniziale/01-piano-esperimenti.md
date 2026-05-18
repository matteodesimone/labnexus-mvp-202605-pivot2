# LabNexus Sprint 1 — Piano degli esperimenti sul modello Qwen

**Progetto**: LabNexus / AICertus — Sprint 1  
**Versione**: 1.0  
**Data**: {{data}}  
**Stato**: Approvato per esecuzione

---

## Premessa

Questo documento descrive gli esperimenti che lo Sprint 1 conduce per rispondere alla domanda contrattuale:

> *Un modello locale (Qwen 3 via Ollama, context 128k) produce output di qualità accettabile per Denis come QM sulle capability ispettive del modello AICertus?*

Lo scope viene derivato direttamente dai materiali tecnici prodotti da LabNexus (`LabNexus_CoWork_KB.md`, `LabNexus_WikiLLM_Orchestrator.md`, KB-ispettore). Quei materiali descrivono undici agenti specializzati e otto stress test del sistema AICertus completo. Lo Sprint 1 non costruisce quel sistema (è materia di una roadmap pluriennale): valida **quale parte del sistema dipende davvero dal modello LLM** e misura come Qwen 3 in locale si comporta sulle parti che dipendono dal modello.

Tutti i test usano l'eseguibile LabNexus che sviluppiamo nello sprint, salvo il test di tenuta su context lungo che si esegue in chat manuale Ollama per ragioni metodologiche descritte sotto.

---

## Capability già validate

Due capability del sistema AICertus sono già state validate prima dell'inizio dei test ufficiali di sprint, tramite esecuzioni manuali in chat Ollama:

### Capability A — Revisione documentale

Dato un documento del SGQ (es. una procedura) + la versione vecchia di una norma + la versione nuova, il modello produce una bozza di revisione che cita correttamente i cambiamenti normativi, propone modifiche pertinenti, mantiene il tono ispettivo. Corrisponde direttamente agli agenti DocAgent + NormAgent del CoWork KB.

Test eseguito sul caso reale **PG_RISK_LAB rev04 → rev05 con RT-08**. Output **approvato formalmente da Denis**.

### Capability B — Gestione rilievi Accredia

Data una lista di rilievi (NC, Osservazioni, Commenti) con codici, sezioni e testo, il modello produce per ciascuno un piano di gestione strutturato secondo la triade meccanismo-estensione-efficacia, decidendo quando basta una correzione e quando serve anche un'azione correttiva. Corrisponde all'agente CAPAAgent del CoWork KB.

Test eseguito sul caso reale **CSV ACIAA A1**. Output **approvato formalmente da Denis**.

---

## Come scomponiamo il Test 3 proposto da LabNexus

Il Test 3 proposto nei materiali LabNexus mette l'agente AICertus davanti al SGQ reale del laboratorio di Denis come banco prova. La proposta è ricca e importante per due ragioni: usa materiale vero, e contiene una lista articolata di task del sistema futuro.

Scomponendo i task in famiglie, emerge una distinzione utile.

### Famiglia deterministica

Lavori in cui un software classico fa meglio di un LLM, in modo veloce, predicibile, con audit trail perfetto: indicizzazione documentale, mappa documentale, scadenze documenti, scadenze tarature, cross-reference moduli mancanti, riferimenti normativi obsoleti, audit trail, watcher di cartella, incrocio nomi tecnici-fornitori per imparzialità, detection celle Excel non protette, calcoli aritmetici su z-score PT.

**Un LLM qui non aggiunge valore.** Sarebbe più lento, più costoso e meno affidabile. Lo Sprint 1 non li valida (sarebbe inutile) e non li costruisce (sono ingegneria di Sprint 2).

### Famiglia LLM-critical

Lavori in cui il modello fa qualcosa che un software classico non può fare: leggere come legge un ispettore, capire il senso, fare collegamenti, scrivere con tono ispettivo. Lettura ispettiva di una procedura, gestione di un rilievo, valutazione di impatto sistemico, scrittura del Management Review Pack, generazione di checklist d'audit, giudizio di coerenza semantica fra documenti, lettura prudenziale di un requisito quando la KB normativa non lo copre.

**Qui Qwen fa la differenza.** Queste sono le capability che lo Sprint 1 misura.

### Famiglia mista

Alcune cose stanno a cavallo: la parte deterministica prepara il contesto, l'LLM ragiona. Esempio: equipment alert. Il deterministico rileva la scadenza di taratura, l'LLM ragiona sull'impatto sui metodi associati. Lo Sprint 1 valida **la parte LLM** del task, dando per assunto che il deterministico sia un esercizio di Sprint 2.

### Cosa resta a validare nello sprint

Dei sotto-task del Test 3, quelli che ricadono nella famiglia LLM-critical e che lo Sprint 1 misura come capability del modello sono cinque (Capability C, D, E, F, G di seguito). Gli altri ricadono in una di queste situazioni:

- Sono già stati validati come Capability A o B (incorporati come profili dell'eseguibile, ma senza ulteriore scoperta)
- Sono deterministici, non c'è capability LLM da misurare (Sprint 2)
- Sono ibridi con una parte LLM così simile a una capability già misurata che testarli separatamente non aggiunge informazione (es. analisi documentale di un foglio Excel data integrity è una variante di lettura ispettiva, già coperta da Capability A)

---

## I sette esperimenti dello sprint

Sette capability LLM da validare, una per profilo dell'eseguibile LabNexus. Più un test indipendente per la tenuta su context lungo.

### Esperimento 1 — Profilo `revisione` (Capability A)

**Scopo**: validare l'eseguibile sulla Capability A già confermata in chat manuale. Banco di shakedown del motore software.

**Cosa misura**: l'eseguibile riproduce l'output del Test 1 manuale, oppure no. Se sì, il motore è validato. Se no, c'è un problema nel software (prompt assembly, context injection, streaming) da debuggare prima di procedere.

**Materiali utilizzati**: tutti già disponibili — riuso dei file del Test 1 (PG_RISK_LAB rev00 + RT-08 rev03 + RT-08 rev05).

**Esecuzione**: lancio dell'eseguibile con profilo `revisione` sulla cartella di input. L'eseguibile concatena la KB-ispettore rilevante come system context, aggiunge il trigger di istruzione, chiama Qwen via Ollama in streaming, scrive output markdown. Confronto qualitativo con il PG_RISK_LAB rev01 bozza approvato da Denis.

**Cosa ottiene LabNexus**: conferma che la macchina automatizza correttamente il risultato già approvato manualmente.

### Esperimento 2 — Profilo `rilievi` (Capability B)

**Scopo**: validare l'eseguibile sulla Capability B già confermata. Secondo banco di shakedown del motore.

**Cosa misura**: stesso pattern dell'esperimento precedente. Confronto con l'output manuale già approvato.

**Materiali utilizzati**: tutti già disponibili — riuso del CSV ACIAA A1 del Test 2.

**Esecuzione**: lancio dell'eseguibile con profilo `rilievi` sul CSV. Produce un CAPA Pack per ciascun rilievo nel formato del CoWork KB sezione 15.2.

**Cosa ottiene LabNexus**: secondo controllo di non regressione dell'eseguibile rispetto alla chat manuale.

### Esperimento 3 — Profilo `review-pack` (Capability C, nuova)

**Scopo**: validare la capability di generare il Management Review Pack annuale del laboratorio, a partire da dati strutturati. Corrisponde all'agente ReviewAgent del CoWork KB sezione 7.6.

**Cosa misura**: tenuta narrativa del modello su un output composito a 13 sezioni, aderenza ai dati di input senza allucinazioni, coerenza tra sintesi iniziale e decisioni finali della direzione.

**Materiali utilizzati**: input strutturati che rappresentano i dati di un anno operativo di un laboratorio (NC dell'anno, esiti audit interni, reclami, risultati PT con z-score, stato apparecchiature, indicatori KPI).

**Esecuzione**: lancio del profilo `review-pack`. Output: il Pack a 13 sezioni come definito nel template del CoWork KB sezione 15.4.

**Cosa ottiene LabNexus**: misura sulla capability più "lunga" che Qwen affronta. Se l'output regge dalla sezione 1 alla 13 senza degradare, è una capability fortemente impressionante in demo.

### Esperimento 4 — Profilo `audit-checklist` (Capability D, nuova)

**Scopo**: validare la capability di generare una checklist d'audit interno mirata, a partire da un risk register e dallo scope dell'audit. Corrisponde all'agente AuditAgent.

**Cosa misura**: capacità del modello di **ragionare da rischi a azioni ispettive concrete**. Tradurre un registro rischi astratto in domande puntuali che un auditor farebbe sul campo.

**Materiali utilizzati**: un risk register realistico (anche sintetico) con 15-25 voci, più una descrizione di scope di un audit (es. "audit su §6.4 dotazioni con focus su tarature esterne").

**Esecuzione**: lancio del profilo `audit-checklist`. Output: 15-30 domande ispettive raggruppate per area, con campioni documentali da verificare, possibili rilievi anticipati con livello di rischio.

**Cosa ottiene LabNexus**: misura sulla capability di generazione strutturata "non narrativa" (la checklist non è una storia, è una griglia mirata).

### Esperimento 5 — Profilo `equipment-alert` (Capability E, nuova)

**Scopo**: validare la capability di valutare l'impatto tecnico di una scadenza taratura o di un'apparecchiatura fuori stato sui metodi di prova associati. Corrisponde all'agente EquipmentAgent.

**Cosa misura**: capacità del modello di **ragionare causalmente** — *"se questo strumento è fuori stato, allora questi metodi sono coinvolti, allora questi rapporti potrebbero dover essere congelati"*. È un task di reasoning su parti tecniche e metrologiche.

**Materiali utilizzati**: una scheda apparecchiatura realistica (codice, descrizione, ultima taratura, prossima scadenza, fornitore, **metodi di prova associati**) e una descrizione testuale dell'evento scatenante (es. *"certificato di taratura rientrato con NC su intervallo X-Y"*).

**Esecuzione**: lancio del profilo `equipment-alert`. Output: Equipment Alert nel formato del template CoWork KB sezione 15.3 (stato attuale, rischio tecnico, azioni proposte, bozza email fornitore, checklist al rientro).

**Cosa ottiene LabNexus**: misura su una capability che combina lettura tecnica + ragionamento causale + scrittura ispettiva.

### Esperimento 6 — Profilo `competence-gap` (Capability F, nuova)

**Scopo**: validare la capability di analizzare un gap di competenze quando una nuova procedura o metodo entra in scope. Corrisponde all'agente CompetenceAgent.

**Cosa misura**: capacità del modello di **leggere una matrice strutturata** (chi è autorizzato a cosa) e collegarla a requisiti procedurali nuovi.

**Materiali utilizzati**: una matrice competenze realistica (anche piccola, 10-15 righe × 8-12 colonne con autorizzazioni/qualifiche), più una descrizione di una procedura o metodo nuovo che richiede una certa competenza.

**Esecuzione**: lancio del profilo `competence-gap`. Output: gap report (chi è autorizzato e chi no), piano formazione proposto, autorizzazioni da rilasciare/aggiornare.

**Cosa ottiene LabNexus**: misura sulla capability di ragionare su griglie di dati strutturati, che è metodologicamente diversa da revisione (lettura testuale) e da review-pack (composizione narrativa).

### Esperimento 7 — Profilo `pt-analysis` (Capability G, nuova)

**Scopo**: validare la capability di analizzare risultati PT con z-score, classificare il rischio del metodo, proporre azioni. Corrisponde all'agente PTQCAgent.

**Cosa misura**: capacità del modello su **ragionamento numerico unito a ragionamento ispettivo**. Questa è la capability più rischiosa del sprint, perché gli LLM sono storicamente meno solidi su task aritmetici e statistici rispetto al testo. Lo sviluppo di questo profilo viene per ultimo apposta: a quel punto il motore software ha sei cicli di shakedown alle spalle e isoliamo eventuali problemi del modello da problemi del codice.

**Materiali utilizzati**: risultati di una o due campagne PT con z-score per round e matrice, elenco dei metodi del laboratorio associati a quei PT, idealmente storico di 2-3 anni di z-score sullo stesso PT per testare l'identificazione di trend.

**Esecuzione**: lancio del profilo `pt-analysis`. Output: Investigation Pack con riassunto risultati, classificazione rischio per metodo, trend storico, azioni immediate per z-score fuori soglia.

**Cosa ottiene LabNexus**: misura sulla capability che potrebbe segnare i limiti di Qwen. Se va bene, lo Sprint 1 chiude con un dato fortissimo. Se va male, è informazione preziosa: indica che in Sprint 2 la parte numerica andrà appoggiata a pre-processing deterministico, non lasciata al modello.

### Esperimento indipendente — Stress test su context lungo

**Scopo**: misurare la tenuta di Qwen quando il context si avvicina al limite dichiarato di 128k token.

**Perché fuori dall'eseguibile**: si esegue in chat manuale Ollama per **isolare la variabile "context lungo"** senza interferenze del software. Se inserito nell'eseguibile, eventuali problemi diventerebbero ambigui (è il modello o è il codice?).

**Cosa misura**: profondità sostenuta del modello (analizza o solo parafrasa), gap reali trovati confrontati con un'analisi manuale, riferimenti obsoleti citati nelle pagine finali del documento (banco di prova della tenuta), eventuale degrado visibile dopo metà output.

**Materiali utilizzati**: una procedura "pesante" del SGQ reale di Denis, idealmente 30-50 pagine con molti moduli e requisiti richiamati.

**Scope normativo**: il test resta sul perimetro ISO 17025 ed esclude esplicitamente ISO 17034. La KB-ispettore disponibile copre 17025 in modo solido; per 17034 il materiale non è formalizzato. Denis ha deciso di mantenere 17034 fuori da questo sprint.

**Esecuzione**: caricamento della procedura + KB-ispettore completa nella chat Ollama. Prompt: review ispettiva completa con sintesi esecutiva, mappatura normativa, gap, riferimenti interni, riferimenti obsoleti, proposta operativa.

**Cosa ottiene LabNexus**: una valutazione qualitativa della profondità sostenuta del modello su input molto grandi. Risponde alla domanda *"fino a che lunghezza Qwen ragiona bene?"*, fondamentale per dimensionare la pipeline deterministica dello Sprint 2.

---

## Materiali richiesti a Denis

Tutto in una lista, da consegnare nei tempi compatibili con l'avanzamento dello sprint (i primi profili sono già coperti dai materiali esistenti, gli altri arrivano via via).

### Materiali nuovi

1. **Procedura "pesante" del SGQ reale** per l'esperimento indipendente di context lungo. Idealmente 30-50 pagine, con molti moduli e requisiti richiamati. In `.docx` o `.pdf`.

2. **Materiali per `review-pack` (Capability C)**: dati strutturati per simulare un riesame annuale. Vanno bene: tabelle di NC dell'anno, esiti audit interni, reclami clienti, risultati PT con z-score, stato apparecchiature, indicatori KPI. **In alternativa**, lo schema delle tabelle/CSV vuoti con le colonne usate tipicamente, da cui costruisco io dati sintetici plausibili.

3. **Materiali per `audit-checklist` (Capability D)**: un risk register realistico (15-25 voci) + descrizione di scope di un audit. **In alternativa**, lo schema del risk register tipico (colonne, tipologia di voci).

4. **Materiali per `equipment-alert` (Capability E)**: una scheda apparecchiatura realistica (con metodi di prova associati) + descrizione di un evento scatenante (taratura scaduta, certificato con NC). **In alternativa**, formato delle schede + descrizione testuale degli eventi tipici.

5. **Materiali per `competence-gap` (Capability F)**: una matrice competenze (10-15 righe × 8-12 colonne) + una procedura/metodo che richiede una certa competenza. **In alternativa**, lo schema della matrice usata tipicamente.

6. **Materiali per `pt-analysis` (Capability G)**: risultati di una o due campagne PT con z-score + elenco metodi associati + idealmente storico 2-3 anni. **In alternativa**, formato delle tabelle.

### Materiali già disponibili

Per i profili `revisione` e `rilievi` (Capability A e B) i materiali sono già in nostre mani, riusati dai Test 1 e 2.

### Nota su materiali reali vs. sintetici

Per la sostanza degli esperimenti, **materiali reali anonimizzati sono preferibili a materiali sintetici** perché rendono le valutazioni di Denis più solide. Tuttavia, se l'estrazione e l'anonimizzazione di materiali reali è onerosa per Denis, gli **schemi delle tabelle e descrizioni di pattern tipici** sono sufficienti: da quelli costruisco io dati sintetici plausibili sul modello dei Test 1 e 2.

L'unica eccezione è l'esperimento indipendente di context lungo: lì la procedura **deve essere reale e di dimensioni rappresentative**, perché stiamo testando come Qwen gestisce un documento di complessità realistica, non una procedura inventata.

### Scope ISO 17034

Su decisione di Denis, ISO 17034 resta fuori da questo sprint. Tutti gli esperimenti si svolgono su perimetro ISO 17025 (per il quale la KB-ispettore è completa). Una eventuale validazione di Qwen su 17034 sarà oggetto di uno sprint successivo, previa formalizzazione di una KB-ispettore dedicata.

---

## Sintesi dei deliverable di validazione

A fine sprint, ognuno dei sette esperimenti produce:

- L'**output di Qwen** sul caso di test (file markdown, allegato al report finale)
- La **valutazione di Denis** sull'output, secondo griglia condivisa (allucinazioni, qualità ispettiva, profondità, tenuta, utilizzabilità pratica)
- Lo **stato** di ciascuna capability: validata / validata con riserva / non validata

L'esperimento indipendente di context lungo produce inoltre una **misura qualitativa del comportamento di Qwen vicino al limite di context window**, dato fondamentale per il dimensionamento dello Sprint 2.

Insieme, questi otto risultati rispondono in modo completo alla domanda contrattuale dello Sprint 1.
