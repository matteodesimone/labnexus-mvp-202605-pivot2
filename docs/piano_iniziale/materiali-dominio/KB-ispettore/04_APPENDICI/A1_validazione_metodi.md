---
codice: A1_validazione_metodi
tipo: appendice_tecnica
livello: misto
titolo: "Appendice A1 — Verifica e validazione dei metodi di prova"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.2.1" }
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.2.2" }
  - { doc: "RT-08", rev: "05", paragrafo: "7.2" }
  - { doc: "RT-26", rev: "07", note: "scopo flessibile" }
fonti_secondarie:
  - { doc: "EURACHEM-FFP", rev: "3rd ed. 2025" }
  - { doc: "EURACHEM-CITAC-Quality", rev: "3rd ed. 2016" }
  - { doc: "EURACHEM-MicroAcc", rev: "2023" }
  - { doc: "ISO-21748", rev: "2017" }
  - { doc: "ILAC-G17", rev: "01/2021", note: "incertezza nelle prove — i dati di validazione ne sono input" }
  - { doc: "ISO-7218", rev: "2024" }
  - { doc: "SANTE-11312", rev: "2021", note: "settore agroalimentare/pesticidi" }
  - { doc: "ISO-IEC-TS-23532-1", rev: "2021/2025", note: "settore cyber" }
  - { doc: "JRC-GMO-Flex", rev: "2nd ed. 2014", note: "settore GMO/scopo flessibile" }
  - { doc: "DT-0002/6", rev: "00", note: "ripetibilità" }
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[00_GLOSSARIO]]"
  - "[[02_PROTOCOLLO_ASK_BEFORE_ANSWER]]"
  - "[[03_PROTOCOLLO_HANDOFF_QM]]"
  - "[[02_TRE_LIVELLI_DI_CONTROLLO]]"
  - "[[04_ERRORI_FATALI]]"
  - "[[sez_7_2_selezione_verifica_validazione_metodi]]"
  - "[[A2_incertezza_misura]]"
  - "[[A5_prove_valutative_PT_ILC]]"
tags: [validazione, verifica, metodi, fitness_for_purpose, scopo_flessibile, parametri_prestazionali]
escalation_qm: true
---

> **Marcatore di livello: 🟢🟡🟣** — Appendice trasversale. Contiene il quadro normativo (livello 1+2), la prassi ispettiva (livello 3) e le raccomandazioni delle guide tecniche internazionali (livello 4). I blocchi sono distinti dai marcatori.

# Appendice A1 — Verifica e validazione dei metodi di prova

> **Riassunto operativo.** Verifica e validazione sono due cose diverse. Confonderle è la causa più ricorrente di NC su § 7.2 nelle visite Accredia. Questa appendice serve al QM/RT per impostare correttamente la scelta tra le due, definire il piano dei parametri prestazionali, predisporre il dossier che l'ispettore si aspetta di trovare. L'agente la carica quando il QM parla di "nuovo metodo", "metodo modificato", "estensione di scopo", "PT scadente che potrebbe rivelare un problema di validazione", o quando entra in modalità di preparazione visita Accredia su attività tecnica.

> **🟡 Ponte validazione → incertezza.** I parametri prestazionali determinati in validazione (precisione RSDr/RSDi, bias/recupero, intervallo) sono gli **input dell'incertezza di misura** (§ 7.6). RT-08 § 7.6.3 consente di usare gli scarti tipo di ripetibilità/riproducibilità per stimare l'UM previa verifica delle proprie prestazioni; il riferimento applicativo EA/ILAC è `[ILAC-G17:01/2021]`, l'approccio top-down è in `[ISO-21748:2017]`. Quindi: il piano di validazione va impostato fin dall'inizio in modo da **produrre anche i dati per l'UM**. Vedi [[A2_incertezza_misura]], [[sez_7_6_incertezza]] § B-bis e [[sez_7_2_selezione_verifica_validazione_metodi]] § B-bis. Il "filo" dato di validazione → componente di UM → incertezza dichiarata sul RdP deve essere tracciabile.

---

## 1. Principio normativo (🟢 LIVELLO 1+2)

### 1.1 Riferimento ISO 17025:2018 § 7.2

> Testo a pagamento: non riprodotto integralmente. Riferirsi alla copia ufficiale UNI in possesso del laboratorio. Riassunto operativo:

- **§ 7.2.1 — Selezione e verifica dei metodi.** Il laboratorio usa metodi appropriati per tutte le attività di laboratorio e, dove applicabile, per la valutazione dell'incertezza di misura e per le tecniche statistiche di analisi dei dati. I metodi devono essere la versione vigente, validi e disponibili al personale. Prima di introdurre un metodo, il laboratorio deve **verificare** di essere in grado di eseguirlo correttamente (verifica delle prestazioni dichiarate dal metodo normalizzato nelle proprie condizioni operative).
- **§ 7.2.2 — Validazione dei metodi.** Il laboratorio deve **validare** i metodi non normalizzati, i metodi sviluppati internamente, i metodi normalizzati usati al di fuori del loro scopo previsto e le modifiche di metodi validati. La validazione deve essere estesa quanto necessario a soddisfare le esigenze dell'applicazione o del campo di applicazione previsto (concetto di **fitness for purpose**). Il laboratorio deve registrare i risultati della validazione, la procedura usata, una dichiarazione di idoneità all'uso e l'autorizzazione del personale che esegue la validazione.

### 1.2 Riferimento RT-08 rev.05 § 7.2

> **Fonte**: `[RT-08 rev.05 § 7.2 — EC 10-02-2022]`. Citazioni testuali ammesse (free).

- **§ 7.2.1** rinvio alla norma con prescrizioni sui metodi normalizzati e sulla verifica delle prestazioni dichiarate nelle condizioni operative del laboratorio.
- **§ 7.2.2** ricorda che la validazione deve produrre evidenza documentata che il metodo, nelle condizioni reali del laboratorio, soddisfa i requisiti per l'uso previsto. La dichiarazione di idoneità all'uso (fitness for purpose) deve essere firmata dal personale autorizzato (cfr. § 6.2.6).
- RT-08 ricorda che per lo **scopo flessibile** si applica integralmente `RT-26`: validazione e gestione varianti seguono regole rafforzate.

### 1.3 Verifica vs validazione: la distinzione che fa la differenza

| Aspetto | Verifica (§ 7.2.1) | Validazione (§ 7.2.2) |
|---|---|---|
| Quando si applica | Metodo normalizzato (ISO, UNI, EN, EPA, APAT, IRSA, ASTM, metodo ufficiale di una autorità) applicato nel proprio scopo previsto e senza modifiche significative | Metodo non normalizzato, sviluppato internamente, normalizzato applicato fuori scopo previsto, normalizzato modificato |
| Cosa si dimostra | Capacità del laboratorio di replicare le prestazioni dichiarate dal metodo nelle proprie condizioni operative | Idoneità all'uso del metodo per lo scopo previsto, su tutti i parametri prestazionali rilevanti |
| Estensione tipica | Sottoinsieme dei parametri (tipico: ripetibilità, recupero su CRM o spike, LOQ verificato, incertezza riconducibile a quella della norma) | Set completo dei parametri prestazionali pertinenti al tipo di metodo (vedi § 3.2) |
| Dossier minimo | Piano di verifica + dati grezzi + report con confronto vs prestazioni dichiarate + dichiarazione di idoneità firmata | Piano di validazione + dati grezzi + report con valori dei parametri + criteri di accettazione + dichiarazione di fitness for purpose firmata |
| Frequenza tipica di ripetizione | Su modifica significativa dello strumento, del personale, della matrice; al ri-onboarding di nuovo operatore | Su modifica significativa del metodo, dei reagenti critici, dello strumento, della matrice; alla scadenza definita dal lab nel piano di riesame |

> **Nota di prassi**: la frase *"il metodo è normalizzato, quindi non serve validarlo"* è vera **solo se** il laboratorio: (a) lo applica esattamente nello scopo dichiarato dalla norma, (b) non ha introdotto varianti, (c) ha comunque eseguito la **verifica** delle prestazioni nelle proprie condizioni. Saltare la verifica è errore frequente quanto saltare la validazione.

---

## 2. Come un ispettore lo verifica (🟡 LIVELLO 3 — prassi ispettiva)

> Pattern: meccanismo → estensione → efficacia. Vedi `[[02_TRE_LIVELLI_DI_CONTROLLO]]`.

### 2.1 Domanda 1 — Meccanismo

**"Mostratemi, per una prova accreditata a scelta del campo, il dossier di validazione o di verifica. Voglio vedere: piano scritto a priori, criteri di accettazione numerici, dati grezzi, report finale con confronto vs criteri, dichiarazione di idoneità all'uso firmata dal personale autorizzato."**

Analisi ispettiva:
- Cerco la **sequenza completa**, non solo il report finale. Un report di validazione senza il piano a priori è sospetto: senza criteri di accettazione preregistrati, il "passa" è retroattivo.
- Verifico la **firma** della dichiarazione di idoneità: deve essere del personale autorizzato a "sviluppo/modifica/validazione metodi" (RT-08 § 6.2.6). Se la firma è del titolare senza autorizzazione specifica nella matrice, è NC su § 6.2.6 oltre che su § 7.2.
- Controllo la **distinzione verifica/validazione**: se la prova è su metodo non normalizzato e in cartella c'è solo "verifica", manca proprio l'oggetto del § 7.2.2.

### 2.2 Domanda 2 — Estensione

**"Quanti metodi accreditati avete? Per quanti di questi il dossier è completo come quello appena mostrato? Esiste un registro/mappa che, per ogni metodo dello scopo, indica stato del dossier (verifica/validazione), data ultimo riesame, prossima scadenza?"**

Analisi ispettiva:
- L'estensione tipica si verifica con un **campionamento sullo scopo**: prendo 3-5 metodi dal CAdA (Campo di Accreditamento), chiedo il dossier; se anche uno solo manca o è incompleto, scatta NC.
- Cerco la **mappa di copertura**: un foglio (anche un foglio Excel) che lega ogni voce dello scopo al dossier corrispondente. Senza mappa, il laboratorio non sa che cosa ha validato.
- Errore comune: dossier di validazione fatto bene per il "metodo madre", e poi nuove matrici o nuovi range aggiunti per anni senza estensione formale. È il caso classico di **scope creep** silenzioso.

### 2.3 Domanda 3 — Efficacia

**"Nell'ultimo PT/ILC scadente o nell'ultima NC interna su un risultato, avete riesaminato il dossier di validazione del metodo per verificare che la causa non fosse un parametro prestazionale fuori specifica? Mi mostrate un caso reale degli ultimi 24 mesi."**

Analisi ispettiva:
- La validazione è viva quando viene **interrogata dopo** un evento (PT scadente, NC su risultato, reclamo, deriva del CQ interno). Se il dossier viene "consegnato" alla validazione iniziale e poi mai più toccato, è morto.
- Cerco la **catena**: PT scadente → analisi causa → riesame parametro prestazionale interessato (es. recupero, ripetibilità) → conferma/revisione dichiarazione di idoneità.
- Se la risposta è *"non è mai successo niente"* su un parco metodi significativo, o il sistema di controllo qualità è cieco o il riesame non avviene davvero.

### 2.4 Pattern di errori comuni (osservati in ispezione)

1. **Verifica spacciata per validazione**, su metodo non normalizzato modificato. NC su § 7.2.2.
2. **Validazione "una tantum"** senza piano di riesame periodico né trigger di ri-validazione documentati.
3. **Parametri prestazionali parziali**: si valida solo ripetibilità e recupero, si omettono LOD/LOQ, robustezza, intervallo di lavoro, selettività. NC su § 7.2.2 (idoneità all'uso non dimostrata sul completo).
4. **Criterio di accettazione assente o post-hoc**: il report dichiara "esito positivo" senza un numero target preregistrato. NC su prassi tecnica + § 7.2.2.
5. **Estensione informale dello scopo**: il metodo è validato per matrice X e poi applicato a matrice Y senza ri-validazione. NC su § 7.2.2 + possibile errore fatale su § 7.8 (rapporto fuori scope).
6. **Validazione firmata da personale non autorizzato** in matrice § 6.2.6. NC combinata.
7. **Scopo flessibile dichiarato ma gestito come fisso** (manca procedura ex `RT-26`, manca elenco varianti, manca matrice di responsabilità del riesame). NC su `RT-26`.

---

## 3. Best practice tecnica (🟣 LIVELLO 4 — best practice)

> Le indicazioni che seguono provengono dalle guide tecniche internazionali. **Non sono obblighi normativi** (la norma chiede "validazione adeguata", non un set specifico di parametri). Sono best practice consolidate che l'ispettore si aspetta di vedere applicate in modo proporzionato.

### 3.1 Approccio fitness for purpose (Eurachem FFP)

Si raccomanda (`EURACHEM-FFP 3rd ed. 2025`) di:
- Definire **prima** lo scopo del metodo (intended use): per chi, per quale decisione, con quale soglia di tolleranza sull'errore.
- Da quello derivare la **specifica di prestazione** (performance specification): quale incertezza target, quale LOQ minimo, quale recupero ammissibile.
- La validazione è il **confronto** tra prestazioni misurate e specifica. Senza specifica preregistrata, non c'è validazione: c'è caratterizzazione.

### 3.2 Parametri prestazionali tipici (insieme di riferimento, da modulare per tipo di prova)

| Parametro | Significato sintetico | Quando obbligatorio (FFP) | Riferimento citabile |
|---|---|---|---|
| **Selettività / specificità** | Capacità del metodo di misurare l'analita in presenza di interferenti | Sempre per metodi quantitativi su matrice complessa | `EURACHEM-FFP` |
| **Linearità / funzione di risposta** | Relazione segnale-concentrazione nell'intervallo di lavoro | Metodi strumentali con calibrazione | `EURACHEM-FFP`, `SANTE-11312` (pesticidi) |
| **Intervallo di lavoro / di misura** | Estremi entro cui le prestazioni sono garantite | Sempre | `EURACHEM-FFP` |
| **LOD (limite di rilevazione)** | Minima quantità rilevabile con confidenza | Metodi di tracce / trace analysis | `EURACHEM-FFP`, `SANTE-11312` |
| **LOQ (limite di quantificazione)** | Minima quantità quantificabile con incertezza accettabile | Sempre quando si quantificano basse concentrazioni | `EURACHEM-FFP`, `SANTE-11312` |
| **Ripetibilità (sr)** | Variabilità nelle stesse condizioni (stesso operatore, strumento, breve intervallo) | Sempre per metodi quantitativi | `EURACHEM-FFP`, `DT-0002/6`, `ISO-21748` |
| **Precisione intermedia / riproducibilità interna** | Variabilità con operatori/giorni/strumenti diversi nello stesso laboratorio | Sempre per metodi quantitativi | `EURACHEM-FFP`, `ISO-21748` |
| **Riproducibilità (sR)** | Variabilità interlaboratorio (richiede ILC o riferimento a studi collaborativi) | Per metodi normalizzati con dati ILC disponibili | `ISO-21748` |
| **Veridicità / recupero** | Vicinanza tra valore misurato e valore di riferimento (CRM o spike) | Metodi chimici quantitativi | `EURACHEM-FFP`, `EURACHEM-CITAC-Quality` |
| **Robustezza / ruggedness** | Insensibilità del metodo a piccole variazioni controllate dei parametri operativi | Metodi sviluppati in casa, applicazione critica | `EURACHEM-FFP` |
| **Incertezza di misura** | Stima del campo plausibile dei valori (vedi `[[A2_incertezza_misura]]`) | Sempre quando applicabile (§ 7.6) | `JCGM-100`, `EURACHEM-CITAC-QUAM`, `ISO-21748` |

### 3.3 Settori specialistici (citabili come riferimento esterno)

- **Microbiologia alimentare**. `ISO-7218:2024` e `EURACHEM-MicroAcc 2023` forniscono il quadro per le particolarità (conteggi, incertezza di conteggio, recupero, scelta del metodo enumerativo) e si combinano con `ISO-19036:2020` per l'incertezza.
- **Pesticidi e residui in alimenti/mangimi**. `SANTE-11312:2021` definisce i criteri specifici di validazione (recupero per livello, ripetibilità, LOQ definito a CCα ecc.). È riferimento tecnico forte ma non sostituisce ISO 17025 § 7.2.
- **GMO e scopo flessibile**. `JRC-GMO-Flex 2nd ed. 2014` (JRC) è la guida tecnica europea per validazione dei metodi GMO e applicazione dello scopo flessibile in quel settore.
- **Cybersecurity / IT security testing**. `ISO-IEC-TS-23532-1:2021` (UNI 2025) integra il § 7.2 con i requisiti di competenza per laboratori di prova di sicurezza IT (CC, ISO/IEC 15408).

### 3.4 Ripetibilità nel tempo (DT-0002/6)

`DT-0002/6` (Accredia, free, livello 2 sul punto specifico ma da intendersi guida applicativa) raccomanda di:
- Calcolare la ripetibilità in fase di validazione iniziale e di **monitorarla nel tempo** con carte di controllo o ripetizioni periodiche.
- Definire un criterio di accettabilità della deriva (es. test F sulla varianza, o limite Shewhart su s).
- Documentare il riesame periodico della ripetibilità come parte integrante dell'assicurazione di validità dei risultati (`§ 7.7` ISO 17025, vedi `[[A5_prove_valutative_PT_ILC]]`).

---

## 4. Template operativo per il QM

### 4.1 Checklist "Quando devo fare una validazione" (decisione preliminare)

| Domanda | Sì → validazione (§ 7.2.2) | No → verifica (§ 7.2.1) |
|---|---|---|
| Il metodo è normalizzato (ISO/UNI/EN/EPA/IRSA/APAT/metodo ufficiale)? | No | Sì |
| È applicato nello scopo previsto dalla norma (matrice, range, analita)? | No | Sì |
| Sono state introdotte modifiche significative (reagenti diversi, step omessi/aggiunti, strumentazione diversa da quella prescritta, condizioni operative fuori range)? | Sì | No |
| È un metodo sviluppato in casa o adattato da letteratura non normata? | Sì | — |
| È un metodo normalizzato applicato su matrice/analita/range non coperti dalla norma? | Sì | — |

> Se **una sola** delle condizioni colonna "Sì → validazione" è vera, occorre validazione completa, non verifica.

### 4.2 Template di "Piano di validazione" (bozza adattabile)

```
PIANO DI VALIDAZIONE METODO [codice]
Rev. [n] del [data]
Redatto da: [tecnico]   Approvato da: [RT autorizzato § 6.2.6]

1. SCOPO DEL METODO
   - Matrice/i: ...
   - Analita/i: ...
   - Intervallo di applicazione: ...
   - Uso previsto del risultato (decisione, criterio di conformità, cliente tipo): ...
   - Specifica di prestazione target (FFP):
       - Incertezza target: ...
       - LOQ minimo: ...
       - Recupero ammissibile: ...
       - Ripetibilità target (sr%): ...

2. RIFERIMENTI
   - Metodo base: [riferimento, edizione]
   - Modifiche introdotte: [elenco puntuale]
   - Documenti tecnici applicati (Eurachem FFP, ISO 21748, settoriali, DT-0002/6): ...

3. PARAMETRI DA VALUTARE (e disegno sperimentale per ciascuno)
   - Selettività / specificità: [come si valuta, su quali campioni]
   - Linearità: [n. livelli, n. replicati, criterio]
   - Intervallo di lavoro: [estremi]
   - LOD / LOQ: [approccio: deviazione del bianco, S/N, calibrazione a basse concentrazioni]
   - Ripetibilità: [n. replicati su almeno n. livelli, criterio (sr% target)]
   - Precisione intermedia: [n. giorni × n. operatori × n. strumenti, criterio]
   - Veridicità / recupero: [CRM o spike, livelli, criterio (% recupero)]
   - Robustezza: [parametri variati, range, criterio]
   - Incertezza di misura: [approccio: vedi A2; bottom-up / top-down]

4. CRITERI DI ACCETTAZIONE
   [Tabella parametro → valore target → criterio di superamento]

5. CRONOPROGRAMMA E RESPONSABILI

6. RIESAME E RI-VALIDAZIONE
   - Trigger di ri-validazione: [modifica metodo, cambio strumento, PT scadente, deriva CQ]
   - Frequenza minima di riesame del dossier: [annuale / biennale / al ciclo Accredia]
```

### 4.3 Checklist minima del dossier (per ispezione)

- [ ] Piano di validazione firmato **prima** dell'esecuzione
- [ ] Riferimenti documentali (norma base + modifiche)
- [ ] Dati grezzi (fogli di calcolo, output strumentali) tracciabili
- [ ] Report di validazione con valore di ogni parametro + criterio + esito
- [ ] Stima dell'incertezza di misura (rinvio a dossier separato, vedi `[[A2_incertezza_misura]]`)
- [ ] Dichiarazione di idoneità all'uso (fitness for purpose) firmata da personale autorizzato (`§ 6.2.6`)
- [ ] Riferimento incrociato al metodo nella matrice autorizzazioni (chi è autorizzato a eseguirlo) e nello scopo di accreditamento
- [ ] Piano di riesame periodico + trigger di ri-validazione

### 4.4 Registrazioni minime da mantenere a sistema

- Registro metodi vs dossier (mappa di copertura dello scopo)
- Storico revisioni del dossier
- Evidenze del riesame periodico (data, motivazione, esito)
- Collegamento PT/ILC ↔ dossier di validazione (se un PT scadente ha innescato riesame)

---

## 5. Quando l'agente attiva questa appendice

1. Il QM chiede *"come imposto la validazione di un nuovo metodo?"* o *"questa è verifica o validazione?"*.
2. Il QM segnala un nuovo metodo in onboarding o un'estensione dello scopo di accreditamento.
3. L'agente rileva, nel dialogo o nei documenti del SGQ, un PT/ILC scadente o una NC interna che potrebbe riflettere un parametro prestazionale non più rispettato.
4. La visita Accredia è in preparazione e l'agente sta costruendo il dossier "metodi" per audit interno propedeutico (vedi `[[A6_audit_interni_riesame_direzione]]`).
5. Il laboratorio dichiara o sta dichiarando scopo flessibile (`RT-26`): la validazione assume regole rafforzate.
6. L'agente sta redigendo una bozza di risposta NC su § 7.2 (vedi `[[03_RISPOSTA_NC]]`).

---

## 6. Cosa l'agente NON può chiudere su questo tema (handoff QM obbligatorio)

Vedi `[[03_PROTOCOLLO_HANDOFF_QM]] § 4`. In materia di validazione, l'agente **non** può:

1. **Firmare la dichiarazione di idoneità all'uso** (fitness for purpose). Spetta al personale autorizzato ex § 6.2.6 (tipicamente RT).
2. **Decidere che la validazione è sufficiente** per un nuovo metodo. L'agente predispone la bozza, il RT/QM giudica l'adeguatezza tecnica.
3. **Definire la specifica di prestazione target** (FFP) per il cliente: dipende dall'uso del risultato, conoscenza che resta del RT/cliente/QM.
4. **Estendere lo scopo di accreditamento** sulla base di una validazione interna. La variazione di scopo va comunicata ad Accredia secondo `RG-02` e `RT-23`.
5. **Approvare l'applicazione dello scopo flessibile**: richiede procedura ex `RT-26` e autorizzazione formale del RT.
6. **Dichiarare chiusa una NC su § 7.2**. La chiusura spetta al QM (vedi `[[03_RISPOSTA_NC]]`).
7. **Decidere su risultati già emessi** con metodo che si rivela non validato sull'estensione effettivamente applicata. Il giudizio tecnico sull'impatto è del RT, la decisione sulla comunicazione clienti è della Direzione.

> **Nota di sistema**: l'agente, su § 7.2, è uno strumento di **istruttoria documentale**. La validazione resta una responsabilità tecnica del personale autorizzato del laboratorio.
