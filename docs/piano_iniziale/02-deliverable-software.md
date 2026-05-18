# LabNexus — Eseguibile per Mac

**Progetto**: LabNexus / AICertus — Sprint 1  
**Componente**: Eseguibile per Mac, denominato *LabNexus*  
**Versione descrittiva**: 1.0  
**Data**: {{data}}

---

## Cosa è

LabNexus è un eseguibile per Mac che esegue una capability ispettiva del modello AICertus alla volta, chiamando il modello AI Qwen 3 in locale. Produce output in markdown nella cartella indicata.

L'eseguibile viene utilizzato come piattaforma per i sette esperimenti dello sprint (descritti nel documento di piano degli esperimenti) e viene rilasciato come deliverable a fine sprint per uso operativo del QM. Gli esperimenti di validazione del modello e i test funzionali dell'eseguibile coincidono: ogni capability validata corrisponde a un profilo dell'eseguibile.

## Come si usa

L'eseguibile riconosce automaticamente la modalità d'uso in base al contesto.

**Doppio click sull'icona**. Si apre un dialog macOS nativo che chiede quale capability eseguire (lista dei profili disponibili con descrizione di una riga ciascuno), poi la cartella di input, poi la cartella di output. L'esecuzione avviene con finestra console aperta che mostra log in tempo reale e progress bar. Al termine, il Finder si apre automaticamente sulla cartella dei risultati.

**Trascinamento di una cartella sull'icona**. Equivale al doppio click, ma con la cartella di input già selezionata.

**Riga di comando**. Si passano i parametri (capability, cartella input, cartella output) e l'esecuzione parte direttamente. Utile per scripting o per ripetere lo stesso test su input differenti.

## Come funziona dentro

Per ciascuna capability esiste un file di configurazione che descrive:

- Quali file della KB-ispettore di Denis caricare come **contesto** per il modello (l'expertise ispettiva scritta da Denis, letta direttamente dal filesystem locale)
- Quale **istruzione di innesco** dare al modello per il task specifico
- Quali parametri tecnici utilizzare

A ogni esecuzione l'eseguibile rilegge la KB-ispettore dal filesystem. Questo significa che le iterazioni di Denis sulla sua KB vengono incorporate automaticamente nei lanci successivi, senza ricompilazione, senza disallineamento tra "quello che è stato scritto" e "quello che il sistema utilizza".

## Quali capability esegue

Le sette capability dell'eseguibile corrispondono direttamente agli agenti specializzati descritti nei documenti `LabNexus_CoWork_KB.md` e `LabNexus_WikiLLM_Orchestrator.md`:

| # | Capability | Agente corrispondente |
|---|---|---|
| A | Revisione documentale dopo cambio di norma | DocAgent + NormAgent |
| B | Gestione rilievi Accredia | CAPAAgent |
| C | Management Review Pack | ReviewAgent |
| D | Checklist d'audit da risk register | AuditAgent |
| E | Equipment alert per scadenze tarature | EquipmentAgent |
| F | Competence gap da matrice competenze | CompetenceAgent |
| G | Analisi z-score PT | PTQCAgent |

## Provider del modello

Per Sprint 1, tutte le capability girano su **Qwen 3 in locale via Ollama**, sull'hardware del laboratorio. È la modalità che rispetta il vincolo *locale mandatory* della domanda contrattuale di sprint.

L'eseguibile include anche un secondo canale tecnico, **EUrouter** (gateway europeo OpenAI-compatible, dati residenti in UE, GDPR-compliant), utilizzato come strumento di debug interno. Serve nei casi in cui un test su Qwen produca risultati inattesi: la stessa identica chiamata può essere lanciata su un modello più grande via EUrouter per discriminare tra problema del modello locale e problema della pipeline software.

**Vincolo di utilizzo**: EUrouter non viene utilizzato sui dati reali del SGQ di Denis senza approvazione esplicita. In Sprint 1 resta uno strumento di debug, non un canale operativo.

## Cosa non fa l'eseguibile

Per chiarezza, vengono dichiarati i confini di scope di questo Sprint 1:

- **Non osserva** la cartella SGQ in background. Niente watcher, niente trigger automatici. L'esecuzione avviene su richiesta.
- **Non orchestra** più capability in sequenza. Una esecuzione corrisponde a una capability.
- **Non modifica** i file del SGQ originale. Legge la cartella di input in sola lettura, scrive esclusivamente nella cartella di output.
- **Non sostituisce il giudizio del QM**. Ogni output è etichettato come bozza, da approvare manualmente.
- **Non ha interfaccia web**. È un eseguibile da Mac.
- **Non ha database centralizzato**. Tutto è filesystem, tutto è markdown.

Tutte queste funzionalità sono presenti nei documenti CoWork e WikiLLM come parti del sistema AICertus completo, e costituiscono la roadmap per Sprint 2 e successivi.

## Modalità d'uso nel tempo

**Durante lo Sprint 1**: l'eseguibile viene utilizzato per i sette esperimenti di validazione descritti nel piano degli esperimenti. Per ogni esperimento, l'output prodotto viene valutato in una sessione strutturata con Denis, secondo griglia condivisa. Il giudizio del QM viene formalizzato e raccolto nel report finale.

**Dopo lo Sprint 1**: l'eseguibile può continuare a essere utilizzato sull'hardware locale per le esigenze ispettive operative del QM. Per ogni nuovo caso (una procedura da revisionare, un set di rilievi da gestire, un riesame da preparare), l'utente lancia la capability appropriata sulla cartella di input. Il modello locale gira sull'hardware del laboratorio, senza connessione internet, con dati che non escono dalla macchina.

L'aggiunta di nuove capability non previste nei sette iniziali è un'attività di Sprint 2: richiede l'aggiunta di un nuovo file di configurazione, senza ricompilazione del software, ma è lavoro tecnico che si svolge sotto la responsabilità del CTO.

## Consegna a fine sprint

A fine Sprint 1 viene consegnato un singolo file zip contenente:

- L'eseguibile `labnexus.app`, da installare nella cartella Applicazioni del Mac
- I sette file di configurazione delle capability (modificabili dal CTO in Sprint 2, non oggetto di modifica diretta da parte di Denis in questa fase)
- La copia della KB-ispettore al momento della consegna; le iterazioni successive vengono fatte da Denis sulla propria versione locale, che viene letta dall'eseguibile a ogni esecuzione
- README di una pagina con: procedura di autorizzazione macOS per software non firmato, modalità di lancio, posizione degli output, gestione di eventuali errori
- Esempi di input per ciascuna capability

## Nota sull'autorizzazione macOS

Alla prima apertura, macOS bloccherà l'eseguibile perché non firmato con certificato Apple Developer (~100€/anno; l'investimento sarà valutato in Sprint 2). La procedura di autorizzazione è standard per software non firmato (control-clic sull'icona, "Apri", conferma), si esegue una sola volta, ed è documentata nel README.

## Significato del deliverable nello scope complessivo del progetto

L'eseguibile e i suoi sette profili rappresentano **il primo strato costruibile del sistema AICertus**. Ogni profilo è un task LLM-critical del sistema completo, validato individualmente. Il pattern architetturale di Sprint 1 — separazione netta tra parte deterministica (Sprint 2 e successivi) e parte LLM (Sprint 1), KB del dominio scritta dall'esperto, esecuzione locale che protegge i dati — è il fondamento su cui si costruirà l'orchestrazione completa negli sprint successivi.

I risultati dei sette esperimenti dello sprint determineranno la direzione di Sprint 2: capability che si dimostreranno solide diventeranno componenti del sistema completo, capability che si dimostreranno deboli indicheranno dove serve pre-processing deterministico aggiuntivo o cambio di modello.
