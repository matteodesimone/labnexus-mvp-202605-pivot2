---
codice: sez_6_6_prodotti_servizi_esterni
tipo: requisito_normativo
livello: misto
titolo: "§ 6.6 — Prodotti e servizi forniti dall'esterno"
sezione_norma: "6.6"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "6.6" }
  - { doc: "RT-08", rev: "05", paragrafo: "6.6" }
fonti_secondarie:
  - { doc: "LS-04", rev: "20" }
  - { doc: "RG-02", rev: "08" }
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[02_TRE_LIVELLI_DI_CONTROLLO]]"
  - "[[04_ERRORI_FATALI]]"
  - "[[sez_6_4_dotazioni]]"
  - "[[sez_6_5_riferibilita]]"
  - "[[sez_7_2_selezione_verifica_validazione_metodi]]"
  - "[[sez_7_7_assicurazione_validita]]"
tags: [fornitori, qualifica_fornitori, subappalto, acquisti, prodotti_critici, servizi_esterni]
escalation_qm: true
---

# § 6.6 — Prodotti e servizi forniti dall'esterno

> Tutti i prodotti e servizi che influenzano le attività di laboratorio devono essere idonei. Il laboratorio definisce procedure per: scelta dei fornitori, qualifica, monitoraggio prestazioni, riesame periodico. Comunica ai fornitori, in modo non ambiguo, i requisiti. La sezione si applica a reagenti, CRM, consumabili, servizi di taratura, servizi di subappalto di prove, servizi di manutenzione, servizi IT che impattano sui dati di laboratorio.

## A. Testo della norma (🟢 LIVELLO 1 — fatto normativo)
> **Fonte**: `[ISO-17025-2018 § 6.6]` — testo a pagamento, non riprodotto integralmente.

- **6.6.1 Idoneità prodotti/servizi forniti dall'esterno** — Il laboratorio deve garantire che siano idonei solo i prodotti e servizi forniti dall'esterno che incorporano nelle proprie attività, quelli che sono utilizzati direttamente nelle attività di laboratorio o per supportare il funzionamento del laboratorio.
- **6.6.2 Procedure di scelta, valutazione e monitoraggio dei fornitori** — Il laboratorio deve disporre di procedura/e e conservare registrazioni per:
  - a) definizione, riesame ed approvazione dei requisiti dei prodotti e servizi forniti dall'esterno;
  - b) definizione dei criteri per la valutazione, scelta, monitoraggio delle prestazioni e ri-valutazione dei fornitori esterni;
  - c) assicurazione che i prodotti e servizi forniti dall'esterno siano conformi ai requisiti stabiliti dal laboratorio o, se applicabile, ai requisiti della norma, prima del loro utilizzo o consegna al cliente;
  - d) attuazione di qualsiasi azione derivante dalla valutazione, monitoraggio delle prestazioni e ri-valutazione dei fornitori esterni.
- **6.6.3 Comunicazione requisiti ai fornitori** — Il laboratorio comunica ai fornitori esterni i propri requisiti per: i prodotti e i servizi da fornire; i criteri di accettazione; la competenza, incluso ogni requisito di qualifica del personale; le attività che il laboratorio (o il cliente) intende svolgere presso il fornitore esterno.

## B. Prescrizioni Accredia (🟢 LIVELLO 2 — prescrizione Accredia)
> **Fonte**: `[RT-08 rev.05 § 6.6 — EC 10-02-2022]`

- **6.6.1** — Si applica il requisito di norma.
- **6.6.2** — Si applica il requisito di norma.
- **6.6.3** — Si applica il requisito di norma.

Nota: RT-08 § 6.6 non aggiunge integrazioni testuali. Le prescrizioni cogenti restano quelle ISO 17025 § 6.6. ATTENZIONE: la valutazione del fornitore di CRM/tarature in casi residuali è esplicitamente richiamata da RT-08 § 6.5 (B.3 chiusura e casi 3a/3b), che fa riferimento al § 6.6 come strumento applicativo.

## C. Documenti applicabili (LS-04 rev.20)
- `ISO-17025-2018 § 6.6`
- `RT-08 rev.05 § 6.6` (richiama anche § 6.5 per CRM/tarature)
- `RG-02 rev.08` — Regolamento accreditamento laboratori di prova (per perimetro accreditamento e gestione subappalto se previsto)
- `Racc-CdIG` — Raccomandazioni del Comitato di Indirizzo e Garanzia (criteri omogenei sulla valutazione fornitori)
- Norme tecniche di settore per requisiti specifici di prodotti/servizi (es. EN, ISO settoriali sui reagenti)

## D. Domande ispettive (🟡 LIVELLO 3 — prassi ispettiva)
> Pattern: meccanismo → estensione → efficacia. Sul § 6.6 l'ispettore cerca un sistema, non un elenco di fornitori. La presenza di "albo fornitori" non equivale a "sistema di qualifica".

### 6.6.1 — Idoneità di prodotti e servizi esterni
#### Domanda 1 — Meccanismo
**"Avete classificato quali prodotti e servizi esterni influenzano le attività di laboratorio o ne supportano il funzionamento? Mi mostri questa classificazione."**
Analisi ispettiva: cerco un documento che mappi le categorie: reagenti analitici, CRM, consumabili (vetreria/colonne/filtri), servizi di taratura, servizi di subappalto di prove, manutenzione strumentale, validazione software, servizi IT/cloud che gestiscono dati di laboratorio, gas tecnici, materiali di calibrazione strumentale. Se l'elenco è "scarso" o riguarda solo "i reagenti chimici", probabilmente i servizi (taratura, manutenzione, IT) sono fuori controllo § 6.6.

#### Domanda 2 — Estensione
**"Il sistema copre i prodotti incorporati nei risultati (es. reagenti, CRM)? Quelli usati direttamente (es. servizi di taratura, manutenzione)? Quelli di supporto (es. cloud, software)?"**
Analisi ispettiva: ISO 17025 § 6.6.1 distingue tre famiglie. Spesso si copre la prima e si trascurano le altre due.

#### Domanda 3 — Efficacia
**"Quando avete introdotto l'ultimo nuovo prodotto/servizio (es. nuovo reagente, nuovo lab di taratura, nuovo SW), come l'avete inserito nel sistema § 6.6 PRIMA dell'uso?"**
Analisi ispettiva: cerco il filtro pre-uso. Se l'uso parte e la qualifica arriva dopo, manca il controllo preventivo.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: "Acquisire mappatura prodotti/servizi esterni con classificazione di criticità per le attività di laboratorio."
- Riferimento: ISO § 6.6.1 + RT-08 § 6.6.1
- Urgenza: medio-alto
- Cosa NON può chiudere agente: la classificazione di criticità è giudizio del QM/RT.

### 6.6.2 — Procedure per scelta, valutazione, monitoraggio e ri-valutazione
#### Domanda 1 — Meccanismo (criteri di scelta)
**"Mi mostri i criteri documentati per la scelta di un nuovo fornitore critico. Cosa valutate? Come decidete in caso di alternative?"**
Analisi ispettiva: criteri tipici — accreditamento (per servizi di taratura: rinvio a RT-08 § 6.5), ISO 17034 (per CRM), referenze, qualità documentale del certificato di prodotto, conformità a norma di prodotto (es. ACS, pharma grade), tempi di consegna, lotto/scadenza. Se i criteri si esauriscono in "prezzo + tempi", manca la qualifica tecnica.

#### Domanda 2 — Meccanismo (monitoraggio prestazioni)
**"Come monitorate nel tempo le prestazioni dei fornitori critici? Voglio vedere indicatori: rese, NC accettazione, problemi di qualità prodotto, esiti di tarature, ritardi."**
Analisi ispettiva: cerco un cruscotto, anche minimale. Un sistema dove "i fornitori si valutano una volta all'anno" senza dati in input è retorica. Cerco trigger: NC accettazione lotto, NC interna riconducibile al prodotto, esito taratura fuori atteso → ricaduta sulla valutazione fornitore.

#### Domanda 3 — Estensione (ri-valutazione)
**"Mi mostri l'ultima ri-valutazione dei fornitori. Quando, su quali dati, con quale esito (conferma, declassamento, sospensione, sostituzione)?"**
Analisi ispettiva: la ri-valutazione deve avere periodicità definita e produrre un esito documentato per ciascun fornitore critico. La ri-valutazione "tacita" (= se non escludo, confermo) non è ri-valutazione.

#### Domanda 4 — Meccanismo (conformità pre-uso)
**"Per il prodotto X consegnato dal fornitore Y, mi mostri la verifica di conformità ai requisiti PRIMA dell'utilizzo: chi controlla, cosa controlla (certificato, lotto, integrità), criterio di accettazione, evidenza."**
Analisi ispettiva: § 6.6.2 c) chiede esplicitamente conformità verificata prima dell'utilizzo. La ricezione "passiva" del materiale senza check è NC.

#### Domanda 5 — Efficacia
**"Per l'ultima NC su un fornitore critico, mi mostri l'azione derivante (§ 6.6.2 d): comunicazione al fornitore, azione interna, eventuale impatto sui risultati già prodotti."**
Analisi ispettiva: § 6.6.2 d) impone azione conseguente. Se la NC su fornitore non genera azione, è una NC documentata ma non gestita.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: "Acquisire procedura qualifica fornitori; verificare a campione: criteri di scelta, monitoraggio attivo, ri-valutazione periodica, conformità pre-uso, gestione azioni post-NC."
- Riferimento: ISO § 6.6.2 + RT-08 § 6.6.2
- Urgenza: alto
- Cosa NON può chiudere agente: la valutazione di adeguatezza dei criteri tecnici di qualifica per categoria specifica richiede QM/RT.

### 6.6.3 — Comunicazione dei requisiti ai fornitori
#### Domanda 1 — Meccanismo
**"Quando ordinate un servizio di taratura o un CRM, in che forma comunicate i vostri requisiti? Mi mostri l'ultimo ordine: parametro/grandezza, range, incertezza massima ammessa, criteri di accettazione, qualifica del personale richiesta."**
Analisi ispettiva: § 6.6.3 elenca quattro aree. Ordini "lapidari" del tipo "taratura bilancia, fattura come al solito" non comunicano requisiti. Pretendo che l'ordine (o un contratto quadro a monte) descriva: grandezza, range, incertezza target, criteri di accettazione, eventuali norme tecniche di riferimento, qualifica del personale (se rilevante).

#### Domanda 2 — Estensione (subappalto)
**"Se subappaltate prove a un altro laboratorio, come comunicate i requisiti? Avete contratto quadro? Avete clausole sul fatto che il fornitore deve permettere accesso al laboratorio (vostro e/o del cliente, § 6.6.3 d)?"**
Analisi ispettiva: per il subappalto, le quattro aree di § 6.6.3 sono tutte rilevanti, in particolare la "d" su attività che il lab o il cliente intende svolgere presso il fornitore (audit di seconda parte). Senza, la qualifica del subappaltatore è opaca.

#### Domanda 3 — Efficacia
**"Per l'ultimo certificato di taratura ricevuto, l'incertezza dichiarata era compatibile con i requisiti che avevate comunicato? Cosa succede se non lo è?"**
Analisi ispettiva: la coerenza requisito comunicato ↔ requisito ricevuto è prova di efficacia. Se si accetta sistematicamente "qualsiasi incertezza il lab esterno dichiari", la comunicazione è stata cosmetica.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: "Verificare evidenza di comunicazione completa dei requisiti ai fornitori (ordini, contratti quadro) e coerenza tra requisito richiesto e prodotto/servizio ricevuto."
- Riferimento: ISO § 6.6.3 + RT-08 § 6.6.3
- Urgenza: medio
- Cosa NON può chiudere agente: la valutazione di adeguatezza dei requisiti tecnici comunicati per un servizio specifico richiede QM/RT.

## E. Errori fatali correlati
Vedi `[[04_ERRORI_FATALI]]`. Sul § 6.6 i pattern fatali ricorrenti:

- **EF-6.6-A** — Albo fornitori formale ma senza criteri tecnici di qualifica (NC su § 6.6.2 b).
- **EF-6.6-B** — Servizi di taratura considerati "qualifica automatica" perché il fornitore "è accreditato" — senza verifica caso per caso che il parametro tarato rientri nello scopo del lab di taratura (NC trasversale § 6.6 + § 6.5).
- **EF-6.6-C** — Servizi IT/cloud che gestiscono dati di laboratorio non inseriti tra "prodotti e servizi che influenzano le attività" (NC § 6.6.1; collegamento § 7.11 controllo dati).
- **EF-6.6-D** — Ri-valutazione fornitori inesistente o "tacita" (NC su § 6.6.2 b ultima parte).
- **EF-6.6-E** — Conformità pre-uso non verificata (lotto, scadenza, integrità, certificato di analisi): § 6.6.2 c).
- **EF-6.6-F** — Subappalto senza comunicazione esplicita dei requisiti e senza diritto di accesso (NC § 6.6.3).
- **EF-6.6-G** — Azione derivante da NC fornitore mancante o non tracciata (§ 6.6.2 d).
- **EF-6.6-H** — Produttore di CRM "competente" senza dossier di valutazione (collegamento § 6.5 B.3 ultima parte).

## F. Quando l'agente DEVE chiedere prima di rispondere
1. **Servizi di taratura/CRM in casi residuali** — qualsiasi qualifica di fornitore di taratura non accreditato o di CRM fuori casi 4/5/6 si lega al § 6.5 e richiede QM (vedi `[[sez_6_5_riferibilita]]`).
2. **Subappalto di prove accreditate** — l'agente non valuta l'ammissibilità del subappalto senza QM e senza riferimento esplicito a RG-02 vigente (subappalto può essere vietato o limitato per certe categorie).
3. **Software e servizi cloud** — la qualifica di fornitori IT che gestiscono dati di laboratorio è materia complessa (§ 6.6 + § 7.11). L'agente non propone criteri senza QM/IT.
4. **CRM non accreditati / produttori esotici** — il dossier di competenza fornitore è giudizio QM/RT; l'agente raccoglie dati ma non valida.
5. **Ri-valutazione fornitori** — criteri quantitativi richiedono input dal QM (soglie di NC, indicatori di prestazione).
6. **Gas tecnici, acqua di laboratorio** — sono spesso sottovalutati. L'agente segnala come categoria critica e chiede QM se non sono nell'elenco.
