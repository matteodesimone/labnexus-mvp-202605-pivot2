---
codice: sez_7_4_manipolazione_oggetti
tipo: requisito_normativo
livello: misto
titolo: "§ 7.4 — Manipolazione degli oggetti di prova/taratura"
sezione_norma: "7.4"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.4" }
  - { doc: "RT-08", rev: "05", paragrafo: "7.4" }
fonti_secondarie:
  - { doc: "RT-08", rev: "05", paragrafo: "7.1.1 (tempi conservazione contrattuali)" }
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[02_TRE_LIVELLI_DI_CONTROLLO]]"
  - "[[sez_7_1_riesame_richieste]]"
  - "[[sez_7_3_campionamento]]"
  - "[[sez_7_5_registrazioni_tecniche]]"
tags: [manipolazione, oggetti_prova, identificazione, conservazione, controcampioni, segregazione]
escalation_qm: true
---

# § 7.4 — Manipolazione degli oggetti di prova/taratura

> Tutto ciò che accade tra il momento in cui il campione entra nel laboratorio e il momento in cui esce (o viene distrutto) ricade qui: identificazione univoca, registrazione condizioni anomale, conservazione, protezione dell'integrità, segregazione dei campioni non idonei. È il punto dove l'ispettore verifica se il lab tratta un campione come un dato tecnico tracciabile o come un oggetto generico.

## A. Testo della norma (🟢 LIVELLO 1 — fatto normativo)

> **Fonte**: `[ISO-17025-2018 § 7.4.1 - 7.4.4]` — testo a pagamento, non riprodotto. Riassunto operativo:

- **7.4.1** Il laboratorio deve avere una **procedura per trasporto, ricezione, manipolazione, protezione, immagazzinamento, conservazione e disposizione/restituzione** degli oggetti di prova/taratura, includendo la protezione dell'integrità e degli interessi del lab e del cliente. Precauzioni per evitare deterioramento, contaminazione, perdita o danno durante manipolazione, trasporto, immagazzinamento/attesa, preparazione. Le istruzioni di manipolazione fornite dall'oggetto devono essere seguite.
- **7.4.2** Sistema di **identificazione univoca** dei campioni. L'identificazione deve essere mantenuta fintanto che l'oggetto è sotto la responsabilità del laboratorio, in modo che non possa esserci confusione fisica o quando ci si riferisce all'oggetto in registrazioni o documenti. Il sistema deve, ove appropriato, gestire una suddivisione in sub-campioni e il trasferimento dell'oggetto.
- **7.4.3** Alla ricezione, eventuali **deviazioni dalle condizioni specificate** devono essere registrate. Quando esistono dubbi sull'idoneità dell'oggetto per la prova/taratura, o quando un oggetto non è conforme alla descrizione fornita, il laboratorio deve consultare il cliente per istruzioni ulteriori prima di procedere, e deve registrare l'esito della consultazione. Quando il cliente richiede attività di prova su un oggetto riconoscendo deviazioni, il laboratorio deve includere disclaimer sui risultati che possono essere influenzati.
- **7.4.4** Quando gli oggetti devono essere conservati o condizionati in condizioni ambientali specifiche, queste devono essere mantenute, monitorate e registrate.

## B. Prescrizioni Accredia (🟢 LIVELLO 2 — prescrizione Accredia)

> **Fonte**: `[RT-08 rev.05 § 7.4]`

### 7.4.1 — Conservazione campioni e controcampioni

> **Citazione integrale RT-08**:
> *"Il Laboratorio deve conservare i campioni sui quali ha eseguito la prova e/o eventuali controcampioni e le registrazioni tecniche relative alle prove effettuate, per il tempo concordato con il richiedente, o stabilito dalle pertinenti leggi o disposto dall'Autorità Giudiziaria (cfr. 7.1.1.)."*

Collegamento esplicito con § 7.1.1: i tempi di conservazione vanno **contrattualizzati** o derivare da prescrizioni cogenti / Autorità Giudiziaria.

### 7.4.4 — Aree di segregazione

> **Citazione integrale RT-08**:
> *"Il Laboratorio deve disporre di adeguate aree di segregazione per conservare i campioni non idonei."*

Punti 7.4.2 e 7.4.3: "Si applica il requisito di norma" senza integrazioni Accredia.

## C. Documenti applicabili (LS-04 rev.20)

- `ISO-17025-2018` § 7.4 — norma base.
- `RT-08 rev.05` § 7.4 e § 7.1.1 (tempi di conservazione contrattuali).
- Eventuali norme settoriali con vincoli specifici (es. campioni biologici → catena del freddo; amianto → segregazione fisica; sostanze infiammabili → stoccaggio normato).

## D. Domande ispettive (🟡 LIVELLO 3 — prassi ispettiva)

> Pattern: meccanismo → estensione → efficacia.

### 7.4.1 — Procedura di manipolazione e conservazione

#### Domanda 1 — Meccanismo
**"Mostratemi la procedura di ricezione e manipolazione campioni. Come distinguete campioni di prova, controcampioni, campioni in attesa, campioni esauriti e campioni non idonei? Quali tempi di conservazione applicate e come li avete derivati dal contratto?"**

Analisi ispettiva: la procedura deve coprire l'intero ciclo. Tempi di conservazione coerenti con quanto contrattualizzato (vincolo RT-08 § 7.1.1 + 7.4.1). Tipico problema: la procedura dice "30 giorni dopo emissione RdP" ma alcuni contratti prevedono di più o di meno, e nessuno verifica per ogni commessa.

#### Domanda 2 — Estensione
**"Mostrate i campioni in conservazione: per 3 di essi, voglio vedere il legame con il RdP, il termine di conservazione, e che il termine sia coerente col contratto."**

Analisi ispettiva: prendo 3 campioni a caso (preferibilmente di età diverse: appena ricevuto, in lavorazione, post-RdP). Verifico:
- Identificazione univoca (collega § 7.4.2);
- Tracciabilità al RdP/commessa;
- Termine di conservazione coerente col contratto;
- Condizioni ambientali rispettate (T, umidità, luce — collega § 7.4.4 e § 6.3).

#### Domanda 3 — Efficacia
**"Cosa fate quando un cliente non passa a ritirare il controcampione entro il termine? Cosa fate per i campioni che superano i limiti contrattuali? Mostratemi la registrazione dello smaltimento."**

Analisi ispettiva: la disposizione finale è parte del § 7.4.1. Se la procedura dice "smaltiamo" ma non c'è registrazione, c'è buco. Se i campioni sono accumulati oltre il termine senza decisione documentata, c'è un problema di controllo.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Procedura di manipolazione + tempi conservazione coerenti con contratto + tracciabilità smaltimento. ISO § 7.4.1, RT-08 § 7.4.1 e § 7.1.1."*
- Urgenza: media — è igiene del processo.
- Cosa NON può chiudere l'agente: dichiarare conforme la conservazione senza vedere coerenza contratto ↔ procedura ↔ realtà.

### 7.4.2 — Identificazione univoca

#### Domanda 1 — Meccanismo
**"Come identificate univocamente i campioni? Etichetta, codice a barre, codifica nel LIMS? Come gestite la suddivisione in sub-campioni e il trasferimento da un'area all'altra del laboratorio?"**

Analisi ispettiva: l'identificazione deve essere unica, persistente e gestire i sub-campioni. Tipico problema: l'etichetta originale cade durante la preparazione e il sub-campione perde tracciabilità. Devono esserci pratiche di pre-etichettatura dei sub-contenitori.

#### Domanda 2 — Estensione
**"Per un campione lavorato la scorsa settimana: ricostruite la traccia dal ricevimento al RdP, includendo TUTTI i sub-campioni e i trasferimenti tra aree."**

Analisi ispettiva: test di ricostruzione tecnica. Devo poter risalire ad ogni movimento. Collega § 7.5 (registrazioni tecniche) — se la traccia non è ricostruibile, è NC del § 7.5 + 7.4.2.

#### Domanda 3 — Efficacia
**"È mai successo che ci sia stato un dubbio sull'identità di un campione (etichetta cancellata, scambio sospetto)? Come è stato gestito?"**

Analisi ispettiva: gli scambi/dubbi capitano, ma devono essere gestiti come **attività non conformi** (§ 7.10). Se la risposta è "non capita mai", è poco credibile su grandi volumi: posso chiedere evidenza del trattamento NC.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Sistema identificazione univoca con gestione sub-campioni, test ricostruzione tecnica su un campione recente, gestione dubbi di identificazione come NC. ISO § 7.4.2 + § 7.10."*

### 7.4.3 — Condizioni di ricezione e idoneità

#### Domanda 1 — Meccanismo
**"Quali criteri di accettabilità applicate al ricevimento? Per quali parametri (temperatura, integrità contenitore, quantità, tempo dal campionamento) avete soglie definite? Come tracciate eventuali deviazioni?"**

Analisi ispettiva: la procedura deve avere criteri oggettivi, non solo "valutazione visiva". Per molti metodi cogenti i criteri sono nel metodo stesso (es. T al ricevimento per microbiologia alimenti).

#### Domanda 2 — Estensione
**"Mostratemi 3 casi in cui sono state registrate deviazioni al ricevimento. Cosa avete fatto: consultazione del cliente, accettazione con disclaimer sul RdP, rifiuto?"**

Analisi ispettiva: il § 7.4.3 dice che in caso di dubbio sull'idoneità il lab DEVE consultare il cliente PRIMA di procedere. Se vediamo deviazioni accettate senza traccia di consultazione, NC.

#### Domanda 3 — Efficacia
**"Per un caso in cui avete proceduto con campione non perfettamente conforme alla descrizione fornita, mostrate il RdP. C'è il disclaimer richiesto dal § 7.4.3?"**

Analisi ispettiva: il § 7.4.3 chiede esplicitamente il disclaimer sui risultati che possono essere influenzati. Manca → NC.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Criteri accettabilità al ricevimento + tracciabilità deviazioni + consultazione cliente + disclaimer sul RdP. ISO § 7.4.3."*

### 7.4.4 — Condizioni ambientali di conservazione

#### Domanda 1 — Meccanismo
**"Per i campioni che richiedono condizioni controllate (T, umidità, luce, atmosfera): quali aree avete, come monitorate, con quale frequenza, e come registrate? Avete aree di segregazione per i campioni non idonei come richiesto da RT-08 § 7.4.4?"**

Analisi ispettiva: monitoraggio continuo (sonda registratore) o spot check? Per campioni alimentari refrigerati serve monitoraggio continuo con allarmi. Le aree di segregazione devono essere fisicamente distinte (RT-08 esplicito).

#### Domanda 2 — Estensione
**"Mostratemi le registrazioni delle condizioni ambientali nelle aree di conservazione per l'ultimo mese. Ci sono escursioni fuori soglia? Come le avete gestite?"**

Analisi ispettiva: cerco escursioni. Se ce ne sono, deve esserci una valutazione di impatto sui campioni conservati (potenziale NC § 7.10 con ricaduta sui RdP emessi).

#### Domanda 3 — Efficacia
**"Se la cella frigo va in blocco di notte e la T sale, come ve ne accorgete e cosa fate?"**

Analisi ispettiva: sistemi di allarme? Procedura emergenziale? Chi è responsabile? Tracciato il follow-up?

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Aree di conservazione + monitoraggio condizioni + aree segregazione campioni non idonei + gestione escursioni. ISO § 7.4.4 + RT-08 § 7.4.4."*

## E. Errori fatali correlati

- **EF-7.4.1-A**: tempi di conservazione campioni non coerenti col contratto (collega § 7.1.1) — NC.
- **EF-7.4.1-B**: smaltimento campioni senza registrazione — NC.
- **EF-7.4.2-A**: identificazione univoca non mantenuta nei sub-campioni — NC, e impatta § 7.5.
- **EF-7.4.2-B**: ricostruzione tecnica di un campione non riproducibile — NC grave su § 7.5 + 7.4.2.
- **EF-7.4.3-A**: deviazione al ricevimento accettata senza consultazione cliente — NC esplicita.
- **EF-7.4.3-B**: campione accettato con deviazione, RdP privo del disclaimer obbligatorio — NC.
- **EF-7.4.4-A**: mancanza di aree fisicamente segregate per campioni non idonei — NC esplicita RT-08.
- **EF-7.4.4-B**: escursioni di T in conservazione non valutate per impatto sui RdP già emessi — NC § 7.10 a cascata.

Vedi anche [[04_ERRORI_FATALI]].

## F. Quando l'agente DEVE chiedere prima di rispondere

L'agente **non risponde** e **chiede al QM** se:

1. Si chiede se un campione "borderline" (es. T ricevimento +1°C oltre soglia) sia accettabile: serve valutazione tecnica caso per caso.
2. Si chiede di smaltire campioni di un cliente che non ha pagato o è irreperibile prima del termine contrattuale: profilo legale + qualità.
3. Si chiede se i campioni di un metodo che richiede conservazione molto specifica (es. campioni per amianto, per sostanze volatili) sono conservati correttamente: verifica tecnica del metodo.
4. Una NC sulle condizioni ambientali di conservazione (es. cella sotto soglia per 2 ore) impatta sui RdP già emessi nei giorni precedenti: profilo § 7.10 (gestione attività non conformi).
5. Un cliente chiede tempi di conservazione fuori standard del laboratorio (molto lunghi o molto brevi): impatto su capacità fisica e contratto.
6. Un campione viene presentato per analisi accreditata ma è già stato manipolato dal cliente in modo non documentato (es. concentrato/diluito senza tracciabilità): collide con § 7.4.3 + § 7.2.1.6.

L'agente NON dichiara conforme una manipolazione senza vedere coerenza tra procedura, contratto, condizioni ambientali e registrazioni reali. NON accetta deviazioni al ricevimento senza traccia di consultazione cliente e disclaimer.
