---
codice: sez_8_7_azioni_correttive
tipo: requisito_normativo
livello: misto
titolo: "§ 8.7 — Azioni correttive (Opzione A)"
sezione_norma: "8.7"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "8.7" }
  - { doc: "RT-08", rev: "05", paragrafo: "8.7" }
fonti_secondarie:
  - { doc: "LS-04", rev: "20" }
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[02_TRE_LIVELLI_DI_CONTROLLO]]"
  - "[[03_RISPOSTA_NC]]"
  - "[[sez_8_5_rischi_opportunita]]"
  - "[[sez_8_6_miglioramento]]"
  - "[[sez_8_8_audit_interni]]"
tags: [azioni_correttive, NC, analisi_cause, attuazione, efficacia, root_cause]
escalation_qm: true
---

# § 8.7 — Azioni correttive (Opzione A)

> Il § 8.7 disciplina come il laboratorio reagisce alle non conformità: contenimento, analisi cause radici, definizione e attuazione delle azioni correttive, **verifica dell'efficacia**. **La distinzione fra "attuazione" ed "efficacia" è il punto critico** scrutinato in ogni ispezione Accredia.

> ⚠️ **Disclaimer di ricostruzione KB**
> Questo file fa parte della ricostruzione della **sezione 8** sull'indice ufficiale **ISO/IEC 17025:2018**. La vecchia KB aveva numerazione errata mutuata da ISO 9001. Per la **gestione del lavoro non conforme** vedi § 7.10 (file separato). I **reclami** sono al § 7.9 (file separato), NON in § 8.

---

## A. Testo della norma (🟢 LIVELLO 1 — fatto normativo)

> **Fonte**: `[ISO-17025-2018 § 8.7]` — testo a pagamento, non riprodotto integralmente.

**Sintesi operativa per sotto-paragrafo**:

- **8.7.1** — Quando si verifica una NC il laboratorio deve:
  - a) reagire alla NC e, ove applicabile: 1) prendere azioni per controllarla e correggerla; 2) gestire le conseguenze;
  - b) valutare la necessità di azioni per eliminare le cause della NC, affinché non si ripeta o non si presenti altrove, attraverso: 1) il riesame ed analisi della NC; 2) la determinazione delle cause; 3) la determinazione se NC simili esistono o potrebbero verificarsi;
  - c) attuare ogni azione necessaria;
  - d) **riesaminare l'efficacia di ogni azione correttiva intrapresa**;
  - e) aggiornare i rischi e le opportunità determinati durante la pianificazione, se necessario;
  - f) apportare modifiche al SGQ se necessario.
- **8.7.2** — Le azioni correttive devono essere **appropriate agli effetti delle NC riscontrate**.
- **8.7.3** — Il laboratorio deve conservare registrazioni come evidenza di:
  - a) natura delle NC, cause, azioni successive;
  - b) risultati di ogni azione correttiva.

---

## B. Prescrizioni Accredia (🟢 LIVELLO 2 — prescrizione Accredia)

> **Fonte**: `[RT-08 rev.05 § 8.7]`

### 8.7.1 (RT-08)
> *"Si applica il requisito di norma. Si rammenta che la verifica dell'efficacia delle azioni correttive è diversa dalla verifica della loro attuazione. Si rammenta che anche le azioni correttive pianificate e comunicate a seguito delle verifiche di seconda e terza parte (es. ACCREDIA) devono essere gestite nell'ambito del sistema di gestione del Laboratorio."*

### 8.7.2-8.7.3 (RT-08)
> *"Si applica il requisito di norma."*

---

### 🔑🔑🔑 ATTUAZIONE vs EFFICACIA — distinzione FONDAMENTALE

> **TESTO RT-08 (citazione fedele)**: *"Si rammenta che la verifica dell'efficacia delle azioni correttive è diversa dalla verifica della loro attuazione."*

Questa è **la più ricorrente NC sistemica** nei laboratori italiani. Definizioni:

| **ATTUAZIONE** (do) | **EFFICACIA** (check) |
|---|---|
| Verifica che l'azione sia stata fatta | Verifica che l'azione abbia eliminato la causa della NC |
| Esempio: "ho fatto il corso di formazione" | Esempio: "dopo il corso, il personale commette zero errori dello stesso tipo da 6 mesi" |
| Si verifica subito dopo la conclusione | Si verifica DOPO un tempo congruo (in genere 3-12 mesi) |
| Output: evidenza documentale dell'attività | Output: dato/KPI che dimostra l'eliminazione/riduzione del problema |
| Necessaria ma NON sufficiente | Obbligatoria per chiudere la NC |

**Esempio pratico — taratura saltata**:
- NC: bilancia non tarata nell'anno
- Causa radice: mancato controllo del piano metrologico
- Azione: introduzione di un calendario automatico con alert in LIMS
- **Attuazione**: alert in LIMS configurato e attivo (verificato 1 settimana dopo)
- **Efficacia**: dopo 12 mesi, 100% delle tarature programmate eseguite entro la data (verificato con report LIMS)

> Un laboratorio che chiude la NC subito dopo l'attuazione, senza verificare l'efficacia, **NON è conforme al § 8.7.1.d**.

---

### 🔑 Tempo congruo per verifica efficacia

Non esiste una regola fissa, ma:
- **NC tecnica isolata** (es. errore singolo): efficacia verificata dopo 1-3 cicli operativi simili
- **NC sistemica** (es. processo): efficacia verificata dopo 3-6 mesi
- **NC formativa** (es. errore del personale): efficacia verificata dopo 6-12 mesi
- **NC su strumenti**: efficacia verificata dopo 1 ciclo di taratura/manutenzione

---

### 🔑 Le azioni correttive da audit Accredia (seconda/terza parte)

> **TESTO RT-08**: *"Si rammenta che anche le azioni correttive pianificate e comunicate a seguito delle verifiche di seconda e terza parte (es. ACCREDIA) devono essere gestite nell'ambito del sistema di gestione del Laboratorio."*

**Implicazioni**:
- Le AC concordate con Accredia post-audit **non sono "separate"** dal sistema SGQ.
- Devono essere registrate nel **registro NC/AC interno** del laboratorio.
- Devono seguire lo **stesso processo** di analisi cause / attuazione / efficacia.
- L'efficacia deve essere verificata anche in autonomia, non solo aspettare la prossima visita Accredia.
- Frequente NC in ispezione: il laboratorio gestisce le AC Accredia su un Excel separato, ignorando il proprio sistema.

---

### 🔑 Metodologie di analisi cause radici (root cause analysis)

Strumenti accettati:
- **5 Whys** (chiedere "perché?" 5 volte)
- **Diagramma di Ishikawa (Fishbone)** — cause-effetto su 6M (Man, Machine, Method, Material, Measurement, Environment)
- **FTA** (Fault Tree Analysis) — analisi ad albero
- **8D** (8 Discipline) — metodologia automotive applicabile in lab
- **A3 problem solving** — Toyota
- **Pareto** — per individuare cause più frequenti

> Non serve un metodo "ufficiale", serve **evidenza che si è andati oltre la causa apparente**.

---

### 🔑 Estensione a NC potenziali (§ 8.7.1.b.3)

Quando una NC viene identificata in un'area, il laboratorio deve chiedersi: **NC simili esistono o potrebbero verificarsi altrove?**

Esempio:
- NC: taratura saltata bilancia laboratorio A
- Estensione: verificare TUTTE le bilance in TUTTI i laboratori
- Risultato: se trova 3 altre bilance scadute → NC sistemica, AC deve coprire tutte

---

### 🔑 Aggiornamento rischi (§ 8.7.1.e)

L'analisi della NC può rivelare nuovi rischi o richiedere rivalutazione di rischi noti.
**Collegamento esplicito § 8.5**: ogni AC chiusa con efficacia dovrebbe portare a rivalutare il rischio associato (probabilità ridotta, impatto ridotto, ecc.).

---

## C. Documenti applicabili (LS-04 rev.20)

- **RG-02** — Regolamento generale accreditamento
- **RT-08 rev.05** — Prescrizioni Accredia
- **Procedura interna gestione NC e AC** (PG-XX)
- **Registro NC/AC** (cartaceo, Excel, LIMS, software dedicato)
- **Verbali riesame direzione** (NC e AC sono input obbligatori § 8.9.2)

---

## D. Domande ispettive (🟡 LIVELLO 3 — prassi ispettiva)

> Vedi anche `[[03_RISPOSTA_NC]]` per la metodologia ispettiva sulle NC.

### 8.7.1 — Reazione e analisi causa

#### Domanda 1 — Meccanismo
*"Mostratemi la procedura di gestione NC/AC. Chi rileva? Chi analizza la causa? Chi approva l'AC? Chi verifica attuazione ed efficacia?"*

#### Domanda 2 — Estensione
*"Mostratemi il registro NC/AC degli ultimi 24 mesi. Quante NC? Quante AC? Quante chiuse? Quante con efficacia verificata?"*
- Trabocchetto: se il numero NC è troppo basso (es. <5/anno per lab medio), il laboratorio probabilmente non rileva.
- Trabocchetto inverso: se il numero NC è altissimo (>50/anno) e ricorrenti, le AC non sono efficaci.

#### Domanda 3 — Efficacia
*"Prendete una NC chiusa 12 mesi fa. Mostratemi: (a) descrizione NC, (b) analisi cause radici con metodo, (c) AC pianificata, (d) evidenza di attuazione, (e) evidenza di efficacia (KPI/dato/trend), (f) eventuale aggiornamento rischi."*

### 8.7.1.d — Verifica efficacia (DOMANDA CRITICA)

#### Domanda 1 — Meccanismo
*"Come distinguete fra verifica di attuazione e verifica di efficacia? Procedura?"*

#### Domanda 2 — Estensione
*"Quante AC chiuse l'anno scorso hanno avuto verifica di efficacia documentata? Mostratemi i casi."*

#### Domanda 3 — Efficacia
*"Avete avuto NC ripetute negli ultimi 12 mesi? (Se sì → AC precedenti non erano efficaci → ulteriore NC § 8.7.1.d)."*

### 8.7.1 — AC da audit esterno (Accredia, seconda parte)

#### Domanda 1 — Meccanismo
*"Le AC concordate dopo la scorsa visita Accredia, come sono gestite? Dove sono registrate?"*

#### Domanda 2 — Estensione
*"Mostratemi tutte le AC della scorsa ispezione Accredia: stato attuazione, stato efficacia."*

#### Domanda 3 — Efficacia
*"Per ognuna, quale evidenza di efficacia? Quando l'avete verificata in autonomia (senza aspettare la prossima Accredia)?"*

### 8.7.1.b.3 — Estensione a NC potenziali

#### Domanda 1 — Meccanismo
*"Quando rilevate una NC, come valutate se possa esistere/ripetersi altrove?"*

#### Domanda 2 — Estensione
*"Mostratemi un caso recente in cui l'analisi ha portato a estendere l'AC ad altre aree."*

#### Domanda 3 — Efficacia
*"L'estensione è stata sufficiente? Ci sono state NC simili dopo?"*

---

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"L'agente ha rilevato: [a] AC chiuse senza verifica di efficacia; [b] confusione fra attuazione ed efficacia; [c] AC da audit Accredia non gestite nel sistema interno; [d] NC ripetute indicano AC precedenti inefficaci; [e] mancata estensione a NC potenziali; [f] analisi cause radici superficiale (no metodo strutturato)."*
- Riferimento: ISO § 8.7 + RT-08 § 8.7
- Urgenza: ALTA (è una delle aree più frequentemente oggetto di NC in ispezione Accredia)
- Cosa NON può chiudere agente: chiusura NC/AC, valutazione efficacia, modifica registro

---

## E. Errori fatali correlati

1. **AC chiusa subito dopo attuazione senza verifica efficacia** → **NC § 8.7.1.d** (più frequente in Italia).
2. **NC ripetute identiche** → AC precedente non efficace, ulteriore NC sistemica.
3. **AC Accredia gestite "fuori sistema"** (Excel separato) → **NC § 8.7.1** (RT-08 esplicito).
4. **Analisi cause = causa apparente** (es. "il tecnico ha sbagliato" senza chiedersi perché) → analisi insufficiente.
5. **Nessuna estensione** a NC potenziali (es. errore su una bilancia, non si controllano le altre) → NC.
6. **AC sproporzionate** all'effetto NC (sia eccessive che inadeguate) → § 8.7.2.
7. **Registro NC/AC inesistente / incompleto / non aggiornato** → **NC sistemica grave**.
8. **NC con rischio NON aggiornato** dopo AC → § 8.7.1.e ignorato.
9. **Confondere correzione (immediata, sul prodotto/risultato) con azione correttiva** (sistemica, sulla causa).

---

## F. Quando l'agente DEVE chiedere prima di rispondere

- Se il laboratorio chiede di **chiudere una NC** → **escalation a QM**, decisione formale.
- Se chiede di valutare se un'AC è efficace → **escalation a QM**, richiede dati/KPI.
- Se chiede di decidere la **profondità** dell'analisi cause → **escalation a QM**.
- Se chiede di scegliere il metodo RCA (5 whys, Ishikawa, ecc.) → **suggerimento OK**, ma scelta operativa al QM.
- Se chiede di gestire una NC critica (es. impatto su RdP già emessi) → **escalation immediata**, collegamento § 7.10.
- Se chiede di "non registrare" una NC perché "minore" → **rifiutare e escalation a QM**, ogni NC va registrata.
