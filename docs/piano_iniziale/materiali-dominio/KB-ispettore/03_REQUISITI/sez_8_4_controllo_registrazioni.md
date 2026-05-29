---
codice: sez_8_4_controllo_registrazioni
tipo: requisito_normativo
livello: misto
titolo: "§ 8.4 — Controllo delle registrazioni (Opzione A)"
sezione_norma: "8.4"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "8.4" }
  - { doc: "RT-08", rev: "05", paragrafo: "8.4" }
fonti_secondarie:
  - { doc: "LS-04", rev: "20" }
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[02_TRE_LIVELLI_DI_CONTROLLO]]"
  - "[[sez_8_3_controllo_documenti]]"
tags: [registrazioni, conservazione, 48_mesi, tarature, backup, accessibilita]
escalation_qm: true
---

# § 8.4 — Controllo delle registrazioni (Opzione A)

> Il § 8.4 disciplina identificazione, conservazione, protezione, recupero, archiviazione, leggibilità ed eliminazione delle registrazioni tecniche e gestionali del laboratorio. È strettamente collegato al § 7.5 (registrazioni tecniche) ma copre TUTTE le registrazioni SGQ.

> ⚠️ **Disclaimer di ricostruzione KB**
> Questo file fa parte della ricostruzione della **sezione 8** sull'indice ufficiale **ISO/IEC 17025:2018**. La vecchia KB aveva numerazione errata mutuata da ISO 9001.

---

## A. Testo della norma (🟢 LIVELLO 1 — fatto normativo)

> **Fonte**: `[ISO-17025-2018 § 8.4]` — testo a pagamento, non riprodotto integralmente.

**Sintesi operativa per sotto-paragrafo**:

- **8.4.1** — Il laboratorio deve stabilire e conservare registrazioni leggibili per dimostrare il soddisfacimento dei requisiti della norma.
- **8.4.2** — Il laboratorio deve implementare i controlli necessari per identificazione, archiviazione, protezione, backup, accessibilità, tempi di conservazione ed eliminazione delle registrazioni. Le registrazioni devono essere mantenute per un periodo coerente con gli obblighi contrattuali. L'accesso a tali registrazioni deve essere coerente con gli impegni di riservatezza e prontamente disponibili.

---

## B. Prescrizioni Accredia (🟢 LIVELLO 2 — prescrizione Accredia)

> **Fonte**: `[RT-08 rev.05 § 8.4]`

### 8.4.1 (RT-08)
> *"Si applica il requisito di norma."*

### 8.4.2 (RT-08)
> *"Si applica il requisito di norma. Tutte le registrazioni devono essere conservate per un periodo minimo di 48 mesi, ove non esistano obblighi cogenti o contrattuali più onerosi, inclusi i documenti di registrazione esterna (es. certificati di taratura, certificati di materiali di riferimento, report di prove valutative interlaboratorio).*
> *I dati relativi alle tarature dovrebbero essere conservati per la vita dell'apparecchiatura di misura, in quanto in base a questi possono essere stabilite o modificate le frequenze di taratura."*

---

### 🔑🔑🔑 LA REGOLA DEI 48 MESI — prescrizione FORTE Accredia

> **TESTO RT-08**: tutte le registrazioni → **minimo 48 mesi** (4 anni), salvo obblighi cogenti/contrattuali più onerosi.

**Cosa include "tutte le registrazioni"** (da RT-08, esplicitamente):
- Certificati di taratura
- Certificati di materiali di riferimento (MR/MRC)
- Report di prove valutative interlaboratorio (PT/EQA)
- Registrazioni tecniche (dati grezzi, fogli di lavoro, RdP)
- Registrazioni gestionali (audit, riesami, NC, AC, reclami, formazione, ecc.)

**Casi di durata SUPERIORE ai 48 mesi**:
- **Settore alimentare**: spesso 5-10 anni in base a obblighi HACCP/cliente
- **Settore farmaceutico GMP**: tipicamente almeno 5 anni dopo scadenza prodotto
- **Acque potabili (D.Lgs. 18/2023)**: obblighi specifici, spesso 5-10 anni
- **Settore ambientale (rifiuti, scarichi)**: D.Lgs. 152/2006 — durate variabili
- **Settore forense / legale**: durate decennali possibili
- **Contratti specifici**: il cliente può richiedere conservazioni più lunghe (vanno verbalizzate in offerta/contratto)

> ⚠️ Il **minimo 48 mesi** vale per i requisiti di accreditamento; se obblighi cogenti/contrattuali sono più lunghi, **prevalgono questi ultimi**.

---

### 🔑🔑 TARATURE — conservazione PER LA VITA DELLO STRUMENTO

> **TESTO RT-08**: *"I dati relativi alle tarature dovrebbero essere conservati per la vita dell'apparecchiatura di misura, in quanto in base a questi possono essere stabilite o modificate le frequenze di taratura."*

**Razionale**:
- Lo **storico tarature** è l'unica base oggettiva per stabilire/aggiustare la **frequenza di taratura** (vedi § 6.4).
- Eliminando lo storico, si perde la capacità di valutare la **deriva** dello strumento.
- "Dovrebbero" (verbo condizionale): è raccomandazione forte, NON obbligo tassativo. Tuttavia in audit Accredia se non sono presenti, il laboratorio deve giustificare come stabilisce le frequenze di taratura.

**Implicazione operativa**:
- Lo strumento dismesso → storico tarature può essere archiviato (es. archivio storico) ma non distrutto fintanto che resta in uso.
- Strumento sostituito → storico del vecchio NON è più obbligatorio dopo i 48 mesi standard.

---

### 🔑 Requisiti di sistema per le registrazioni elettroniche

Se le registrazioni sono elettroniche (LIMS, fogli Excel, database):
1. **Backup periodici** (frequenza definita, verificata, testata almeno annualmente per il restore)
2. **Protezione da modifiche non autorizzate** (controllo accessi, ruoli, log di audit trail)
3. **Tracciabilità modifiche** (chi-cosa-quando, con possibilità di ripristino della versione precedente)
4. **Leggibilità nel tempo** (formato esportabile, no dipendenza da software obsoleto)
5. **Conformità GDPR** se contengono dati personali
6. **Validazione del sistema** (LIMS validato, foglio Excel critico con controllo formule)

---

### 🔑 Categorie tipiche di registrazioni in un lab ISO 17025

| Categoria | Esempi | Durata minima |
|---|---|---|
| Registrazioni tecniche | Fogli di lavoro, dati grezzi, RdP | 48 mesi |
| Tarature | Certificati, dati intermedi, conferme | Vita strumento (raccomandato) |
| Controlli qualità interni | CQI, carte di controllo, MR | 48 mesi |
| Prove valutative | Report PT, valutazione esiti | 48 mesi |
| Reclami | Registro reclami, evidenza gestione | 48 mesi |
| Non conformità e AC | Registro NC/AC, analisi cause | 48 mesi |
| Audit interni | Piano, rapporti, evidenze | 48 mesi |
| Riesami direzione | Verbali, input/output | 48 mesi |
| Formazione | Verbali, attestati, valutazioni efficacia | 48 mesi (consigliato vita lavorativa) |
| Contratti / offerte / riesami | Documenti contrattuali | 48 mesi (spesso più lunghi per fisco) |

---

## C. Documenti applicabili (LS-04 rev.20)

- **RG-02** — Regolamento generale accreditamento
- **RT-08 rev.05** — Prescrizioni Accredia
- **PG interna gestione registrazioni** (deve esistere)
- **Procedura backup** (se registrazioni elettroniche)
- **Procedura audit trail** (per LIMS validati)

---

## D. Domande ispettive (🟡 LIVELLO 3 — prassi ispettiva)

### 8.4.1 — Esistenza e leggibilità

#### Domanda 1 — Meccanismo
*"Quali categorie di registrazioni mantenete? Mostratemi la procedura."*

#### Domanda 2 — Estensione
*"Vorrei vedere la registrazione di [test specifico] effettuato in data [3 anni fa]."*
- Verificare reperibilità e leggibilità.

#### Domanda 3 — Efficacia
*"In quanto tempo riuscite a recuperare una registrazione di 4 anni fa? Provatelo davanti a me."*

### 8.4.2 — Conservazione e protezione

#### Domanda 1 — Meccanismo (48 mesi)
*"Quali tempi di conservazione applicate per le diverse categorie? Esistono obblighi contrattuali specifici per qualche cliente?"*

#### Domanda 2 — Estensione (tarature)
*"Mostratemi lo storico tarature degli ultimi 10 anni dello strumento [X]. Come usate questi dati per definire la frequenza di taratura?"*
- Trabocchetto: se non c'è storico → la frequenza di taratura è stabilita su quale base?

#### Domanda 3 — Efficacia (backup)
*"Quando è stato l'ultimo test di restore del backup? Esito? Chi l'ha eseguito? È documentato?"*
- Trabocchetto Accredia: backup esistono ma nessuno ha mai provato a fare il restore.

### 8.4.2 — Accessibilità e riservatezza

#### Domanda 1 — Meccanismo
*"Chi ha accesso alle registrazioni dei clienti? Come è gestita la riservatezza?"*

#### Domanda 2 — Estensione
*"Le registrazioni elettroniche hanno audit trail? Mostratemi un esempio di modifica tracciata."*

#### Domanda 3 — Efficacia
*"Un cliente vi chiede una registrazione di una sua analisi di 2 anni fa. In quanto tempo gliela fornite? Procedura?"*

---

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"L'agente ha rilevato: [a] registrazioni conservate per <48 mesi senza giustificazione; [b] storico tarature non disponibile/incompleto; [c] backup mai testato per restore; [d] audit trail non attivo / non funzionante; [e] accessibilità inadeguata o riservatezza compromessa."*
- Riferimento: ISO § 8.4 + RT-08 § 8.4
- Urgenza: ALTA (perdita registrazioni = perdita evidenza accreditamento)
- Cosa NON può chiudere agente: definizione/modifica politiche conservazione, eliminazione registrazioni, modifica configurazione LIMS

---

## E. Errori fatali correlati

1. **Distruzione registrazioni prima dei 48 mesi** (es. politica "tutto a 24 mesi" copiata da altri SGQ) → **NC maggiore**.
2. **Storico tarature inesistente o eliminato** dopo pochi anni → impossibile giustificare frequenze di taratura → **NC**.
3. **Backup esistenti ma mai testati** per restore → **NC** alla prossima prova (spesso il restore non funziona).
4. **Registrazioni elettroniche modificabili senza audit trail** → **NC grave** (compromette integrità del dato).
5. **Registrazioni illeggibili dopo qualche anno** (es. file in formato proprietario obsoleto, fogli di lavoro a matita scoloriti) → **NC**.
6. **Accessibilità eccessiva** (es. tutti i dipendenti vedono tutto, incluse offerte/dati clienti riservati) → **NC § 4.2 riservatezza**.
7. **Certificati di taratura o MR scartati** anche se MR/strumento ancora in uso → **NC**.
8. **Mancanza di registro reclami / NC / AC** → **NC sistemica grave**.

---

## F. Quando l'agente DEVE chiedere prima di rispondere

- Se il laboratorio chiede se può eliminare registrazioni "vecchie" → **escalation a QM**, decisione contrattuale/normativa.
- Se chiede di abbreviare i tempi di conservazione sotto i 48 mesi → **rifiutare**, è prescrizione Accredia forte.
- Se chiede di modificare/cancellare una registrazione esistente → **rifiutare e escalation a QM**, integrità del dato.
- Se chiede di gestire un caso di perdita registrazioni (es. crash server, file cancellato per errore) → **escalation immediata a QM**, gestione come NC.
- Se chiede di dare accesso a registrazioni di un cliente a un altro cliente / terza parte → **rifiutare e escalation a QM**, violazione § 4.2.
