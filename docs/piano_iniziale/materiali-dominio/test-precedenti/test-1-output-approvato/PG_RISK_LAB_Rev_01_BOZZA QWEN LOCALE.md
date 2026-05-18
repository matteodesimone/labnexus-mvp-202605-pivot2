---
codice: PRO-GES-RIS-00
titolo: Analisi e Gestione del Rischio in Laboratorio
tipo: procedura
origine: interna
revisione: "01"
data_emissione: 2024-05-20
data_prossima_revisione: 2027-05-20
stato: bozza
sezione_iso: ["6.1", "8.1", "8.5"]
norme_riferimento:
  - ISO-9001-2015-6.1
  - ISO-31000-2018
  - RT-08-REV05
  - ISO-IEC-17025-2018
documenti_collegati:
  - MOD_SWOT_01
  - MOD_PDM_01
  - MOD_TDL_01
  - come-rispondere-NC-ACCREDIA.md
  - errori-fatali-da-evitare.md
scadenze:
  - titolo: Riesame annuale piano rischi
    frequenza_giorni: 365
    data_prossima: 2025-05-20
    sezione_iso: "8.9"
    output_template: MOD_SWOT_01
ultima_modifica_trigger: esterno_accredia
ultima_modifica_nota: Allineamento a RT-08 rev.05 (EC 10-02-2022), ISO 9001:2015 §6.1 e ISO 31000:2018. Ricalcolo criteri matrice, esplicitazione processo risk management, linkage a NC/AC, PT risk-based e validazione metodi. Chiarimento responsabilità e trigger di aggiornamento.
approvato_da: [IN ATTESA VERIFICA QM/CERTIFICATORE]
approvato_at: 
tags: [gestione-rischio, RT-08-rev05, ISO-31000, risk-based, SGQ]
---

# Sintesi per il QM — Revisione Procedura di Gestione del Rischio

> [!ATTENZIONE] **Lettura sintetica del rischio complessivo**: La procedura originale era formalmente completa ma tecnicamente generica. Mancava l'allineamento esplicito a ISO 9001:2015 §6.1 e ISO 31000:2018, non collegava la matrice dei rischi ai criteri di RT-08 rev.05 (PT risk-based, validazione metodi, imparzialità), e trattava SWOT/FMEA come metodo invece che come strumenti. Il rischio di NC è concreto su audit e ispezioni se non si dimostrano i meccanismi di controllo, l'estensione e l'efficacia del processo di risk management.
> 
> **Priorità 1** — Bloccare o valutare prima dell'uso/emissione: nessuna attività di PT o validazione metodi può essere pianificata o approvata senza riferimento al piano rischi aggiornato.
> **Priorità 2** — Chiudere con evidenza entro una data: aggiornare tutti i moduli collegati (MOD_SWOT_01, MOD_PDM_01, MOD_TDL_01) e comunicare agli operatori prima dell'entrata in vigore.
> **Priorità 3** — Pianificare e monitorare: inserire il piano rischi nel riesame della direzione (§8.9) e negli audit interni (§8.8).

## Revisione Documento

-------------------------------------
| Rev | Data emissione | Motivo/Modifica |
-------------------------------------
| 01  | 20/05/2024     | Allineamento a RT-08 rev.05 (EC 10-02-2022), ISO 9001:2015 §6.1 e ISO 31000:2018. Ricalcolo criteri matrice, esplicitazione processo risk management, linkage a NC/AC, PT e validazione metodi. Aggiornamento riferimenti normativi e responsabilità. |
---------------------------

# 1. SCOPO
> [!MODIFICA] Aggiornato allineamento normativo e chiarimento obiettivo operativo.
Lo scopo di questa procedura è descrivere il processo sistematico di gestione del rischio, conforme a **ISO 9001:2015 §6.1** e **ISO 31000:2018**, al fine di:
* garantire il raggiungimento degli obiettivi del SGQ e la conformità ai requisiti di **UNI CEI EN ISO/IEC 17025:2018**;
* identificare, analizzare e trattare proattivamente i rischi e le opportunità, minimizzando gli impatti sulle attività di laboratorio;
* fornire al Riesame della Direzione (§8.9) e all'Audit Interno (§8.8) evidenze misurabili di efficacia del sistema;
* supportare l'approccio risk-based richiesto da **RT-08 rev.05** per la pianificazione dei PT (RT-24), la validazione dei metodi (§7.2) e la gestione della validità dei risultati (§7.7).

La gestione del rischio è incorporata in tutte le attività operative e gestionali. L'efficacia è verificata periodicamente e documentata in **MOD_SWOT_01** e riportata al §6 della presente procedura.

# 2. TERMINI E DEFINIZIONI
> [!MODIFICA] Allineato a ISO 31000:2018 e ISO/IEC 17025:2018. Rimossi termini generici non normativi.
* **Rischio**: effetto dell'incertezza sugli obiettivi (ISO 31000:2018, 3.6.1).
* **Gestione del rischio**: processo coordinato per dirigere e controllare un'organizzazione rispetto al rischio (ISO 31000:2018, 3.6.2).
* **Processo di gestione del rischio**: applicazione sistematica di politiche, procedure e prassi per comunicare, consultare, definire il contesto, identificare, analizzare, valutare, trattare, monitorare e riesaminare il rischio (ISO 31000:2018, 4.1).
* **Fonte di rischio**: elemento che, da solo o in combinazione, possiede il potenziale intrinseco di originare il rischio.
* **Evento**: verificarsi o modificarsi di un particolare insieme di circostanze.
* **Conseguenza**: esito di un evento che influenza gli obiettivi.
* **Livello di rischio**: combinazione della probabilità di accadimento e della gravità delle conseguenze.
* **Trattamento del rischio**: processo per modificare il rischio (es. evitare, mitigare, trasferire, accettare).
* **Controllo**: misura che tiene sotto controllo il rischio.
* **Rischio residuo**: rischio rimanente dopo il trattamento.
* **Opportunità**: evento con potenziali conseguenze positive.
* **Imparzialità**: assenza di pregiudizi o conflitti d'interesse che possano influire sulle attività di laboratorio (ISO/IEC 17025:2018, 4.1).
* **Validità dei risultati**: grado in cui i dati ottenuti sono affidabili e possono essere utilizzati per lo scopo previsto (ISO/IEC 17025:2018, 7.7).
* **Pianificazione risk-based**: approccio che dirige le risorse di controllo verso le attività a più alto rischio, ottimizzando l'efficacia del SGQ.

# 3. RIFERIMENTI NORMATIVI
> [!MODIFICA] Aggiornati e allineati a RT-08 rev.05 e standard internazionali di riferimento.
* UNI CEI EN ISO/IEC 17025:2018 — Requisiti generali per la competenza dei laboratori di prova e taratura
* ISO 9001:2015 — Sistemi di gestione per la qualità. Requisiti (§6.1 Azioni per affrontare rischi e opportunità)
* ISO 31000:2018 — Gestione del rischio. Principi e linee guida
* RT-08 rev.05 (EC 10-02-2022) — Prescrizioni per l'accreditamento dei laboratori di prova
* RT-24 — Partecipazione a prove valutative interlaboratorio (PT)
* RG-09 — Uso del Marchio di accreditamento e riferimento all'accreditamento
* MOD_ELE_DOC_09 — Elenco documenti di riferimento del laboratorio

# 4. CAMPO DI APPLICAZIONE
> [!MODIFICA] Chiariti i criteri di selezione, aggiornata la matrice e specificati i trigger di revisione.
La gestione del rischio è applicata a tutti i processi del SGQ che possono impattare sulla competenza tecnica, sull'imparzialità, sulla riservatezza e sulla validità dei risultati. Il laboratorio sottopone a analisi del rischio:
* Le attività esplicitamente richieste da **ISO/IEC 17025:2018** e **RT-08 rev.05** (es. validazione metodi, pianificazione PT, gestione dotazioni, controllo risultati);
* I processi gestionali che influenzano direttamente la capacità del laboratorio di produrre risultati validi e difendibili;
* I cambiamenti organizzativi, tecnologici, normativi o di personale.

### 4.1 Criteri adottati per la definizione del campo di applicazione
L'identificazione dei rischi e delle opportunità viene aggiornata:
* Annualmente, in occasione del riesame della direzione (§8.9);
* In seguito a cambiamenti rilevanti (nuovi metodi, dotazioni, sedi, subappalti, variazioni normative);
* In seguito a NC, reclami, PT negativi o incidenti tecnici significativi.

La frequenza e l'impatto dei rischi sono valutati secondo la matrice sottostante. I valori sono ponderati in base al contesto operativo del laboratorio e all'obiettivo (tecnico o gestionale).

#### Tabella dell'Impatto (Gravità)
> [!MODIFICA] Criteri resi operativi e allineati a RT-08 rev.05 §7.10 (NC) e §7.7 (validità risultati).
| Valore | Scala | Descrizione |
|------|-------|-----------|
| 1-2  | Basso | Errore isolato, correggibile senza impatto sul risultato o sul cliente. Nessuna violazione normativa. |
| 3-4  | Moderato | Errore ripetitivo o sistematico. Perdita di efficienza, costo aggiuntivo, potenziale reclamo o NC minore. Impatto limitato sulla validità del risultato. |
| 5    | Alto  | Errore grave con impatto diretto sulla validità del risultato, violazione di requisiti cogenti/accidentali, pericolo per la salute/sicurezza, sospensione accreditamento o comunicazione al cliente. |

#### Tabella della Frequenza (Probabilità)
> [!MODIFICA] Allineata a ISO 31000:2018.
| Valore | Scala | Descrizione |
|------|-------|-----------|
| 1-2  | Molto raro | Controllo efficace, risorse adeguate, storico impeccabile. Probabilità < 5%. |
| 3-4  | Raro | Controllo parziale o dipendente da intervento manuale. Probabilità 5-15%. |
| 5    | Frequente | Processo complesso, risorse limitate, controlli insufficienti o assenti. Probabilità > 15%. |

#### Matrice di Valutazione del Rischio
| Impatto ↓ / Frequenza → | Molto raro (1-2) | Raro (3-4) | Frequente (5) |
|----------------|------------|--------|-------|
| **Basso (1-2)** | Basso | Basso | Medio |
| **Moderato (3-4)** | Basso | Medio | Alto |
| **Alto (5)** | Medio | Alto | Alto |

> [!ATTENZIONE] Il trattamento dei rischi con livello **ALTO** è prioritario. Richiede azioni correttive sistemiche, monitoraggio intensivo e verifica di efficacia documentata. I rischi **MEDIO** richiedono controlli specifici e revisione periodica. I rischi **BASSO** sono accettati con monitoraggio ordinario, purché non cumulativi.

### 4.2 Processi coinvolti
I rischi sono mappati per macro-processi. Ogni rischio identificato deve essere associato al processo di competenza e alla funzione responsabile:
* **M1-GESTIONE DEL CLIENTE**: Riesame contratti, selezione metodi, collaborazione cliente. Rischio principale: incomprensioni contrattuali, metodi non validati per la matrice richiesta.
* **M2-PROCESSO ANALITICO**: Campionamento, manipolazione, esecuzione metodi, registrazioni, refertazione. Rischio principale: errore operativo, perdita tracciabilità, validità risultato compromessa.
* **M3-RISORSE UMANE**: Competenza, autorizzazioni, imparzialità, riservatezza. Rischio principale: personale non autorizzato, conflitti d'interesse, mancanza di formazione.
* **M4-APPARECCHIATURE**: Disponibilità, manutenzione, taratura, riferibilità metrologica, CRM. Rischio principale: dotazione fuori stato, catena di riferibilità interrotta, criteri di accettazione non definiti.
* **M5-APPROVVIGIONAMENTO/INFRASTRUTTURA**: Locali, condizioni ambientali, acquisti, flussi. Rischio principale: condizioni non controllate, fornitori non valutati, interruzione servizi critici.
* **M6-VERIFICA E VALIDAZIONE METODI**: Incertezza, validità risultati, validazione/verifica metodi. Rischio principale: metodo non verificato/validato, incertezza non linearità, piano PT inadeguato.
* **M7-SISTEMA DI GESTIONE**: Documentazione, controllo dati, registrazioni (conservazione ≥48 mesi). Rischio principale: documenti obsoleti, dati non ricostruibili, conservazione inadeguata.
* **M8-MIGLIORAMENTO E PIANIFICAZIONE**: Reclami, NC, azioni correttive, audit, direzione, miglioramento. Rischio principale: NC ricorrenti, AC inefficaci, miglioramento non monitorato.

Il laboratorio valuta esplicitamente il rischio di **imparzialità** (§4.1 ISO 17025) e di **confidentialità** (§4.2) in tutti i processi sopra elencati. La valutazione è continua e documentata.

# 5. RESPONSABILITÀ
> [!MODIFICA] Chiariti i ruoli e allineati a §5.6 e §6.2 ISO 17025.
* **Top Management / Direzione**: assicura le risorse, approva la politica di gestione del rischio e il piano rischi annuale, e valuta l'efficacia nel riesame della direzione (§8.9).
* **Responsabile Qualità (RL / RTQL)**: coordina il processo, mantiene il registro dei rischi, verifica l'allineamento a RT-08 rev.05 e ISO 17025, garantisce la formazione del personale e l'efficacia dei controlli.
* **Responsabile Tecnico (RT)**: valuta l'impatto tecnico dei rischi (validità risultato, criteri metrologici, limiti di metodo), approva il trattamento dei rischi tecnici e verifica l'efficacia sulle attività reali.
* **Responsabili di Macro-Processo (M1-M8)**: identificano i rischi nei propri ambiti, applicano i controlli, segnalano variazioni o eventi anomali, aggiornano i moduli collegati.
* **Personale operativo**: applica i controlli definiti, segnala potenziali rischi, partecipa alla formazione e al riesame periodico.

Tutti i membri del laboratorio sono formati sulle procedure di gestione del rischio e sulla propria responsabilità nel presidio dei controlli.

# 6. GESTIONE DEL RISCHIO
> [!MODIFICA] Ricalibrato su ISO 31000:2018. SWOT e FMEA sono ora strumenti, non il metodo. Chiarito il flusso completo e il linkage a NC/AC/PT.

### 6.1 Generalità
La gestione del rischio è un processo continuo e iterativo, non un'attività occasionale. Il laboratorio adotta il ciclo di gestione del rischio conforme a **ISO 31000:2018**:
1. **Definizione del contesto** (interno/esterno e obiettivi);
2. **Identificazione dei rischi e delle opportunità**;
3. **Analisi del rischio** (probabilità × impatto, fonti, cause);
4. **Valutazione del rischio** ( confronto con criteri, priorità);
5. **Trattamento del rischio** (evitare, mitigare, trasferire, accettare);
6. **Monitoraggio e riesame** (verifica efficacia, aggiornamento).

Il processo alimenta il miglioramento continuo, la gestione delle non conformità (§7.10), le azioni correttive (§8.7) e la pianificazione risk-based dei PT (§7.7 / RT-24).

### 6.2 Tecniche operative
Il laboratorio utilizza strumenti strutturati per supportare le fasi del processo. Tali strumenti non sono il metodo in sé, ma supportano l'identificazione e l'analisi:
* **SWOT Analysis**: utile per mappare Punti di Forza/Debolezze (contesto interno) e Opportunità/Minacce (contesto esterno). Supporta la pianificazione strategica e il riesame della direzione.
* **FMEA di processo (Failure Mode and Effects Analysis)**: applicata ai processi critici (es. validazione metodi, gestione dotazioni, refertazione). Identifica modi di guasto, effetti e cause prime. Assegna un RPN (Risk Priority Number) o equivalente per prioritizzare il trattamento.
* **Diagramma di Ishikawa (Causa-Effetto)**: utilizzato in fase di analisi di causa per NC, fuori controllo QC o eventi tecnici.
* **Classificazione metodi e dotazioni**: i metodi di prova sono classificati per classe di rischio in base a complessità, criticità del misurando, frequenza di errore e impatto sul cliente. La classe determina la frequenza di CQ, la robustezza della validazione e la severità del controllo.

> [!ATTENZIONE] L'uso di questi strumenti deve essere documentato con evidenze oggettive (verbali, moduli compilati, report). Non sono ammesse analisi "a mente" o non registrate.

### 6.3 Tecnica adottata dal laboratorio
Il laboratorio adotta il processo descritto al §6.1, supportato dagli strumenti al §6.2. Il flusso operativo è:
1. **Identificazione**: RL/RTQL, con le funzioni coinvolte, mappa i rischi sui macro-processi M1-M8. Si utilizzano dati storici, audit, PT, reclami, variazioni normative e valutazioni tecniche.
2. **Analisi e Valutazione**: i rischi sono classificati nella matrice §4.1. Si valuta l'estensione (ricorrenza, altri processi/dotazioni/metodi impattati) e l'impatto tecnico/gestionale.
3. **Trattamento**: per ogni rischio, si definisce un'azione specifica. Le opzioni sono:
   * *Mitigazione*: introduzione/modifica controllo, formazione, segregazione, miglioramento procedura.
   * *Accettazione*: solo se il rischio è basso o il costo di mitigazione sproporzionato, previa approvazione RL/RT. Documentazione esplicita del rischio residuo.
   * *Trasferimento/Evitamento*: dove applicabile (es. subappalto a laboratorio accreditato, rinuncia a prova non valida).
4. **Pianificazione e Azione**: le azioni sono inserite in **MOD_PDM_01** (miglioramento proattivo) o **MOD_TDL_01** (intervento reattivo). Sono assegnati responsabili, tempi e criteri di efficacia.
5. **Monitoraggio e Riesame**: l'efficacia è verificata su casistica reale (nuovi dati, rapporti, registrazioni). I risultati sono riportati in **MOD_SWOT_01** e presentati al Riesame della Direzione (§8.9). La matrice dei rischi è rivista annualmente o in seguito a eventi significativi.

> [!ISPETTORE] In visita ACCREDIA, l'ispettore cercherà: (1) evidenze che il piano rischi ha diretto le risorse verso le attività critiche, (2) correlazione tra livello di rischio e frequenza di PT/CQ, (3) traccia di trattamento dei rischi di imparzialità e validità risultato, (4) prove che le AC derivano dai rischi residui e che la loro efficacia è verificata. Se il piano rischi è solo un elenco statico senza linkage a PT, validazione o NC, il rilievo è concreto.

## Sintesi per approvazione
> [!MODIFICA] **Sintesi per approvazione**: modifiche sostanziali in 4 aree; nessun impatto retroattivo se il piano rischi è già attivo e monitorato; richiesta verifica del Responsabile Tecnico per matrice criteri e linkage a PT risk-based/validazione metodi; necessaria comunicazione agli operatori e aggiornamento dei moduli MOD_SWOT_01, MOD_PDM_01, MOD_TDL_01 prima dell'entrata in vigore. Firmare questo documento prima del rilascio in versione vigente.

## Checklist di Verifica Umana (Certificatore / QM)
> [!ISPETTORE] Il certificatore deve validare ogni modifica percorrendo i tre livelli di controllo. Inserire esito in colonna "Verifica Umana" prima dell'approvazione.

| Articolazione Modifica | Meccanismo (barriera/documento?) | Estensione (problema cercato altrove?) | Efficacia (evidenza su casistica reale?) | Verifica Umana [✓/✗] | Note/Interventi |
|---|---|---|---|---|---|
| Matrice criteri risk | Criteri espliciti e allineati a RT-08/ISO 31000 | Valutazione estesa a tutti i macro-processi M1-M8 e a dotazioni/metodi | Piano rischi utilizzato per prioritizzare PT e validazione metodi | | |
| Processi coinvolti (M1-M8) | Mappatura completa e trigger di aggiornamento | Rischi di imparzialità, riservatezza, validità risultati inclusi | Evidenza di monitoraggio continuo e linkage a NC/AC | | |
| Strumenti (SWOT/FMEA) | Documentati come supporto, non metodo | Applicati ai processi critici identificati | Report/verbali con evidenze di analisi e trattamento | | |
| Linkage a PT/Validazione | Approccio risk-based esplicito | Metodi/dotazioni ad alto rischio coperti da PT/CQ più frequenti | Piano PT e validazione metodi allineati alla matrice dei rischi | | |
| Riesame e miglioramento | Input al §8.9 e §8.8 | Trend di rischio monitorato e azioni verificate | Verbale riesame direzione con evidenze di efficacia | | |

> [!ATTENZIONE] Prima della firma, verificare che: (1) i moduli collegati siano aggiornati, (2) la formazione sul nuovo processo sia documentata, (3) il piano rischi sia integrato nel piano di audit interno (§8.8) e nel piano PT (RT-24), (4) nessun rischio residuo ad alto livello rimanga senza trattamento documentato.

---
*Documento generato da LabNexus Agent — revisione tecnica assistita. Non sostituisce la validazione del Responsabile Qualità, del Responsabile Tecnico o del Certificatore. L'output è bozza soggetta ad approvazione formale prima del passaggio a stato vigente.*
*La struttura e il contenuto sono allineati a RT-08 rev.05 (EC 10-02-2022), ISO/IEC 17025:2018 e ISO 31000:2018. Ogni modifica è tracciata e verificabile su evidenze reali.*