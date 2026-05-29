---
codice: sez_7_6_incertezza
tipo: requisito_normativo
livello: misto
titolo: "§ 7.6 — Valutazione dell'incertezza di misura"
sezione_norma: "7.6"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.6" }
  - { doc: "RT-08", rev: "05", paragrafo: "7.6" }
fonti_secondarie:
  - { doc: "JCGM-100", rev: "2008", paragrafo: "GUM" }
  - { doc: "ILAC-G17", rev: "01/2021", paragrafo: "incertezza nelle prove (riferimento EA/ILAC)" }
  - { doc: "DT-0002", rev: "01", paragrafo: "guida Accredia incertezza" }
  - { doc: "DT-0002/3", rev: "00", paragrafo: "incertezza analisi chimica" }
  - { doc: "DT-0002/4", rev: "00", paragrafo: "esempi chimici" }
  - { doc: "DT-0002/6", rev: "00", paragrafo: "ripetibilità" }
  - { doc: "EURACHEM-CITAC-QUAM", rev: "3rd ed. 2012", paragrafo: "QUAM" }
  - { doc: "Nordtest-TR-537", rev: "ed.4 2017", paragrafo: "ambiental" }
  - { doc: "EA-4/02", paragrafo: "tarature interne" }
  - { doc: "ISO-21748", rev: "2017", paragrafo: "uso r/R per UM" }
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[02_TRE_LIVELLI_DI_CONTROLLO]]"
  - "[[sez_7_2_selezione_verifica_validazione_metodi]]"
  - "[[sez_7_3_campionamento]]"
  - "[[sez_7_7_assicurazione_validita]]"
  - "[[sez_7_8_presentazione_risultati]]"
  - "[[A2_incertezza_misura]]"
tags: [incertezza_misura, GUM, top_down, bottom_up, contributi, ripetibilita, riproducibilita, fattore_copertura]
escalation_qm: true
---

# § 7.6 — Valutazione dell'incertezza di misura

> Il § 7.6 chiede al laboratorio di **identificare i contributi** all'incertezza e di **valutarli**. La norma non impone una metodologia (bottom-up GUM vs top-down basato su r/R + bias): impone che il lab abbia un approccio coerente, documentato, e che l'incertezza che riporta sia coerente coi dati prestazionali disponibili in letteratura. Per Accredia, l'incertezza è anche uno strumento di **decisione di conformità** (collega regola decisionale § 7.1.3 e § 7.8.6).

## A. Testo della norma (🟢 LIVELLO 1 — fatto normativo)

> **Fonte**: `[ISO-17025-2018 § 7.6.1 - 7.6.3]` — testo a pagamento, non riprodotto. Riassunto operativo:

- **7.6.1** Il laboratorio deve **identificare i contributi** all'incertezza di misura. Quando si valuta l'incertezza, tutti i contributi significativi, inclusi quelli derivanti dal campionamento, devono essere presi in considerazione usando metodi appropriati di analisi.
- **7.6.2** Un laboratorio che esegue **tarature**, incluse di proprio strumenti, deve valutare l'incertezza per tutte le tarature.
- **7.6.3** Un laboratorio che esegue **prove** deve valutare l'incertezza di misura. Quando il metodo preclude una valutazione rigorosa, una stima deve essere fatta basata sulla comprensione dei principi teorici o sull'esperienza pratica dell'esecuzione del metodo (Nota 1: in alcuni casi il metodo di prova consolidato e ben riconosciuto pone limiti ai valori delle principali fonti di incertezza e fissa la forma di presentazione dei risultati calcolati — in tali casi il laboratorio è considerato soddisfare il requisito seguendo il metodo e le istruzioni di reporting. Nota 2: per un metodo specifico, dove l'incertezza è stabilita e verificata, non è necessario rivalutare l'incertezza per ciascun risultato se il laboratorio può dimostrare che i fattori critici di influenza sono sotto controllo).

## B. Prescrizioni Accredia (🟢 LIVELLO 2 — prescrizione Accredia)

> **Fonte**: `[RT-08 rev.05 § 7.6]`

### 7.6.1 — Campionamento accreditato come fonte di incertezza

> **Citazione integrale RT-08**:
> *"Se il Laboratorio richiede in accreditamento la sola attività di campionamento (vedere § 7.3.) deve rendere disponibili le informazioni necessarie per il successivo calcolo dell'incertezza di misura associata al risultato."*

Cioè: se il lab fa solo il campionamento, deve fornire al lab di prova le info per stimare l'UM dal campionamento (collega [[sez_7_3_campionamento]] e EURACHEM-Sampling).

### 7.6.2 — Tarature interne

> **Citazione integrale RT-08**:
> *"Se il Laboratorio effettua tarature interne, l'incertezza deve essere determinata in accordo alla guida JCGM 100, ovvero alla Guida EA-4/02, o secondo quanto previsto dalle norme tecniche di settore."*

Riferimento esplicito a **JCGM 100 (GUM)** o **EA-4/02** o norme settoriali.

### 7.6.3 — Uso dei parametri di validazione per stimare UM

> **Citazione integrale RT-08**:
> *"Quando un metodo normalizzato/non normalizzato/ufficiale riporta i parametri statistici della validazione (scarto tipo di ripetibilità e scarto tipo di riproducibilità) e il Laboratorio decide di utilizzarli per il calcolo dell'incertezza di misura, deve verificare almeno che le proprie prestazioni sono compatibili con quelle indicate (es. ripetibilità, esattezza)."*

### 7.6.3 — Avvertenze specifiche per chimica

> **Citazione integrale RT-08**:
> *"Per le prove chimiche, è importante verificare che il livello di concentrazione per cui viene riportata la riproducibilità sia prossimo al risultato della prova, in quanto l'incertezza (e la riproducibilità che ne fornisce una stima) potrebbe essere una funzione non lineare della concentrazione (es. legge di Horwitz); non è pertanto sempre corretto esprimere l'incertezza in termini relativi su tutto il campo di applicazione di un metodo, quando questa è stata determinata solo su un valore di concentrazione."*

### 7.6.3 — Valori significativi per il cliente e confronto col limite di ripetibilità

> **Citazione integrale RT-08**:
> *"Il Laboratorio deve controllare che l'incertezza da associare al risultato sia stata valutata a livelli significativi per il cliente (es. limiti di legge) e inoltre deve confrontare il limite di ripetibilità, del metodo o calcolato dal Laboratorio, con il valore dell'incertezza di misura (2U > r)."*

Cioè: la **regola di sanity check** `2U > r` (incertezza estesa raddoppiata maggiore del limite di ripetibilità r). Se `2U < r`, l'incertezza è sottostimata.

### 7.6.3 — Congruenza coi dati di letteratura

> **Citazione integrale RT-08**:
> *"E' importante inoltre verificare che l'incertezza valutata sia congruente con i dati prestazionali disponibili in letteratura, ove esistenti (es. riportati nel metodo, ricavati dalla partecipazione ad appropriati confronti interlaboratorio, ecc.)."*

### 7.6.3 — Documenti di riferimento

> **Citazione integrale RT-08**:
> *"Alcuni documenti di riferimento per la valutazione dell'incertezza di misura sono elencati nella nota 3 della norma. Per settori specifici sono disponibili appositi documenti e linee guida."*

## B-bis. Ponte con la validazione (§ 7.2) — 🟡 integrazione

> L'incertezza non nasce da sé: i suoi ingredienti (precisione RSDr/RSDi, bias/recupero, contributi dominanti) sono prodotti in **validazione/verifica** (§ 7.2). RT-08 § 7.6.3 lo rende esplicito quando consente di usare gli scarti tipo di ripetibilità/riproducibilità del metodo per stimare l'UM, **previa verifica delle proprie prestazioni**. Il riferimento applicativo EA/ILAC è `[ILAC-G17:01/2021]`; per l'approccio top-down vedi `[ISO-21748:2017]`, `[EURACHEM-CITAC-QUAM]` e [[A2_incertezza_misura]]. Per il percorso inverso (dal dato di validazione all'UM) vedi [[sez_7_2_selezione_verifica_validazione_metodi]] § B-bis. **Il filo dev'essere tracciabile**: dato di validazione → componente di incertezza → UM dichiarata sul RdP.

## C. Documenti applicabili (LS-04 rev.20)

- `ISO-17025-2018` § 7.6.
- `RT-08 rev.05` § 7.6.
- `JCGM-100` / GUM 2008 — riferimento metrologico fondamentale (free BIPM).
- `ILAC-G17:01/2021` — Measurement Uncertainty in Testing: riferimento EA/ILAC che riflette l'approccio degli enti di accreditamento alla valutazione/espressione dell'UM nelle prove (free, ilac.org). NB: la storica `EA-4/16` è ritirata (vedi [[00_FONTI_NORMATIVE]] § 7).
- `EA-4/02` — guida EA per tarature (UM in taratura).
- `DT-0002 rev.01` — Guida Accredia valutazione/espressione incertezza.
- `DT-0002/1-6` — esempi applicativi per settori (elettrico, meccanico, chimico, materiali) + DT-0002/6 ripetibilità.
- `EURACHEM-CITAC-QUAM` 3rd ed. 2012 — Quantifying Uncertainty in Analytical Measurement.
- `ISTISAN-03/30` — traduzione italiana QUAM.
- `Nordtest-TR-537` ed.4 2017 — incertezza in laboratori ambientali (approccio top-down).
- `EUROLAB-1/2007` — approcci alternativi all'incertezza.
- `UNI-CEI-70098-3` — incertezza, parte 3.
- `ISO-21748:2017` — uso di r/R/trueness nella stima dell'incertezza.
- `ISO-19036:2020` — incertezza in microbiologia alimenti.
- `IEC-Guide-115:2023` — incertezza elettrotecnica.
- `EURACHEM-Sampling` 2nd ed. 2019 + `ISTISAN-22/39` — incertezza da campionamento.

## D. Domande ispettive (🟡 LIVELLO 3 — prassi ispettiva)

> Pattern: meccanismo → estensione → efficacia.

### 7.6.1 — Identificazione contributi e metodologia

#### Domanda 1 — Meccanismo
**"Per un metodo accreditato (scelgo io): qual è l'approccio scelto per l'UM? Bottom-up GUM (modello + budget), top-down (basato su r, R, bias da PT/CRM), ibrido? Su quale guida vi siete basati? Avete un documento di budget dell'incertezza?"**

Analisi ispettiva: la norma non impone una metodologia ma chiede coerenza. Per la chimica analitica spesso è top-down (Nordtest TR-537, ISTISAN) o ibrido. Per la metrologia fisica spesso bottom-up GUM. Per la microbiologia c'è ISO 19036. Il laboratorio deve poter giustificare la scelta. Tipico problema: usa la sola r di ripetibilità intermedia spacciandola per incertezza estesa.

#### Domanda 2 — Estensione
**"Mostratemi il budget dell'incertezza (contributi, tipo A/tipo B, distribuzioni, fattori di sensibilità) per questo metodo. I contributi del campionamento e del subcampionamento sono inclusi? La componente "matrice" è considerata?"**

Analisi ispettiva: il documento di budget deve esistere ed essere ragionato. Controllo:
- Sono inclusi i contributi dominanti? (es. recupero, bias, matrice, taratura strumento)
- Se il lab fa anche campionamento accreditato, il contributo da campionamento è inserito (RT-08 esplicito al § 7.6.1)?
- Sono dichiarate distribuzioni e fattori (rettangolare, triangolare, normale)?

#### Domanda 3 — Efficacia
**"Avete mai aggiornato l'incertezza dopo: (a) un cambio strumento, (b) un cambio reagente critico, (c) una serie di PT con esito borderline, (d) un cambio matrice? Mostratemi un caso."**

Analisi ispettiva: la stima di UM va riesaminata quando cambiano i fattori critici di influenza (Nota 2 della norma). Tipico problema: l'incertezza è stata calcolata in fase di validazione 7 anni fa e da allora nessuno l'ha rivista.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Metodologia UM scelta + budget documentato + inclusione contributi campionamento/matrice + riesame in seguito a cambi. ISO § 7.6.1 + RT-08 § 7.6.1."*
- Urgenza: alta — l'incertezza inaffidabile invalida la decisione di conformità sul RdP.

### 7.6.2 — Tarature interne

#### Domanda 1 — Meccanismo
**"Effettuate tarature interne (anche solo di vostri strumenti)? Se sì, su quale guida vi basate per l'UM: GUM JCGM-100 o EA-4/02? Avete procedura scritta?"**

Analisi ispettiva: RT-08 esplicito su questo. Se il lab fa tarature interne, la base è GUM/EA-4/02. Se la risposta è "stimiamo con il metodo del costruttore", verificare che la procedura del produttore sia documentata e che il risultato sia coerente con GUM.

#### Domanda 2 — Estensione
**"Mostratemi un certificato di taratura interna recente e il relativo budget di incertezza."**

Analisi ispettiva: completezza certificato, contributi tipici (deriva strumento di riferimento, ripetibilità, risoluzione, ambiente).

#### Domanda 3 — Efficacia
**"L'UM delle vostre tarature interne è coerente con quella delle tarature esterne accreditate equivalenti?"**

Analisi ispettiva: sanity check.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Tarature interne: aderenza a JCGM-100/EA-4/02, budget UM, coerenza con tarature accreditate. ISO § 7.6.2 + RT-08 § 7.6.2."*

### 7.6.3 — Uso dei parametri di validazione + sanity check 2U>r

#### Domanda 1 — Meccanismo
**"Se usate la riproducibilità del metodo (scarto tipo di riproducibilità R) per stimare l'incertezza: avete verificato che le vostre prestazioni di ripetibilità siano coerenti con quelle del metodo? Avete eseguito anche un check di esattezza (recupero o materiali di riferimento)?"**

Analisi ispettiva: RT-08 esplicito. Usare R del metodo richiede verifica delle proprie prestazioni. Se mancano i dati di verifica, l'incertezza è dichiarata ma non giustificata.

#### Domanda 2 — Estensione (sanity check Accredia)
**"Per 3 metodi accreditati: confrontate `2U` (incertezza estesa) con `r` (limite di ripetibilità del metodo o calcolato dal lab). La regola RT-08 chiede `2U > r`. Soddisfatta?"**

Analisi ispettiva: prova del nove. Se `2U < r`, l'incertezza è SOTTOSTIMATA. È una NC tecnica.

#### Domanda 3 — Efficacia
**"Per la chimica: la riproducibilità del metodo è stata determinata a quale livello di concentrazione? È vicino al livello di routine dei vostri campioni? Se no, come avete gestito la non-linearità (es. legge di Horwitz)?"**

Analisi ispettiva: l'allerta RT-08. Per i lab chimici, esprimere UM in termini relativi su tutto il campo quando è stata determinata a un solo livello è un errore. Verifica anche che la stima sia significativa al livello dei **limiti di legge** del cliente (RT-08 esplicito).

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Stima UM da r/R del metodo: verifica prestazioni proprie, sanity check 2U>r, livello significativo per il cliente (limiti di legge), congruenza con letteratura. RT-08 § 7.6.3."*
- Urgenza: alta.
- Cosa NON può chiudere l'agente: NON può dichiarare l'UM conforme se 2U<r o se mancano le verifiche di propria prestazione.

### 7.6 trasversale — congruenza con risultati di PT/ILC

#### Domanda 1 — Meccanismo
**"Confrontate la vostra incertezza dichiarata con i dati di precisione ottenuti dai PT/ILC? Esistono casi in cui i vostri risultati sono "outlier" in un PT ma la vostra incertezza dichiarata copriva la differenza?"**

Analisi ispettiva: collegamento esplicito con § 7.7. Se l'UM "magicamente" copre sempre la differenza con il consensus PT, è sospetta di gonfiatura. Se invece c'è incongruenza ricorrente (outlier non spiegati), è sottostimata.

#### Domanda 2 — Estensione
**"Mostratemi i z-score (o En) dei vostri PT degli ultimi 12 mesi. Per i risultati borderline (|z|>2 o |En|>0.7), come avete reagito sulla stima UM?"**

Analisi ispettiva: cerco evidenza che PT borderline abbiano triggerato riesame UM.

#### Domanda 3 — Efficacia
**"Se domani un cliente vi facesse causa contestando l'incertezza dichiarata su un RdP perché ha portato a una falsa conformità, cosa producete come difesa?"**

Analisi ispettiva: domanda provocatoria, ma utile. La difesa è il dossier UM + dati PT + budget contributi.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Congruenza UM con risultati PT/ILC, riesame UM in caso di z-score borderline, dossier UM difendibile. ISO § 7.6 + § 7.7."*

## E. Errori fatali correlati

- **EF-7.6.1-A**: budget UM inesistente o ridotto alla sola ripetibilità intermedia (sR) — manca identificazione dei contributi (NC esplicita § 7.6.1).
- **EF-7.6.1-B**: campionamento accreditato che NON include UM da campionamento (collega RT-08 § 7.6.1).
- **EF-7.6.2-A**: tarature interne con UM stimata senza riferimento a JCGM-100 / EA-4/02 — NC esplicita RT-08.
- **EF-7.6.3-A**: `2U < r` — incertezza sottostimata, NC tecnica RT-08.
- **EF-7.6.3-B**: UM dichiarata in termini relativi su tutto il campo quando R del metodo è stata determinata solo a un livello (legge di Horwitz violata).
- **EF-7.6.3-C**: UM non valutata al livello del limite di legge per il cliente (rilevanza decisionale persa).
- **EF-7.6-D**: incongruenza ripetuta UM ↔ z-score PT senza riesame — NC § 7.7 con ricaduta su § 7.6.
- **EF-7.6-E**: stima UM "ferma" da anni nonostante cambi strumento/reagente/matrice — NC § 7.6 + § 7.2.2.2.

Vedi anche [[04_ERRORI_FATALI]] e [[A2_incertezza_misura]].

## F. Quando l'agente DEVE chiedere prima di rispondere

L'agente **non risponde** e **chiede al QM** se:

1. Si chiede quale approccio UM applicare a un nuovo metodo (bottom-up GUM vs top-down): scelta tecnica con impatti consistenti.
2. Si discute se un contributo (es. matrice, recupero, stabilità campione) sia "significativo" e quindi da includere: richiede valutazione tecnica.
3. La sanity check `2U > r` fallisce su un metodo: il QM deve decidere se rivalutare UM o se riconsiderare r/R del metodo.
4. I risultati PT mostrano z-score sistematicamente borderline: profilo § 7.7 + revisione UM.
5. Il cliente chiede UM espressa diversamente (intervallo di fiducia, deviazione standard, percentuale): verificare ammissibilità e coerenza § 7.8.
6. Si fa "rinuncia all'UM" per metodi consolidati (Nota 1 della norma): il caso va analizzato — non vale per tutto, solo per metodi consolidati e ben riconosciuti.
7. Cambio di metodo / strumento / matrice / range che potrebbe richiedere rivalutazione UM: § 7.6 + § 7.2.2.2.
8. Per regola decisionale (§ 7.1.3 + § 7.8.6) si sceglie guard band: dipende da UM e da rischio accettato dal cliente.

L'agente NON dichiara UM conforme senza vedere: (a) metodologia documentata, (b) budget contributi, (c) sanity check 2U>r per chimica, (d) verifica delle proprie prestazioni se si usano r/R del metodo, (e) congruenza coi PT.

## G. Ruolo dell'agente AI su questo requisito (governance) — 🟡 LIVELLO 3

> Strato generico; le specifiche del laboratorio vivono nella KB verticale. Vedi [[06_GOVERNANCE_AI]], [[03_PROTOCOLLO_HANDOFF_QM]].

**Cosa l'agente PUÒ fare su § 7.6:**
- Estrarre i budget di incertezza presenti nel SGQ e mapparli sui metodi accreditati; segnalare i metodi privi di budget.
- Eseguire il **sanity check `2U > r`** sui dati disponibili e segnalare i casi in cui fallisce (UM potenzialmente sottostimata).
- Confrontare l'UM dichiarata con i dati di precisione di validazione e con gli esiti PT, segnalando incongruenze (filo rotto).
- Segnalare stime di UM "ferme" nonostante cambi di strumento/reagente/matrice (collega § 7.2.2.2).

**Cosa l'agente NON PUÒ fare (handoff QM/RT):**
- Dichiarare un'UM "adeguata/conforme".
- Scegliere la metodologia (bottom-up GUM vs top-down) al posto del RT.
- Definire la guard band della regola decisionale (§ 7.1.3 / § 7.8.6).

**Privacy/riservatezza:** elaborazione in locale dei dossier UM, nessun dato fuori dal perimetro del laboratorio ([[sez_4_imparzialita_riservatezza]]).
