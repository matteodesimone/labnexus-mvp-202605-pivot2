---
codice: sez_7_3_campionamento
tipo: requisito_normativo
livello: misto
titolo: "§ 7.3 — Campionamento"
sezione_norma: "7.3"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.3" }
  - { doc: "RT-08", rev: "05", paragrafo: "7.3" }
fonti_secondarie:
  - { doc: "RT-23", rev: "05", paragrafo: "combinazione metodi campionamento + determinazione" }
  - { doc: "EURACHEM-Sampling", rev: "2nd ed. 2019", paragrafo: "incertezza da campionamento" }
  - { doc: "ISTISAN-22/39", rev: "2023", paragrafo: "incertezza da campionamento (trad. it.)" }
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[02_TRE_LIVELLI_DI_CONTROLLO]]"
  - "[[sez_7_1_riesame_richieste]]"
  - "[[sez_7_4_manipolazione_oggetti]]"
  - "[[sez_7_6_incertezza]]"
  - "[[sez_7_8_presentazione_risultati]]"
  - "[[A7_campionamento]]"
tags: [campionamento, accreditamento_campionamento, ricevimento, subcampionamento]
escalation_qm: true
---

# § 7.3 — Campionamento

> Il § 7.3 si applica **solo se il laboratorio fa anche il campionamento**. Se il campionamento è del cliente o di terzi, il laboratorio gestisce solo il **ricevimento** del campione (oggetto del § 7.4). Il campionamento è accreditabile **solo se associato a una prova successiva accreditata** (RT-08 esplicito). Confondere "il laboratorio prende il campione" con "il laboratorio è accreditato per il campionamento" è uno degli equivoci più frequenti.

## A. Testo della norma (🟢 LIVELLO 1 — fatto normativo)

> **Fonte**: `[ISO-17025-2018 § 7.3.1 - 7.3.3]` — testo a pagamento, non riprodotto. Riassunto operativo:

- **7.3.1** Quando il laboratorio esegue il campionamento, deve disporre di un **piano e metodo di campionamento**. Il metodo affronta i fattori da controllare per assicurare la validità dei risultati. Il piano si basa su metodi statistici appropriati. Il campionamento descrive selezione, piano, preparazione/trattamento, secondo lo scopo della prova/taratura.
- **7.3.2** Il metodo di campionamento deve descrivere:
  - **a)** selezione di campioni o siti;
  - **b)** piano di campionamento;
  - **c)** preparazione/trattamento del campione per produrre l'oggetto richiesto per prova/taratura successiva.
- **7.3.3** Il laboratorio deve conservare le **registrazioni** del campionamento, includendo (se rilevante): riferimento al metodo di campionamento usato; data e ora; identificazione/descrizione del campione; identificazione del campionatore; identificazione apparecchiatura usata; condizioni ambientali o di trasporto; diagrammi o equivalenti per identificare il punto di campionamento; deviazioni, aggiunte o esclusioni dal metodo di campionamento.

## B. Prescrizioni Accredia (🟢 LIVELLO 2 — prescrizione Accredia)

> **Fonte**: `[RT-08 rev.05 § 7.3]`

### 7.3.1 — Quando il campionamento è accreditabile

> **Citazione integrale RT-08**:
> *"Il campionamento è accreditabile solo se associato ad una successiva prova accreditata."*

### 7.3.1 — Due scenari operativi

**Scenario 1**: il metodo di campionamento è specificato in un metodo normalizzato/ufficiale **diverso** da quello della determinazione (es. ISO 18593, UNI EN 1948-1).

Condizioni da soddisfare:
- la successiva prova associata deve essere eseguita sotto accreditamento ISO/IEC 17025, dal laboratorio stesso o da altro laboratorio accreditato per la specifica prova;
- il metodo della determinazione deve essere scelto rispettando le indicazioni del metodo di campionamento e in ogni caso quanto previsto dal **Regolamento RT-23** sulla combinazione dei metodi;
- il laboratorio si impegna a concordare con il cliente, in fase contrattuale, le prove successive associate al campionamento e a rispettare quanto previsto al § 7.1.1; per il punto b) il laboratorio esterno deve essere accreditato secondo ISO/IEC 17025 per la specifica prova.

**Scenario 2**: un singolo metodo comprende sia campionamento sia determinazione.

Il laboratorio può:
- richiedere l'accreditamento dell'intero metodo (campionamento + determinazione);
- richiedere l'accreditamento della sola determinazione analitica, esplicitando l'**esclusione del capitolo relativo al campionamento**.

### 7.3.1 — Quando il cliente fa il campionamento

> **Citazione integrale RT-08**:
> *"È opportuno che il Laboratorio fornisca al cliente, qualora questi effettui il campionamento, idonea assistenza (istruzioni per il campionamento e relative registrazioni, contenitori, ecc.)."*

Quando il campionamento è del cliente, sul RdP il laboratorio NON dichiara l'accreditamento del campionamento e si applicano le prescrizioni del § 7.8.2.2 (vedi [[sez_7_8_presentazione_risultati]]).

**Punti 7.3.2 e 7.3.3**: "Si applica il requisito di norma" senza integrazioni Accredia.

## C. Documenti applicabili (LS-04 rev.20)

- `ISO-17025-2018` § 7.3 — norma base.
- `RT-08 rev.05` § 7.3 — prescrizioni Accredia.
- `RT-23 rev.05` — combinazione di metodi (vincolante quando campionamento e determinazione sono in metodi distinti).
- `EURACHEM-Sampling` 2nd ed. 2019 — incertezza da campionamento.
- `ISTISAN-22/39` 2023 — traduzione italiana EURACHEM Sampling.
- Norme/metodi settoriali di campionamento (es. ISO 18593 per superfici alimentari, UNI EN 1948-1 per emissioni in atmosfera).

## D. Domande ispettive (🟡 LIVELLO 3 — prassi ispettiva)

> Pattern: meccanismo → estensione → efficacia.
> Distinzione preliminare: **il laboratorio è accreditato per il campionamento?** Se no, salto al § 7.4 (ricevimento campione). Se sì, applico le domande seguenti.

### 7.3.1 — Scope del campionamento

#### Domanda 1 — Meccanismo
**"Per quali prove siete accreditati anche per il campionamento? Mostrate l'elenco prove accreditate e indicate dove è esplicitato che il campionamento è incluso."**

Analisi ispettiva: leggo l'elenco prove (Accredia pubblica). Cerco l'esplicita indicazione del campionamento accreditato. Se il lab dice "facciamo campionamento per tutto" ma l'elenco non lo dice, c'è disallineamento: il campionamento esiste come servizio commerciale ma non è accreditato.

#### Domanda 2 — Estensione
**"Per uno dei campionamenti accreditati: mostrate il metodo/piano di campionamento, l'analisi successiva (deve essere accreditata) e il rispetto del Regolamento RT-23 per la combinazione."**

Analisi ispettiva: i 3 vincoli RT-08 dello Scenario 1. Tipico problema: il laboratorio fa campionamento accreditato ma poi subappalta la determinazione a un lab non accreditato per quella specifica prova — NC grave.

#### Domanda 3 — Efficacia
**"Avete esempi recenti di campionamento accreditato in cui non avete potuto rispettare il piano di campionamento? Come avete gestito la deviazione?"**

Analisi ispettiva: il § 7.3.3 chiede di registrare deviazioni/aggiunte/esclusioni. Se la risposta è "non capita mai" ma vediamo verbali con scostamenti, c'è un buco di tracciabilità.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Scope del campionamento accreditato: rispetto Scenario 1/2 RT-08 § 7.3.1, applicazione RT-23 sulla combinazione metodi, registrazione deviazioni. Voglio vedere il fascicolo di un campionamento accreditato recente."*
- Urgenza: alta se il lab si presenta come "laboratorio di campionamento".
- Cosa NON può chiudere l'agente: dichiarare conforme un campionamento senza verificare che la prova successiva sia accreditata.

### 7.3.2 — Metodo di campionamento

#### Domanda 1 — Meccanismo
**"Per un campionamento sotto accreditamento, fatemi vedere il metodo/piano: selezione punti, piano statistico, preparazione/trattamento del campione."**

Analisi ispettiva: i tre punti a)/b)/c). Il piano deve essere basato su metodi statistici appropriati (norma esplicita), non solo "prendiamo 3 punti".

#### Domanda 2 — Estensione
**"Come gestite il subcampionamento (riduzione di un campione primario a campione di laboratorio)? È compreso nella validazione del metodo di campionamento?"**

Analisi ispettiva: il subcampionamento è una fase critica spesso dimenticata. Se il metodo di campionamento (es. ISO 18593) la prevede, deve essere coperta. Se è metodo sviluppato dal laboratorio che include campionamento, anche la validazione include questa fase (vedi RT-08 § 7.2.2.1.c).

#### Domanda 3 — Efficacia
**"Quando le condizioni ambientali in campo sono fuori specifica (es. pioggia improvvisa durante campionamento aria), cosa fate?"**

Analisi ispettiva: criterio decisionale documentato? Annullamento del campionamento, riprogrammazione, segnalazione al cliente.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Metodo di campionamento: piano statistico, subcampionamento, criteri di accettazione condizioni in campo. RT-08 § 7.3 e Regolamento RT-23."*

### 7.3.3 — Registrazioni del campionamento

#### Domanda 1 — Meccanismo
**"Mostrate il modulo di verbale di campionamento. Coprite tutti gli elementi del § 7.3.3 (data, ora, campionatore, apparecchiatura, condizioni, diagramma punto, deviazioni)?"**

Analisi ispettiva: la lista è esaustiva. Tipici buchi: manca il diagramma del punto di campionamento (rilevante per aria/acque/superfici), o manca l'identificazione dell'apparecchiatura usata (con tracciabilità taratura — collega § 6.4).

#### Domanda 2 — Estensione
**"Estraete 5 verbali di campionamento recenti. Voglio verificare la completezza, l'identificazione del campionatore e la coerenza con il metodo."**

Analisi ispettiva: campionamento del verbale reale.

#### Domanda 3 — Efficacia
**"Quando una deviazione dal piano avviene (es. punto non accessibile), come la riportate sul verbale e sul RdP?"**

Analisi ispettiva: deve essere visibile sia sul verbale sia sul RdP. Manca uno dei due → NC.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Verbali di campionamento: completezza per § 7.3.3, gestione deviazioni e propagazione al RdP."*

### Per laboratori che NON sono accreditati per il campionamento

#### Domanda 1 — Meccanismo
**"Il vostro accreditamento NON include il campionamento. Quando ricevete un campione dal cliente, come trattate le informazioni sul campionamento (chi, dove, come, condizioni)?"**

Analisi ispettiva: si passa a [[sez_7_4_manipolazione_oggetti]] e § 7.8.2.2 (presentazione risultati). Il RdP deve dichiarare chi ha campionato (cliente / altro) e che il campionamento NON è accreditato. Le misure effettuate in fase di campionamento (volume, area, portata) NON possono entrare nelle unità di misura sul RdP del lab se non esplicitate (vedi sez_7_8 per i casi accettabili).

#### Domanda 2 — Estensione
**"Fornite ai vostri clienti istruzioni per il campionamento o contenitori dedicati? Come lo tracciate?"**

Analisi ispettiva: RT-08 lo raccomanda ("idonea assistenza"). È buona prassi ma non vincolante.

#### Domanda 3 — Efficacia
**"Se ricevete un campione che mostra segni di campionamento improprio (contenitore sbagliato, tempo di trasporto eccessivo, T fuori specifica), cosa fate?"**

Analisi ispettiva: criteri di accettabilità del campione (collega § 7.4.3 sulle condizioni). Documentati? Comunicati al cliente?

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Pur non essendo accreditati per il campionamento, c'è assistenza al cliente e criteri di accettabilità del campione ricevuto? § 7.4.3 + § 7.8.2.2."*

## E. Errori fatali correlati

- **EF-7.3-A**: presentarsi come "laboratorio di campionamento" senza che il campionamento sia nello scope di accreditamento — uso improprio del marchio (collega RG-09).
- **EF-7.3-B**: campionamento accreditato seguito da determinazione **non accreditata** o **subappaltata a non accreditato** — viola RT-08 § 7.3.1.
- **EF-7.3-C**: campionamento accreditato senza rispetto del Regolamento RT-23 sulla combinazione metodi.
- **EF-7.3-D**: verbale di campionamento incompleto rispetto al § 7.3.3 (manca diagramma punto, apparecchiatura, ecc.).
- **EF-7.3-E**: deviazione dal piano di campionamento non tracciata sul verbale e/o non propagata sul RdP.
- **EF-7.3-F**: subcampionamento (riduzione campione primario) non incluso nella validazione del metodo (per metodi sviluppati che includono campionamento).
- **EF-7.3-G**: per lab non accreditato per campionamento, RdP che riporta unità di misura derivate dal campionamento (es. mg/m³) senza tracciabilità delle misure effettuate dal cliente (§ 7.8.2.2).

Vedi anche [[04_ERRORI_FATALI]].

## F. Quando l'agente DEVE chiedere prima di rispondere

L'agente **non risponde** e **chiede al QM** se:

1. Si chiede se una specifica attività di prelievo del lab è "campionamento accreditato" o "ricevimento campione": dipende dallo scope dell'elenco prove pubblicato da Accredia.
2. Si valuta se applicare lo Scenario 1 o 2 RT-08 quando il metodo include campionamento ma il lab non vuole accreditare quel capitolo.
3. Si chiede se la combinazione "metodo campionamento X + metodo determinazione Y" è ammessa: richiede lettura RT-23.
4. Subcampionamento o trattamento preliminare (omogeneizzazione, riduzione, conservazione): collide con incertezza da campionamento (EURACHEM-Sampling).
5. Campionamento accreditato in subappalto: vale solo se il subcontraente è accreditato per quella specifica attività di campionamento.
6. Le condizioni in campo escono dalle specifiche del metodo: il QM decide annullamento vs prosecuzione con deviazione.
7. Un cliente chiede di unire risultati di campionamento accreditato + determinazione di altro lab (interlaboratorio): collide con § 7.1.1 punti subappalto e con § 7.8.2.1.

L'agente NON dichiara mai un campionamento conforme senza verificare: (a) scope accreditamento, (b) rispetto Scenario 1/2 RT-08, (c) coerenza con RT-23, (d) completezza verbale § 7.3.3.
