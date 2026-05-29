---
codice: A3_riferibilita_taratura
tipo: appendice_tecnica
livello: misto
titolo: "Appendice A3 — Riferibilità metrologica e scelta del laboratorio di taratura"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "6.4" }
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "6.5" }
  - { doc: "RT-08", rev: "05", paragrafo: "6.4" }
  - { doc: "RT-08", rev: "05", paragrafo: "6.5" }
  - { doc: "UNI-CEI-70099", rev: "VIM 2008/2012" }
fonti_secondarie:
  - { doc: "LS-09", note: "elenco documenti taratura applicabili ai lab di prova" }
  - { doc: "ISO-10012", rev: "2003-2004" }
  - { doc: "DT-08-DL", rev: "00", note: "EMC" }
  - { doc: "DT-0002", rev: "01", note: "guida incertezza" }
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[00_GLOSSARIO]]"
  - "[[03_PROTOCOLLO_HANDOFF_QM]]"
  - "[[02_TRE_LIVELLI_DI_CONTROLLO]]"
  - "[[04_ERRORI_FATALI]]"
  - "[[sez_6_4_dotazioni]]"
  - "[[sez_6_5_riferibilita]]"
  - "[[A2_incertezza_misura]]"
  - "[[A4_materiali_riferimento_CRM]]"
tags: [riferibilita, taratura, ILAC_MRA, EA_MLA, accredia_LAT, taratura_interna, catena_metrologica]
escalation_qm: true
---

> **Marcatore di livello: 🟢🟡🟣** — Appendice trasversale. Contenuti normativi (livello 1+2), prassi ispettiva (livello 3), best practice tecnica (livello 4).

# Appendice A3 — Riferibilità metrologica e scelta del laboratorio di taratura

> **Riassunto operativo.** La riferibilità è una catena: ogni anello deve essere documentato, ogni anello deve dichiarare l'incertezza, ogni anello deve essere riconducibile al SI (o, dove il SI non è applicabile, a un riferimento concordato). Il § 6.5 ISO 17025 + RT-08 sono fra i requisiti **più presidiati** dagli ispettori Accredia: una rottura della catena è quasi sempre NC. Questa appendice serve al QM/RT per impostare la matrice delle tarature, scegliere correttamente i laboratori di taratura esterni e gestire la taratura interna nei limiti consentiti. L'agente la carica quando il QM parla di "certificato di taratura", "scelta lab di taratura", "taratura interna", "riferibilità SI", "fornitore esterno", o quando un certificato in cartella non risulta accreditato.

---

## 1. Principio normativo (🟢 LIVELLO 1+2)

### 1.1 ISO 17025:2018 § 6.4 — Apparecchiature

> Testo a pagamento, non riprodotto. Riassunto operativo:

- Il laboratorio deve avere accesso a tutte le apparecchiature richieste dalle proprie attività e che influenzano i risultati.
- **§ 6.4.4** — verifica di idoneità all'uso prima della messa o rimessa in servizio.
- **§ 6.4.5** — le apparecchiature devono permettere il raggiungimento dell'accuratezza richiesta e/o dell'incertezza di misura per fornire risultati validi.
- **§ 6.4.6** — le apparecchiature devono essere tarate quando: (a) l'accuratezza/incertezza influisce sulla validità del risultato; (b) la taratura è richiesta per stabilire la riferibilità.
- **§ 6.4.7** — il laboratorio deve stabilire un programma di taratura, riesaminato e modificato come necessario.
- **§ 6.4.8 / 6.4.9 / 6.4.10** — etichettatura/stato, gestione delle apparecchiature soggette a uso improprio, controlli intermedi quando necessario.
- **§ 6.4.13** — registrazioni delle apparecchiature: identificazione, posizione, stato di taratura, evidenza di controllo, documenti di taratura.

### 1.2 ISO 17025:2018 § 6.5 — Riferibilità metrologica

> Testo a pagamento, non riprodotto. Riassunto operativo:

- **§ 6.5.1** — il laboratorio deve stabilire e mantenere la riferibilità metrologica dei propri risultati attraverso una catena documentata e ininterrotta di tarature, ciascuna delle quali contribuisce all'incertezza di misura, che collega i risultati a un riferimento appropriato.
- **§ 6.5.2** — la riferibilità al SI è garantita tramite:
  - (a) taratura fornita da un **laboratorio competente** (es. accreditato secondo ISO/IEC 17025), oppure
  - (b) valori certificati di materiali di riferimento certificati (CRM) forniti da produttore competente, oppure
  - (c) realizzazione diretta delle unità SI garantita da confronto con NMI.
- **§ 6.5.3** — quando la riferibilità al SI non è tecnicamente possibile, il laboratorio dimostra la riferibilità a un riferimento appropriato (es. valori certificati di CRM, metodo concordato, scala convenzionale, confronti interlaboratorio).

### 1.3 RT-08 rev.05 § 6.4 e § 6.5 (free, citabile)

> **Fonte**: `[RT-08 rev.05 § 6.4 e § 6.5 — EC 10-02-2022]`

RT-08 è particolarmente prescrittivo sui § 6.4-6.5. Punti chiave (sintesi operativa, citazione integrale nei punti chiave):

- **§ 6.4.6 / 6.5.2** — Le tarature esterne devono essere effettuate da laboratori:
  - accreditati **ACCREDIA-LAT**, oppure
  - accreditati da Enti firmatari degli **Accordi di Mutuo Riconoscimento ILAC-MRA** o **EA-MLA** per l'attività di taratura specifica,
  - oppure realizzate dagli **Istituti Metrologici Nazionali (NMI)** firmatari del CIPM-MRA, nell'ambito delle CMC (Calibration and Measurement Capabilities) pubblicate nel BIPM KCDB.
  
  > Le tarature effettuate da laboratori che non rientrano in queste categorie **non sono accettabili** per dimostrare la riferibilità ai fini dell'accreditamento.

- **§ 6.4.7** — Programma di taratura: il laboratorio definisce intervalli di taratura su base tecnica (uso, deriva, stabilità, criticità rispetto all'incertezza target), riesamina e modifica gli intervalli sulla base dell'esperienza.

- **§ 6.5.3** — RT-08 ammette la **taratura interna** nei casi in cui:
  - il laboratorio dispone di personale, riferimenti e procedure adeguate,
  - l'attività è formalizzata e tracciata,
  - i riferimenti utilizzati sono a loro volta riferibili (tipicamente CRM certificati o strumenti tarati ACCREDIA-LAT/ILAC-MRA),
  - l'incertezza della taratura interna è valutata.
  
  La taratura interna **non sostituisce** la taratura primaria della catena: serve a estendere la riferibilità all'interno del laboratorio o a strumenti secondari.

- **Documenti di riferimento per la taratura applicabili ai laboratori di prova**: cfr. `LS-09` (elenco curato Accredia/EURAMET) e `LS-04` § 2.

### 1.4 VIM (UNI-CEI-70099) — termini di riferimento

> **Fonte**: `[UNI-CEI-70099 / VIM]`

I termini chiave che l'agente usa con rigore:
- **Riferibilità metrologica**: proprietà di un risultato per cui esso può essere ricondotto a un riferimento attraverso una catena documentata e ininterrotta di tarature.
- **Catena di riferibilità**: sequenza di standard di misura e tarature usata per relazionare un risultato a un riferimento.
- **Taratura**: operazione che, sotto condizioni specificate, stabilisce una relazione tra valori indicati e valori di riferimento, considerando le incertezze.
- **Verifica**: fornitura di evidenza oggettiva che un'entità soddisfi requisiti specificati. Verifica ≠ taratura.

---

## 2. Come un ispettore lo verifica (🟡 LIVELLO 3 — prassi ispettiva)

> Pattern: meccanismo → estensione → efficacia.

### 2.1 Domanda 1 — Meccanismo

**"Mostratemi il vostro programma di taratura: per ogni strumento critico, voglio vedere intervallo di taratura, criterio tecnico che lo giustifica, fornitore prescelto, evidenza che il fornitore è accreditato ACCREDIA-LAT o ILAC-MRA o EA-MLA per la grandezza specifica."**

Analisi ispettiva:
- Cerco una **matrice o registro strumenti** che colleghi ogni apparecchiatura a: campo accreditato in cui è usata, intervallo di taratura, ultimo certificato, prossima scadenza, fornitore, accreditamento del fornitore.
- Verifico il **logo Accredia o ILAC-MRA/EA-MLA** sui certificati e la presenza della grandezza tarata nello scope del fornitore. **Errore frequente**: il fornitore è accreditato LAT ma per grandezze diverse da quella necessaria; il certificato è "non accreditato" pur essendo il fornitore accreditato per altro.
- Per le tarature da **NMI**: cerco riferimento alle CMC pubblicate nel BIPM KCDB.

### 2.2 Domanda 2 — Estensione

**"Avete un elenco completo degli strumenti che hanno impatto sui risultati accreditati? Quanti hanno taratura accreditata ACCREDIA-LAT/ILAC-MRA/EA-MLA? Quanti hanno solo taratura non accreditata (e perché)? Quanti sono tarati internamente (e con quale dossier)? Per gli strumenti soggetti a controllo intermedio (§ 6.4.10), il piano è documentato?"**

Analisi ispettiva:
- Campionamento sullo scopo: prendo 3-5 strumenti critici. Per ognuno: certificato + verifica accreditamento + verifica intervallo.
- **Errore fatale ricorrente**: strumento con certificato scaduto da mesi ancora in uso operativo. NC critica.
- **Errore comune**: lo strumento è dichiarato "verificato internamente" ma in cartella non c'è il dossier della taratura interna (riferimenti usati, procedura, incertezza, autorizzazione personale). NC su § 6.5.3 + RT-08.
- Verifico se i **riferimenti usati per la taratura interna** sono a loro volta riferibili (es. CRM con certificato di accreditamento ISO 17034, o strumento tarato ACCREDIA-LAT). Catena interrotta = NC.

### 2.3 Domanda 3 — Efficacia

**"Negli ultimi 24 mesi, è successo che un certificato di taratura abbia evidenziato una deriva fuori criterio o un'incertezza superiore a quella target? Cosa è stato fatto? Mi mostrate il caso."**

Analisi ispettiva:
- La taratura è viva quando viene **letta** al ritorno dal fornitore, non solo archiviata. Cerco evidenza di:
  - confronto deriva vs criterio interno,
  - verifica che l'incertezza dichiarata dal fornitore sia compatibile con l'incertezza target del laboratorio,
  - eventuale azione (riduzione intervallo, riparazione, sostituzione, retroazione sui risultati passati se la deriva è significativa).
- Se la risposta è *"i certificati vanno in archivio, non li riesaminiamo"*, il § 6.4 è formalmente coperto ma sostanzialmente cieco.

### 2.4 Pattern di errori comuni

1. **Certificato non accreditato** per uno strumento usato in attività accreditata. NC su § 6.5.2 + RT-08 § 6.5 (può essere NC grave se il certificato è "di fabbrica" o di fornitore non riconosciuto).
2. **Fornitore accreditato per grandezza diversa**: il certificato porta il logo Accredia LAT, ma quella grandezza specifica non è nello scope del fornitore. Stessa gravità del precedente.
3. **Catena di riferibilità non documentata** per taratura interna. NC su § 6.5.3.
4. **Intervalli di taratura "stabiliti per prassi"** senza base tecnica documentata né riesame. NC su § 6.4.7.
5. **Controlli intermedi assenti** per strumenti che lo richiederebbero (es. bilance critiche, termometri di riferimento). NC su § 6.4.10.
6. **Certificato letto solo come "OK/KO"** senza verifica della compatibilità incertezza fornitore vs incertezza target del laboratorio. NC su § 6.5.2 oppure § 7.6.
7. **Strumento usato dopo la data di scadenza taratura**: errore fatale (uscite del laboratorio in quel periodo potenzialmente compromesse, vedi `[[04_ERRORI_FATALI]]`).

---

## 3. Best practice tecnica (🟣 LIVELLO 4 — best practice)

### 3.1 Riferimento operativo: LS-09 e LS-04 § 2

`LS-09` (Accredia + EURAMET) elenca i documenti applicabili per le attività di taratura. Il laboratorio di prova lo consulta quando deve:
- definire i requisiti tecnici della taratura da richiedere al fornitore,
- valutare se la taratura interna è giustificata e con quali standard.

`LS-04 § 2` rinvia ai cataloghi documentali per i requisiti settoriali. Si raccomanda di mantenere in cartella un'estrazione dei documenti `LS-09` rilevanti per il proprio campo.

### 3.2 Come riconoscere un certificato di taratura accettabile

Indicatori di un certificato che regge in audit Accredia:

| Elemento | Cosa controllare |
|---|---|
| **Logo accreditamento** | Logo ACCREDIA-LAT o logo dell'Ente firmatario ILAC-MRA/EA-MLA, oppure intestazione NMI |
| **Numero accreditamento + scope** | Verificare che la grandezza tarata sia coperta dallo scope (database Accredia / sito dell'Ente / KCDB per NMI) |
| **Riferibilità dichiarata** | Riferimento alla catena di tarature e agli standard usati |
| **Incertezza dichiarata** | Con fattore di copertura k e livello di confidenza (in coerenza con `JCGM-100`) |
| **Condizioni di taratura** | Condizioni ambientali, metodo, configurazione strumento |
| **Data della taratura + identificativo strumento** | Numero di serie tracciabile, data univoca |
| **Firma del personale autorizzato** del fornitore | Presente |

### 3.3 Tarature interne — quando è ammessa e come si documenta

Riferimenti: `RT-08 § 6.5.3`, `ISO-10012:2003-2004` (gestione misurazione), buone prassi metrologiche.

Condizioni minime per ammettere una taratura interna:
1. **Riferimento riferibile**: lo standard usato per la taratura interna (CRM, masse di riferimento, blocchetti, termometro di riferimento, ecc.) è a sua volta tarato ACCREDIA-LAT/ILAC-MRA o è un CRM certificato.
2. **Procedura scritta**: passi, condizioni, criteri, modulistica.
3. **Personale autorizzato** in matrice § 6.2.6 per la taratura interna.
4. **Dossier di incertezza della taratura interna** (vedi `[[A2_incertezza_misura]]`).
5. **Riesame periodico** dell'adeguatezza della taratura interna (ad es. confronto periodico con taratura esterna accreditata su sottoinsieme di strumenti).

> Una taratura interna non documentata in tutti questi punti non è "interna": è solo una verifica non strutturata. Non sostituisce la taratura formale richiesta.

### 3.4 Intervalli di taratura

La determinazione e il riesame degli intervalli (`§ 6.4.7`) si fondano su:
- specifica tecnica del costruttore,
- esperienza storica del laboratorio sulla deriva (vedi `DT-0002/6` su ripetibilità nel tempo, in chiave generale),
- criticità rispetto all'incertezza target,
- frequenza d'uso e condizioni ambientali,
- esiti dei controlli intermedi.

`ISO-10012` fornisce un quadro gestionale per il "sistema di gestione della misurazione" che integra questi elementi. È citabile come buona pratica.

### 3.5 EMC (compatibilità elettromagnetica)

Per i laboratori EMC, la guida settoriale Accredia è `DT-08-DL` (free). Definisce le specificità della catena di riferibilità per gli strumenti EMC (analizzatori di spettro, antenne, network analyzer, generatori di segnale).

### 3.6 Quando l'NMI è la fonte (non-accreditato ma riconosciuto)

Le tarature da NMI firmatari del **CIPM-MRA** sono accettate purché ricadano nelle CMC pubblicate nel BIPM KCDB (`https://www.bipm.org/kcdb`). Per laboratori italiani, il riferimento naturale è **INRIM** (e **INMRI-ENEA** per misure ionizzanti).

### 3.7 CRM come fonte di riferibilità

Vedi `[[A4_materiali_riferimento_CRM]]`. La riferibilità via CRM (§ 6.5.2 b) richiede produttore qualificato (tipicamente accreditato ISO 17034) e certificato che dichiari incertezza e catena.

---

## 4. Template operativo per il QM

### 4.1 Matrice "Registro strumenti e riferibilità" (struttura minima)

| Cod. strumento | Descrizione | Sede | Grandezza/range | Usato in (metodi accreditati) | Tipo taratura (ext LAT / ext ILAC-MRA / ext NMI / interna / verifica) | Fornitore + n. accreditamento + scope | Ultimo cert. (data, n.) | Incertezza dichiarata | Incertezza target lab | OK/NOK | Intervallo (mesi) | Prossima scadenza | Controlli intermedi (sì/no, freq.) | Note / azioni |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|

### 4.2 Checklist "Verifica certificato di taratura in ingresso"

- [ ] Logo accreditamento presente e leggibile (ACCREDIA-LAT / ILAC-MRA / EA-MLA / NMI)
- [ ] Grandezza tarata coperta dallo scope del fornitore (verifica su database)
- [ ] Riferibilità SI dichiarata (catena agli standard nazionali/internazionali)
- [ ] Incertezza dichiarata con k e livello di confidenza
- [ ] Incertezza dichiarata ≤ incertezza target del laboratorio per quell'uso
- [ ] Condizioni di taratura compatibili con condizioni d'uso del laboratorio
- [ ] Identificativo strumento corrisponde al registro
- [ ] Eventuale deriva rispetto a taratura precedente: verificata e annotata
- [ ] Firma personale autorizzato del fornitore presente
- [ ] Certificato archiviato e protocollato; data di scadenza riportata nel registro
- [ ] Esito: utilizzabile / non utilizzabile / utilizzabile con correzione

### 4.3 Decision tree "Posso usare questo certificato per attività accreditata?"

```
Fornitore = ACCREDIA-LAT?  →  Scope copre la grandezza?  →  OK (verifica incertezza)
                          ↓ no
Fornitore = lab estero ILAC-MRA/EA-MLA per la grandezza?  →  OK (verifica incertezza)
                          ↓ no
Fornitore = NMI firmatario CIPM-MRA, dentro CMC pubblicate?  →  OK (verifica incertezza)
                          ↓ no
Taratura interna conforme a § 6.5.3 + RT-08 (riferimento riferibile, procedura, incertezza, autorizzazione)?  →  OK per uso interno; verificare se basta per accreditato
                          ↓ no
NON utilizzabile per accreditato. Riconvocare il fornitore o ri-tarare altrove.
```

### 4.4 Format minimo "Dossier taratura interna"

```
DOSSIER TARATURA INTERNA — STRUMENTO [codice]
Rev. [n] del [data]   Redatto: [tecnico]   Approvato: [RT autorizzato § 6.2.6]

1. SCOPO E CAMPO DI APPLICAZIONE
2. STANDARD DI RIFERIMENTO UTILIZZATO
   - Tipo (CRM / strumento tarato): ...
   - Riferibilità dello standard: [certificato n., fornitore, accreditamento]
   - Incertezza dello standard: ...
3. PROCEDURA DI TARATURA INTERNA
4. CRITERI DI ACCETTAZIONE
5. VALUTAZIONE INCERTEZZA DELLA TARATURA INTERNA
   (vedi A2)
6. PERSONALE AUTORIZZATO (§ 6.2.6)
7. FREQUENZA E RIESAME
8. CONFRONTO PERIODICO CON TARATURA ESTERNA ACCREDITATA (sì/no, quando)
```

---

## 5. Quando l'agente attiva questa appendice

1. Il QM chiede *"questo certificato basta?"*, *"posso fare la taratura in casa?"*, *"come scelgo il lab di taratura?"*.
2. L'agente, ispezionando la cartella strumenti, trova certificati senza logo accreditamento o per fornitori non verificabili.
3. È in lavorazione il programma annuale di taratura (cfr. § 6.4.7) e serve la revisione degli intervalli.
4. Un certificato in ritorno rileva deriva fuori criterio: serve istruttoria su impatto risultati e azione correttiva.
5. Il laboratorio sta valutando di internalizzare alcune tarature (riduzione costi): serve il quadro § 6.5.3 + RT-08.
6. È in lavorazione una NC su § 6.4 / § 6.5.
7. Per laboratori EMC: serve raccordare il quadro con `DT-08-DL`.

---

## 6. Cosa l'agente NON può chiudere su questo tema (handoff QM obbligatorio)

Vedi `[[03_PROTOCOLLO_HANDOFF_QM]] § 4`. In materia di riferibilità, l'agente **non** può:

1. **Dichiarare accettabile un certificato** ai fini dell'accreditamento. Propone l'esito dell'istruttoria; la decisione è del RT/QM.
2. **Autorizzare una taratura interna** come fonte di riferibilità per attività accreditate. Spetta al RT con autorizzazione formale.
3. **Modificare gli intervalli di taratura** del programma. Propone la revisione, decide il RT.
4. **Decidere il fornitore di taratura**. Può istruire la qualificazione, ma la scelta del fornitore è gestionale.
5. **Dichiarare conforme uno strumento** dopo controllo intermedio. La dichiarazione è del personale autorizzato.
6. **Valutare l'impatto sui risultati emessi** quando un certificato rivela deriva significativa. Resta valutazione tecnica del RT (vedi `[[03_RISPOSTA_NC]]`).
7. **Dichiarare chiusa una NC** su § 6.4 / § 6.5. La chiusura spetta al QM.
