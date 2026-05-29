---
titolo: KB-Ispettore — Guida Completa alla Knowledge Base per LabNexus Agent
tipo: kb_ispettore
origine: LabNexus Knowledge Base per Reasoning Ispettivo ACCREDIA
autore: Denis Brazzo — Ispettore Tecnico ACCREDIA DL0871
data_creazione: 2026-05-06
versione: 1.0
stato: completo
---

# KB-ISPETTORE — Guida Completa

> **Questa Knowledge Base contiene l'expertise ispettivo di Denis Brazzo, Ispettore Tecnico ACCREDIA DL0871, strutturato in domande e logiche di controllo per guidare il LabNexus Agent nell'analizzare i laboratori accreditati.**

---

## STRUTTURA DELLA KB-ISPETTORE

```
/KB-ispettore/
├── MAPPA_ESPLOSA_REQUISITI_ISO17025.md (il documento madre)
├── come-pensa-un-ispettore.md (logica ispettiva)
├── come-rispondere-NC-ACCREDIA.md (struttura risposta NC)
├── errori-fatali-da-evitare.md (alert prioritari)
├── CLAUDE.md (system prompt per agente)
│
├── sezioni-ISO/
│   ├── sezione-4-requisiti-generali-domande.md (IMPARZIALITÀ, RISERVATEZZA)
│   ├── sezione-5-requisiti-strutturali-domande.md (ORGANIZZAZIONE)
│   ├── sezione-6-risorse-personale-dotazioni-domande.md (PERSONALE, DOTAZIONI, TARATURA)
│   ├── sezione-7-processo-metodi-validazione-domande.md (METODI, VALIDAZIONE, QC, RAPPORTI)
│   ├── sezione-8-sgq-domande.md (SISTEMA DI GESTIONE) [in preparazione]
│   └── [sezioni 1-3: concettuali, non richiedono domande ispettive]
│
└── [appendici: PT, incertezza, taratura, etc.]
```

---

## COSA CONTIENE LA KB

### 1. **MAPPA_ESPLOSA_REQUISITI_ISO17025.md** (Documento Fondazionale)
- **Scopo**: Elenco COMPLETO e STRUTTURATO di TUTTI i requisiti ISO/IEC 17025:2018 come implementati in RT-08 ACCREDIA
- **Organizzazione**: 8 Sezioni (1-8) con sottorequisiti per ognuna
- **Contenuto per requisito**:
  - Testo esatto del requisito ISO
  - Elementi chiave
  - Riferimenti normativi (§ ISO e § RT-08)
- **Uso**: Base per generazione di domande ispettive, mappa del perimetro accreditato

### 2. **come-pensa-un-ispettore.md**
- **Scopo**: Descrizione della LOGICA ispettiva ACCREDIA
- **Contiene**:
  - Come un ispettore legge un requisito
  - Cosa cerca nel laboratorio
  - Quali sono i "segnali" di rischio
  - Differenza tra "conforme sulla carta" vs "conforme in realtà"
- **Uso**: Priming del LabNexus Agent per ragionamento ispettivo

### 3. **come-rispondere-NC-ACCREDIA.md**
- **Scopo**: Guida COMPLETA su come un laboratorio deve rispondere a una Non-Conformità
- **Contiene**:
  - Struttura corretta di una risposta NC
  - Gli 7 step logici di una risposta efficace
  - Errori comuni (e perché falliscono)
  - Differenza tra "risposta formalmente corretta" e "risposta tecnicamente fondata"
- **Uso**: Guida il LabNexus Agent nel valutare le risposte NC dei laboratori

### 4. **errori-fatali-da-evitare.md**
- **Scopo**: Elenco di ERRORI CRITICI che generano immediatamente NC grave o sospensione
- **Contiene**:
  - 8+ errori che "non ammettono attenuanti"
  - Esempi concreti di come si manifestano
  - Perché l'ispettore reagisce severamente
  - Procedure di alert prioritario
- **Uso**: Guida il LabNexus Agent nell'identificazione di situazioni critiche

### 5. **Sezioni ISO con 3 Domande Ispettive per Requisito** (sezione-X-domande.md)
- **Scopo**: Per ogni requisito ISO, 3 domande che un ispettore pone per verificare conformità
- **Struttura di ogni domanda**:
  
  **DOMANDA 1 — Il Meccanismo di Controllo**
  - Verifica se il laboratorio HA un controllo documentato che governa quel aspetto
  - Esempio: "Quale barriera procedurale impedisce X?"
  - Analisi ispettiva: cosa l'ispettore cerca, come valuta, note critiche
  
  **DOMANDA 2 — Ricerca dell'Estensione**
  - Verifica se il laboratorio HA CERCATO il problema altrove
  - Esempio: "Dove/come avete cercato se questo problema esiste in altri contesti?"
  - Analisi ispettiva: differenza tra episodio isolato vs problema sistemico
  
  **DOMANDA 3 — L'Evidenza di Efficacia**
  - Verifica se il controllo FUNZIONA DAVVERO su casistica reale
  - Esempio: "Portate un caso concreto dove il controllo ha intercettato il problema"
  - Analisi ispettiva: differenza tra controllo teorico vs operativo

- **File attuali**:
  - ✅ sezione-4-requisiti-generali-domande.md (9 requisiti × 3 domande = 27 domande)
  - ✅ sezione-5-requisiti-strutturali-domande.md (7 requisiti × 3 domande = 21 domande)
  - ✅ sezione-6-risorse-personale-dotazioni-domande.md (sottosezioni 6.2, 6.4, 6.5 = 18 domande)
  - ✅ sezione-7-processo-metodi-validazione-domande.md (sottosezioni 7.2, 7.5, 7.7 = 24 domande)
  - ⏳ sezione-8-sgq-domande.md (in preparazione)

---

## COME USARE LA KB-ISPETTORE

### A. Analisi di un Laboratorio

**Scenario**: Il LabNexus Agent riceve un rapporto di self-assessment da un laboratorio.

1. **Carica la MAPPA_ESPLOSA** per capire i requisiti ISO che si applicano allo scopo del laboratorio
2. **Leggi "come-pensa-un-ispettore.md"** per allinearsi alla logica ispettiva
3. **Consulta la sezione ISO corrispondente** (es. sezione-6-risorse se è su personale/dotazioni)
4. **Poni le 3 domande ispettive** per ogni requisito critico
5. **Valuta le risposte** secondo la logica ispettiva
6. **Identifica rischi** e genera alert se necessario

### B. Valutazione di una Risposta NC

**Scenario**: Un laboratorio ha ricevuto una NC e devo valutare se ha risposto adeguatamente.

1. **Consulta "come-rispondere-NC-ACCREDIA.md"** per capire la struttura attesa
2. **Leggi "errori-fatali-da-evitare.md"** per identificare se la risposta contiene errori critici
3. **Applica i 7 step logici** per valutare se la causa è stata realmente compresa
4. **Genera alert** se mancano: causa reale, estensione, impatto, efficacia

### C. Generazione di Suggerimenti per il Laboratorio

**Scenario**: Il LabNexus Agent suggerisce al laboratorio come migliorare la conformità.

1. **Consulta la sezione ISO rilevante** per le 3 domande ispettive
2. **Identifica quale domanda il laboratorio non sa rispondere bene**
3. **Fornisci suggerimenti basati su "come-pensa-un-ispettore.md"**
4. **Contesta non con norme astratte, ma con logica ispettiva concreta**

---

## LOGICA ISPETTIVA FONDAMENTALE

### I 3 Livelli di Controllo (presente in ogni domanda)

| Livello | Domanda | Cosa Cerca | Rischio |
|---------|---------|-----------|--------|
| **Meccanismo** | Avete il controllo? | Documentazione, procedura, barriera | Controllo assente o informale |
| **Estensione** | Avete cercato il problema altrove? | Ricerca sistematica, non episodio isolato | Problema sistemico non visto |
| **Efficacia** | Funziona davvero? | Evidenza su casistica reale, non teorica | Controllo inefficace in pratica |

> Se manca anche uno di questi tre livelli, la risposta del laboratorio è **INCOMPLETA**.

---

## ERRORI COMUNI CHE LA KB AIUTA A IDENTIFICARE

### ✗ **ERRORE 1: La causa è un effetto** (da come-rispondere-NC-ACCREDIA.md)
- Laboratorio scrive: "La causa è l'errore umano dell'operatore"
- Ispettore sa: L'operatore è un sintomo, la causa vera è quale BARRIERA NON HA FUNZIONATO
- **Come la KB aiuta**: Le 3 domande ispettive forzano il laboratorio a pensare ai controlli di sistema

### ✗ **ERRORE 2: Non valutare l'estensione**
- Laboratorio scrive: "Abbiamo corretto il caso campionato"
- Ispettore sa: Un campione è un segnale, serve ricerca dove ricorre
- **Come la KB aiuta**: DOMANDA 2 esplicita chiede "Avete cercato altrove?"

### ✗ **ERRORE 3: Omettere l'impatto sui risultati passati**
- Laboratorio scrive: "D'ora in poi applicheremo la procedura correttamente"
- Ispettore sa: La domanda prioritaria è "Quali risultati del passato potrebbero essere compromessi?"
- **Come la KB aiuta**: DOMANDA 3 chiede evidenza che i dati passati sono stati protetti

### ✗ **ERRORE 4: Dichiarare efficacia prima di verificare**
- Laboratorio scrive: "La procedura è stata aggiornata, quindi efficace"
- Ispettore sa: Attuazione ≠ Efficacia. Serve verifica su casistica reale
- **Come la KB aiuta**: DOMANDA 3 richiede "Portate un caso concreto dove ha funzionato"

---

## UTILIZZO CON IL LABNEXUS AGENT

### System Prompt Predisposto
Il file **CLAUDE.md** contiene il system prompt già allineato alla KB-ispettore:
- Istruzioni per applicare la logica ispettiva
- Riferimenti alle domande delle sezioni ISO
- Workflow di analisi
- Alert prioritari da "errori-fatali-da-evitare.md"

### Workflow Suggerito
```
1. LabNexus Agent riceve dati laboratorio
   ↓
2. Consulta MAPPA_ESPLOSA per capire requisiti applicabili
   ↓
3. Legge come-pensa-un-ispettore.md per allineamento logico
   ↓
4. Seleziona sezione ISO rilevante
   ↓
5. Pone le 3 DOMANDE ISPETTIVE per requisiti critici
   ↓
6. Valuta risposte secondo la logica (meccanismo → estensione → efficacia)
   ↓
7. Identifica rischi e genera output
   └── Suggerimenti per laboratorio
   └── Alert per ispettore
   └── Valutazione di conformità
```

---

## PROSSIMI PASSI DI SVILUPPO

### ✅ **COMPLETATI**
- Mappa Esplosa (sezioni 1-8)
- come-pensa-un-ispettore.md
- come-rispondere-NC-ACCREDIA.md
- errori-fatali-da-evitare.md
- Sezione 4 — 3 domande per requisito
- Sezione 5 — 3 domande per requisito
- Sezione 6 (sottosezioni) — 3 domande per requisito
- Sezione 7 (sottosezioni) — 3 domande per requisito

### ⏳ **IN PREPARAZIONE**
- Sezione 8 — 3 domande per requisito (SGQ, audit interni, azioni correttive)
- Sezioni 1-3 — Analisi concettuale (non richiedono domande, ma guida interpretativa)

### 🎯 **INTEGRAZIONI FUTURE**
- FAQ ispettivo (risposte a domande ricorrenti)
- Case study di NC reali e come risolverle
- Appendici tecniche (PT, incertezza, taratura, CRM)
- Glossario ISO-ACCREDIA
- Template di documenti conformi
- Integrazione nel LabNexus Agent come knowledge engine

---

## NOTE IMPORTANTI PER L'USO

### ⚠️ Non sostituisce la visita ispettiva
La KB-ispettore è uno **strumento di auto-assessment e analisi documentale**. Non sostituisce:
- La visita ispettiva fisica
- L'osservazione diretta dei processi
- Le interviste con il personale
- Le verifiche di laboratorio reali

### ⚠️ Personalizzazione per laboratorio
Le domande sono **strutturate per laboratori di analisi generici**. Potrebbero necessitare adattamenti per:
- Laboratori di taratura (domande su metrologiae riferibilità enfatizzate)
- Laboratori specializzati (biologici, microbiologici, etc.)
- Laboratori con attività molto diverse

### ⚠️ Aggiornamenti
- Questa KB è versione 1.0 (2026-05-06)
- Sarà aggiornata quando nuovi criteri ACCREDIA o nuove norme ISO entrano in vigore
- Incorporerà feedback da casi di visita reali

---

## CONTATTI E FEEDBACK

**Denis Brazzo**  
Ispettore Tecnico ACCREDIA DL0871  
Chimico Albo Nazionale n. 935 sez. A  
DASP Srl  
info@labnexus.app

Per integrazioni, feedback, o applicazioni specifiche della KB-ispettore, contattare Denis direttamente.

---

*KB-Ispettore Versione 1.0*  
*Creata: 2026-05-06*  
*Stato: Completa e Operativa*  
*Ultima revisione: 2026-05-06*
