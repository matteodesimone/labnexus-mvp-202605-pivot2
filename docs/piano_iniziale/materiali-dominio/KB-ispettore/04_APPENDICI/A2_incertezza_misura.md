---
codice: A2_incertezza_misura
tipo: appendice_tecnica
livello: misto
titolo: "Appendice A2 — Valutazione dell'incertezza di misura"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.6" }
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.8.3" }
  - { doc: "RT-08", rev: "05", paragrafo: "7.6" }
  - { doc: "JCGM-100", rev: "2008", note: "GUM" }
  - { doc: "UNI-CEI-70098-3", rev: "2016" }
  - { doc: "UNI-CEI-70099", rev: "VIM 2008/2012" }
fonti_secondarie:
  - { doc: "DT-0002", rev: "01" }
  - { doc: "DT-0002/1", rev: "01", note: "elettrico" }
  - { doc: "DT-0002/2", rev: "00", note: "meccanico" }
  - { doc: "DT-0002/3", rev: "00", note: "chimico avvertenze" }
  - { doc: "DT-0002/4", rev: "00", note: "chimico esempi" }
  - { doc: "DT-0002/5", rev: "01", note: "materiali strutturali" }
  - { doc: "DT-0002/6", rev: "00", note: "ripetibilità" }
  - { doc: "EURACHEM-CITAC-QUAM", rev: "3rd ed. 2012" }
  - { doc: "ISTISAN-03/30", rev: "2ª ed. 2000", note: "trad. QUAM" }
  - { doc: "Nordtest-TR-537", rev: "ed. 4, 2017" }
  - { doc: "EUROLAB-1/2007", rev: "2007" }
  - { doc: "ISO-21748", rev: "2017" }
  - { doc: "ILAC-G17", rev: "01/2021", note: "incertezza nelle prove — riferimento EA/ILAC vigente (EA-4/16 ritirata)" }
  - { doc: "ISO-19036", rev: "2020", note: "microbiologia" }
  - { doc: "IEC-Guide-115", rev: "2023", note: "elettrico/conformity assessment" }
  - { doc: "JCGM-106", rev: "2012", note: "incertezza e dichiarazioni di conformità" }
  - { doc: "EURACHEM-VIM3-Intro", rev: "2nd ed. 2023" }
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[00_GLOSSARIO]]"
  - "[[03_PROTOCOLLO_HANDOFF_QM]]"
  - "[[02_TRE_LIVELLI_DI_CONTROLLO]]"
  - "[[sez_7_6_incertezza]]"
  - "[[sez_7_8_presentazione_risultati]]"
  - "[[A1_validazione_metodi]]"
  - "[[A3_riferibilita_taratura]]"
  - "[[A5_prove_valutative_PT_ILC]]"
  - "[[A7_campionamento]]"
tags: [incertezza, GUM, MU, top_down, bottom_up, conformity_assessment, decision_rule]
escalation_qm: true
---

> **Marcatore di livello: 🟢🟡🟣** — Appendice trasversale. Contenuti normativi (livello 1+2), prassi ispettiva (livello 3), best practice tecnica (livello 4). I marcatori distinguono ogni blocco.

> **🟡 Ponte con la validazione (§ 7.2).** L'incertezza si costruisce sugli **input prodotti in validazione**: precisione (RSDr/RSDi), bias/recupero, contributi dominanti. RT-08 § 7.6.3 ammette l'uso degli scarti tipo di ripetibilità/riproducibilità del metodo previa verifica delle proprie prestazioni. Riferimento EA/ILAC vigente per l'UM nelle prove: `[ILAC-G17:01/2021]` (la storica `EA-4/16` è **ritirata**, vedi [[00_FONTI_NORMATIVE]] § 7). Approccio top-down: `[ISO-21748:2017]`, `[EURACHEM-CITAC-QUAM]`, `[Nordtest-TR-537]`. Vedi [[A1_validazione_metodi]] e [[sez_7_2_selezione_verifica_validazione_metodi]] § B-bis. **Il filo dev'essere tracciabile**: dato di validazione → componente → UM dichiarata sul RdP.

# Appendice A2 — Valutazione dell'incertezza di misura

> **Riassunto operativo.** L'incertezza di misura è il punto dove la maggior parte dei laboratori inciampa in visita Accredia: spesso esiste un calcolo, ma non esiste il **modello** né il **dossier**. Questa appendice serve al QM/RT per impostare la scelta del metodo (GUM vs top-down), strutturare la documentazione e raccordare l'incertezza alle dichiarazioni di conformità (§ 7.8.6) e alla decision rule. L'agente la carica quando il QM parla di "calcolo incertezza", "MU", "k=2", "decision rule", "guard band", "PT non superato per cui dobbiamo verificare l'incertezza", o quando un rapporto di prova deve riportare l'incertezza ex § 7.8.3.

---

## 1. Principio normativo (🟢 LIVELLO 1+2)

### 1.1 ISO 17025:2018 § 7.6 — Valutazione dell'incertezza di misura

> Testo a pagamento, non riprodotto. Riassunto operativo:

- **§ 7.6.1** — i laboratori devono identificare tutti i contributi all'incertezza di misura. Per la taratura si considerano tutti i contributi significativi, inclusi quelli derivanti dal campionamento (se rilevante).
- **§ 7.6.2** — un laboratorio che esegue tarature, comprese quelle delle proprie apparecchiature, deve valutare l'incertezza di misura.
- **§ 7.6.3** — un laboratorio di prova deve valutare l'incertezza di misura. Quando il metodo di prova preclude una valutazione rigorosa, il laboratorio deve almeno tentare una stima basata sulla comprensione dei principi teorici o sull'esperienza pratica del metodo.

### 1.2 ISO 17025:2018 § 7.8.3 — Quando l'incertezza si riporta nel rapporto di prova

- Va riportata quando rilevante per la validità o l'applicazione dei risultati, quando richiesta dall'istruzione del cliente, oppure quando influenza la conformità a un limite di specifica.
- Va riportata nella stessa unità della misurazione (o come unità relativa, es. %).

### 1.3 RT-08 rev.05 § 7.6 (free, citabile)

> **Fonte**: `[RT-08 rev.05 § 7.6 — EC 10-02-2022]`

RT-08 rinvia alla norma e ricorda:
- L'incertezza deve essere **documentata** in un dossier specifico, riesaminato periodicamente.
- Per la chimica, il laboratorio applica `DT-0002/3` (avvertenze) e `DT-0002/4` (esempi).
- Per l'elettrico `DT-0002/1`, per il meccanico `DT-0002/2` e `DT-0002/5`, per la ripetibilità nel tempo `DT-0002/6`.
- `DT-0002 rev.01` è la guida generale che mette in fila il quadro GUM in modo accessibile.

### 1.4 Vocabolario di riferimento (VIM)

> **Fonte**: `[UNI-CEI-70099 / VIM 2008/2012]`

L'agente, quando usa i termini "incertezza standard", "incertezza combinata", "incertezza estesa", "fattore di copertura", "contributo di tipo A/B", "ripetibilità", "precisione intermedia", "riproducibilità", "veridicità", li riferisce sempre al VIM. Per l'introduzione divulgativa: `EURACHEM-VIM3-Intro 2023` (free).

---

## 2. Come un ispettore lo verifica (🟡 LIVELLO 3 — prassi ispettiva)

> Pattern: meccanismo → estensione → efficacia.

### 2.1 Domanda 1 — Meccanismo

**"Per la prova X dello scopo, mi mostrate il dossier di incertezza? Voglio: modello di misurazione (mensurando + equazione + contributi identificati), tipo di approccio applicato (GUM bottom-up, top-down su PT, ISO 21748 da R&r), valori dei contributi, incertezza combinata, fattore di copertura, incertezza estesa, firma dell'autorizzato."**

Analisi ispettiva:
- Un foglio Excel che dà un numero senza spiegare il modello è sospetto. Cerco il **modello matematico** o, se top-down, l'evidenza dei dati di precisione/veridicità su cui poggia.
- Verifico la coerenza con `DT-0002` rev.01 e la guida settoriale (DT-0002/x). Se il lab è chimico e cita solo "GUM" senza riferimento a `DT-0002/3-4` o `QUAM`, segnale debole.
- Controllo la **firma**: chi ha approvato il dossier? Deve essere personale autorizzato (cfr. § 6.2.6 — l'autorizzazione su "analisi risultati" è prerequisito).

### 2.2 Domanda 2 — Estensione

**"Per quanti dei vostri metodi accreditati esiste un dossier di incertezza completo? Esiste mappa di copertura metodo ↔ dossier? Per i metodi a matrici/range multipli, l'incertezza è dichiarata in modo che copra tutti i casi?"**

Analisi ispettiva:
- Campionamento sullo scopo: 3-5 metodi, dossier richiesto. Mancanze = NC su § 7.6.
- Errore frequente: l'incertezza calcolata per la matrice "tipica" applicata anche alle altre matrici senza giustificazione. Se le matrici hanno comportamento diverso (recupero, interferenze), l'incertezza non è la stessa.
- Verifico se l'incertezza è **dichiarata nei rapporti** quando dovuta (§ 7.8.3) — incrociando dossier e rapporti emessi.

### 2.3 Domanda 3 — Efficacia

**"Il vostro dossier di incertezza è coerente con i risultati dei PT/ILC degli ultimi 3 anni? Quando un PT è stato 'questionable' o 'unsatisfactory', avete verificato se l'incertezza dichiarata era sottostimata? Mi mostrate un caso."**

Analisi ispettiva:
- L'incertezza vive nel confronto con il mondo reale. Un PT con z-score |z| > 2 ricorrente con incertezza dichiarata "molto stretta" è incongruente: o l'incertezza è sottostimata, o il bias è sistematico.
- Cerco la **catena**: PT scadente → riesame incertezza → eventuale revisione del dossier. Se i dossier sono fermi da anni nonostante i PT, l'incertezza è ornamentale.
- Verifico la coerenza con la **decision rule** (§ 7.8.6.2) sulle dichiarazioni di conformità (vedi § 3.5 di questa appendice).

### 2.4 Pattern di errori comuni

1. **Incertezza dichiarata senza modello.** Il rapporto riporta "U = X%" ma in cartella non c'è il dossier. NC su § 7.6 + § 7.8.3.
2. **Riproducibilità interlaboratorio (sR) usata come unica fonte top-down, senza considerare il bias del laboratorio.** Approccio incompleto. La guida `ISO-21748` raccomanda di combinare sR con bias e contributi specifici.
3. **k=2 dichiarato meccanicamente senza riferimento al livello di confidenza** (~95% assumendo distribuzione normale). Va precisato.
4. **Incertezza non riportata nei rapporti** quando il risultato è confrontato con un limite di specifica (§ 7.8.3.1 c). Errore frequente nei settori cogente (alimentari, ambiente).
5. **Dichiarazione di conformità senza decision rule documentata** (§ 7.8.6.2). NC indipendente dal § 7.6 ma direttamente collegata.
6. **Per la microbiologia: confusione tra incertezza ex `ISO-19036` (basata su sR) e "errore tecnico" del conteggio.** Va seguito l'approccio specifico.
7. **Per l'elettrico: omissione del contributo del banco di taratura/strumento di riferimento** nella catena (vedi `IEC-Guide-115`).

---

## 3. Best practice tecnica (🟣 LIVELLO 4 — best practice)

### 3.1 I due approcci principali

#### 3.1.1 Approccio bottom-up (GUM / QUAM)

Riferimenti: `JCGM-100 (GUM) 2008`, `UNI-CEI-70098-3:2016`, `EURACHEM-CITAC-QUAM 3rd ed. 2012`, `ISTISAN-03/30` (trad. it. QUAM), `DT-0002 rev.01`, `DT-0002/3-4` (chimica), `DT-0002/1-2-5` (settoriali).

**Logica**:
1. Definire il **mensurando** (cosa si sta misurando, in che condizioni).
2. Scrivere il **modello di misurazione** (equazione che lega il mensurando alle grandezze in ingresso).
3. Identificare tutte le **grandezze in ingresso** che contribuiscono significativamente.
4. Per ogni contributo: valutare l'incertezza **tipo A** (statistica, ripetizioni) o **tipo B** (valutazione su dati pregressi, certificati, manuali, distribuzioni assunte).
5. Combinare i contributi (radice della somma dei quadrati, con eventuali coefficienti di sensibilità dal modello) → **incertezza standard combinata u_c**.
6. Moltiplicare per il **fattore di copertura k** (tipicamente k=2 per ~95% di livello di confidenza, distribuzione normale) → **incertezza estesa U**.

**Quando preferirlo**: tarature, metodi chimici ben modellabili, situazioni in cui i contributi sono noti.

#### 3.1.2 Approccio top-down / empirico

Riferimenti: `ISO-21748:2017`, `Nordtest-TR-537 ed.4 2017`, `EUROLAB-1/2007`, `EURACHEM-CITAC-QUAM 3rd ed.` (cap. dedicato), `ISO-19036:2020` (microbio).

**Logica**:
- Si stima l'incertezza a partire dai dati di **prestazione complessiva** del metodo: ripetibilità (sr), precisione intermedia (sIP), riproducibilità interlaboratorio (sR), bias del laboratorio (da CRM/PT/spike), eventuale contributo dal campionamento.
- `ISO-21748` formalizza la combinazione di sR (dallo studio collaborativo o dalle prove valutative) con un termine di bias del laboratorio.
- `Nordtest-TR-537` formalizza l'approccio per laboratori ambientali (anche con dati PT/ILC interni).

**Quando preferirlo**: metodi complessi non modellabili in modo rigoroso, abbondanti dati storici di precisione/PT, microbiologia, ambiente.

#### 3.1.3 Combinazione dei due

`EUROLAB-1/2007` e l'attuale prassi accreditata raccomandano un approccio **ibrido**: top-down per la base, bottom-up per i contributi che il top-down non cattura (es. contributo del campionamento sull'aliquota, contributi specifici del trattamento).

### 3.2 Guide settoriali Accredia (DT-0002/x — livello 2, free)

| Documento | Settore | Cosa fornisce |
|---|---|---|
| `DT-0002 rev.01` | Generale | Guida applicativa Accredia all'espressione dell'incertezza |
| `DT-0002/1 rev.01` | Misurazioni elettriche | Esempi applicativi |
| `DT-0002/2 rev.00` | Misurazioni meccaniche | Esempi applicativi |
| `DT-0002/3 rev.00` | Analisi chimica | Avvertenze (cosa il chimico spesso sbaglia) |
| `DT-0002/4 rev.00` | Misurazioni chimiche | Esempi applicativi |
| `DT-0002/5 rev.01` | Materiali strutturali | Esempio applicativo |
| `DT-0002/6 rev.00` | Trasversale | Calcolo della ripetibilità e verifica nel tempo |

> Citare `DT-0002` nel dossier di incertezza non è obbligo, ma in visita Accredia un dossier che si appoggia esplicitamente alla guida Accredia di settore è **molto più solido** di uno che cita solo GUM/QUAM.

### 3.3 Guide internazionali e regionali

| Documento | Quando è particolarmente utile |
|---|---|
| `JCGM-100` (GUM) 2008 | Base teorica universale; sempre citabile |
| `UNI-CEI-70098-3:2016` | Recepimento italiano GUM; per coerenza con il quadro UNI |
| `EURACHEM-CITAC-QUAM 3rd ed. 2012` | Riferimento d'elezione per la chimica analitica |
| `ISTISAN-03/30` 2000 | Traduzione italiana di QUAM (2ª ed.) — utile per lab di lingua italiana |
| `Nordtest-TR-537` ed. 4 2017 | Laboratori ambientali; approccio pragmatico con PT/ILC |
| `EUROLAB-1/2007` | Sintesi degli approcci alternativi GUM/top-down |
| `ISO-21748:2017` | Combinazione di R&r interlab con bias del laboratorio |
| `ISO-19036:2020` | Microbiologia alimentare — specifico per conteggi |
| `IEC-Guide-115:2023` | Elettrico — incertezza e conformity assessment |

### 3.4 Settori specialistici

#### 3.4.1 Microbiologia
`ISO-19036:2020` definisce un approccio basato su precisione intermedia/riproducibilità intra-laboratorio dei conteggi (trasformazione log). Non si applica il modello GUM classico. La pratica accreditata è dichiarare l'incertezza di conteggio come deviazione tipo del log10 (s'_R o sR') con fattore di copertura concordato.

#### 3.4.2 Elettrico (taratura/prova)
`IEC-Guide-115:2023` integra GUM con i contributi tipici del settore (catena di taratura, drift, condizioni ambientali, EMC residua). Il `DT-08-DL` Accredia copre lo specifico EMC. Vedi `[[A3_riferibilita_taratura]]`.

#### 3.4.3 Meccanico/strutturale
`DT-0002/2` e `DT-0002/5` sono il riferimento operativo. I contributi tipici includono: incertezza del riferimento (CRM o macchina campione), risoluzione strumentale, ripetibilità, condizioni ambientali, operatore.

#### 3.4.4 Campionamento
Quando il laboratorio è accreditato anche per il campionamento, l'incertezza del campionamento si combina con quella analitica. Vedi `[[A7_campionamento]]` e `EURACHEM-Sampling 2nd ed. 2019` / `ISTISAN-22/39:2023` (trad. it.) / `Nordtest-TR-537 § dedicato`.

### 3.5 Incertezza e dichiarazioni di conformità (decision rule)

Riferimenti: `ISO-17025-2018 § 7.8.6.2`, `JCGM-106:2012`, `IEC-Guide-115:2023`.

Quando il laboratorio dichiara conformità a una specifica (es. "valore ≤ limite"), deve dichiarare la **decision rule** applicata, che tiene conto dell'incertezza. Le opzioni più comuni:

| Decision rule | Significato | Note |
|---|---|---|
| **Binary (simple acceptance)** | Conforme se valore misurato ≤ limite, indipendentemente da U | Ammesso solo se il cliente lo accetta esplicitamente e se U è trascurabile rispetto alla soglia |
| **Guard band (w = k·U) — non-conformità protetta** | Conforme solo se valore + U ≤ limite | Riduce il rischio di falsi conformi (favorisce il "produttore") |
| **Guard band — conformità protetta** | Non-conforme solo se valore − U > limite | Riduce il rischio di falsi non-conformi (favorisce il "consumatore") |
| **Conditional / con zona di incertezza** | Tre esiti: conforme, non-conforme, indeterminato (entro ± U dal limite) | Esplicita la zona dove l'incertezza non permette dichiarazione |

> Senza decision rule documentata e comunicata al cliente, una dichiarazione di conformità è formalmente non conforme al § 7.8.6.2.

---

## 4. Template operativo per il QM

### 4.1 Decisione preliminare — Quale approccio applicare?

| Situazione | Approccio raccomandato | Riferimenti citabili |
|---|---|---|
| Taratura strumentale, modello chiaro | Bottom-up GUM | `JCGM-100`, `UNI-CEI-70098-3`, `DT-0002` |
| Chimica analitica con metodo modellabile | Bottom-up QUAM | `EURACHEM-CITAC-QUAM`, `DT-0002/3-4`, `ISTISAN-03/30` |
| Chimica con metodo complesso e abbondanti dati PT/precisione | Top-down ISO 21748 / Nordtest | `ISO-21748`, `Nordtest-TR-537`, `DT-0002/3` |
| Microbiologia (conteggi) | Approccio specifico ISO 19036 | `ISO-19036`, `EURACHEM-MicroAcc` |
| Elettrico | Bottom-up GUM + specifico settoriale | `DT-0002/1`, `IEC-Guide-115` |
| Meccanico/strutturale | Bottom-up GUM con DT settoriale | `DT-0002/2`, `DT-0002/5` |
| Metodo non modellabile e nessun dato collaborativo | Stima cautelativa basata su esperienza + PT/ILC dedicato (cfr. § 7.6.3 ultima frase) | `ISO-17025-2018 § 7.6.3` |

### 4.2 Template "Dossier incertezza di misura"

```
DOSSIER INCERTEZZA DI MISURA — METODO [codice]
Rev. [n] del [data]   Redatto: [tecnico]   Approvato: [RT autorizzato § 6.2.6]

1. IDENTIFICAZIONE
   - Metodo: [codice + edizione + matrice + range]
   - Mensurando: [grandezza, condizioni, unità]
   - Uso previsto del risultato: [decisione, criterio]

2. APPROCCIO METODOLOGICO
   [ ] Bottom-up GUM/QUAM
   [ ] Top-down ISO 21748
   [ ] Top-down Nordtest TR-537
   [ ] Ibrido (specificare)
   [ ] Specifico settoriale (ISO 19036 / IEC Guide 115 / DT-0002/x)
   - Riferimenti citati: [elenco]

3. MODELLO DI MISURAZIONE (se bottom-up)
   - Equazione: y = f(x1, x2, …, xN)
   - Coefficienti di sensibilità: ∂f/∂xi

4. CONTRIBUTI ALL'INCERTEZZA
   | Sorgente | Tipo (A/B) | Distribuzione | u_i | Coeff. sens. c_i | c_i · u_i |
   |----------|-----------|---------------|-----|------------------|-----------|
   | …        |           |               |     |                  |           |

5. INCERTEZZA COMBINATA E ESTESA
   - u_c = ...
   - k = ... (livello di confidenza ~95%, distribuzione assunta: ...)
   - U = k · u_c = ...
   - Espressione finale: y = (valore) ± U (unità), k = ...

6. VALIDAZIONE DELL'INCERTEZZA DICHIARATA
   - Confronto con PT/ILC degli ultimi … anni: [tabella z-score / E_n]
   - Coerenza con dati di precisione intermedia: [esito]
   - Eventuali bias rilevati: [valore + trattamento]

7. PIANO DI RIESAME
   - Trigger: PT scadente, modifica metodo, cambio strumento, deriva CQ
   - Frequenza minima: [annuale / al ciclo Accredia]

8. UTILIZZO NEI RAPPORTI DI PROVA
   - Quando si riporta U: [criterio applicato ex § 7.8.3]
   - Decision rule applicata per dichiarazioni di conformità: [vedi A2 § 3.5]
```

### 4.3 Checklist minima per audit interno su § 7.6

- [ ] Esiste un dossier incertezza per ogni metodo accreditato (mappa di copertura)
- [ ] Per ogni dossier: modello/approccio esplicitato, contributi quantificati, U dichiarata con k
- [ ] Approvazione del dossier da personale autorizzato § 6.2.6
- [ ] Coerenza con PT/ILC degli ultimi 3 anni (z-score, E_n)
- [ ] Coerenza con dati di precisione intermedia interna
- [ ] Riesame periodico documentato (almeno annuale o al ciclo Accredia)
- [ ] Incertezza riportata nei rapporti nei casi previsti da § 7.8.3
- [ ] Decision rule documentata e comunicata al cliente per ogni dichiarazione di conformità (§ 7.8.6.2)
- [ ] Per i metodi su matrici multiple: dimostrazione che l'incertezza copre tutti i casi o dichiarazione differenziata

### 4.4 Format minimo nei rapporti di prova

| Componente | Esempio |
|---|---|
| Valore misurato + U + k | `Pb = 0.045 ± 0.009 mg/L (U, k=2, livello di confidenza ~95%)` |
| Riferimento al dossier | `Stima incertezza ex dossier MU-[codice] rev. [n]` |
| Dichiarazione di conformità (se richiesta) | `Conforme al limite di X mg/L (decision rule: guard band protezione consumatore w = U; cfr. accordo cliente del [data])` |

---

## 5. Quando l'agente attiva questa appendice

1. Il QM chiede *"come calcolo l'incertezza per…"*, *"GUM o top-down?"*, *"posso usare la sR del metodo normato?"*.
2. Un PT è "questionable" o "unsatisfactory" e l'agente deve istruire la verifica della coerenza con l'incertezza dichiarata (collegamento con `[[A5_prove_valutative_PT_ILC]]`).
3. È in lavorazione un rapporto di prova con dichiarazione di conformità (§ 7.8.6) e non risulta decision rule documentata.
4. Si sta validando un nuovo metodo (collegamento con `[[A1_validazione_metodi]]`) e l'incertezza è un parametro del dossier di validazione.
5. È in preparazione la visita Accredia e l'agente costruisce l'istruttoria su § 7.6.
6. Si è introdotto un nuovo strumento o si è cambiata la catena di taratura: l'incertezza va riesaminata (collegamento con `[[A3_riferibilita_taratura]]`).
7. Per laboratori accreditati anche al campionamento: serve raccordare l'incertezza analitica con quella da campionamento (collegamento con `[[A7_campionamento]]`).

---

## 6. Cosa l'agente NON può chiudere su questo tema (handoff QM obbligatorio)

Vedi `[[03_PROTOCOLLO_HANDOFF_QM]] § 4`. In materia di incertezza, l'agente **non** può:

1. **Approvare il dossier di incertezza**. La firma di approvazione spetta al personale autorizzato (RT) ex § 6.2.6.
2. **Decidere l'approccio metodologico definitivo** (bottom-up vs top-down) per un metodo specifico. L'agente propone, il RT decide.
3. **Validare la decision rule** verso il cliente. La decision rule deve essere concordata con il cliente e approvata dal RT/Direzione (vedi § 7.8.6.2 + autorizzazione § 6.2.6).
4. **Dichiarare un rapporto conforme** ad un limite di specifica. È atto del personale autorizzato a "dichiarazioni di conformità" (§ 6.2.6).
5. **Decidere che un'incertezza dichiarata è "adeguata"** in assenza di evidenze comparative (PT, CRM, precisione intermedia). Resta giudizio tecnico del RT.
6. **Rivalutare l'incertezza di risultati già emessi** che si scoprono affetti da sottostima. La rivalutazione è del RT e la comunicazione al cliente è atto della Direzione (vedi `[[03_RISPOSTA_NC]]`).
7. **Dichiarare chiusa una NC su § 7.6 o § 7.8.3-6**. La chiusura è del QM.

> **Nota di sistema**: l'incertezza è uno dei punti dove l'illusione di rigore è massima (numeri, formule, k=2) ma la sostanza è spesso fragile. L'agente, su questo tema, opera sempre in modalità "istruttoria + domanda al QM".
