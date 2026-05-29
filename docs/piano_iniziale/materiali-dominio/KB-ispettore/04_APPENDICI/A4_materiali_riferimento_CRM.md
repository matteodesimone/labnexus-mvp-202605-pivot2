---
codice: A4_materiali_riferimento_CRM
tipo: appendice_tecnica
livello: misto
titolo: "Appendice A4 — Materiali di riferimento (RM) e Materiali di riferimento certificati (CRM)"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "6.5" }
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.7" }
  - { doc: "RT-08", rev: "05", paragrafo: "6.5" }
  - { doc: "UNI-CEI-70099", rev: "VIM 2008/2012" }
fonti_secondarie:
  - { doc: "ISO-17034", rev: "2017", note: "riferimento esterno, produttori RM" }
  - { doc: "ISO-33401", rev: "2024", note: "contenuti certificati RM" }
  - { doc: "ISO-33403", rev: "2024", note: "requisiti per l'uso" }
  - { doc: "ISO-33405", rev: "2024", note: "caratterizzazione e valutazione omogeneità/stabilità" }
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[00_GLOSSARIO]]"
  - "[[03_PROTOCOLLO_HANDOFF_QM]]"
  - "[[02_TRE_LIVELLI_DI_CONTROLLO]]"
  - "[[sez_6_5_riferibilita]]"
  - "[[sez_7_7_assicurazione_validita]]"
  - "[[A1_validazione_metodi]]"
  - "[[A2_incertezza_misura]]"
  - "[[A3_riferibilita_taratura]]"
  - "[[A5_prove_valutative_PT_ILC]]"
tags: [CRM, RM, materiali_riferimento, ISO17034, riferibilita, omogeneita, stabilita]
escalation_qm: true
---

> **Marcatore di livello: 🟢🟡🟣** — Appendice trasversale. Contenuti normativi (livello 1+2), prassi ispettiva (livello 3), best practice tecnica (livello 4).

# Appendice A4 — Materiali di riferimento (RM e CRM)

> **Riassunto operativo.** Il CRM è uno dei tre canali ammessi dalla ISO 17025 § 6.5.2 per dimostrare la riferibilità (insieme alla taratura accreditata e alla realizzazione diretta dell'unità SI). Ma "CRM" è un termine spesso usato in modo improprio: molti laboratori chiamano "CRM" qualunque materiale fornito con un valore di riferimento. La distinzione **RM** vs **CRM** e la qualifica del produttore (ISO 17034) sono punti verificati con attenzione dall'ispettore. Questa appendice serve al QM/RT per impostare la gestione documentale e operativa dei materiali di riferimento. L'agente la carica quando il QM parla di "CRM", "standard di riferimento", "spike", "controllo di recupero", "ricondizionamento del riferimento", "fornitore CRM".

---

## 1. Principio normativo (🟢 LIVELLO 1+2)

### 1.1 ISO 17025:2018 — punti di contatto

> Testo a pagamento, non riprodotto. Sintesi operativa:

- **§ 6.5.2 b** — la riferibilità metrologica al SI può essere garantita tramite "valori certificati di materiali di riferimento certificati forniti da un produttore competente, con riferibilità metrologica dichiarata al SI".
- **§ 6.5.3** — quando il SI non è applicabile, il riferimento può essere un valore certificato di un RM (anche non SI), un metodo concordato, una scala convenzionale, un risultato di confronto interlaboratorio.
- **§ 7.7.1** — il laboratorio deve avere una procedura per monitorare la validità dei risultati. L'elenco non esaustivo dei mezzi include esplicitamente: uso di materiali di riferimento (lett. b) e uso di materiali di riferimento certificati o controlli di qualità (lett. c).

### 1.2 RT-08 rev.05 § 6.5

> **Fonte**: `[RT-08 rev.05 § 6.5 — EC 10-02-2022]`

RT-08 ricorda che, ai fini della riferibilità via CRM:
- il **produttore** deve essere "competente": è considerato tale chi è accreditato secondo **`ISO-17034`** per il tipo di materiale fornito, o riconosciuto come tale (es. NMI, JRC, NIST, BCR, IRMM, LGC, NMIA, sulla base degli accordi internazionali e dei database CIPM-MRA).
- il **certificato** del CRM deve dichiarare valore di riferimento, incertezza, riferibilità, condizioni di conservazione, scadenza, modalità d'uso.

### 1.3 Vocabolario VIM

> **Fonte**: `[UNI-CEI-70099 / VIM]`

| Termine | Significato (sintesi) |
|---|---|
| **Materiale di riferimento (RM)** | Materiale sufficientemente omogeneo e stabile rispetto a proprietà specificate, idoneo all'uso previsto in una misurazione o nell'esame di proprietà qualitative |
| **Materiale di riferimento certificato (CRM)** | RM caratterizzato da un processo metrologicamente valido per una o più proprietà specificate, accompagnato da un certificato che fornisce il valore della proprietà, la sua incertezza e una dichiarazione di riferibilità metrologica |

> La differenza è sostanziale: un **RM** non ha necessariamente certificato di riferibilità con incertezza dichiarata; un **CRM** sì. **Solo il CRM, prodotto da fornitore competente, soddisfa § 6.5.2 b.**

### 1.4 ISO/IEC 17034 e famiglia ISO 33000 (riferimenti esterni)

> Questi documenti sono **fuori scope** dell'accreditamento del laboratorio di prova: regolano i produttori di RM. Sono citati come riferimento esterno per capire cosa il laboratorio deve cercare nel fornitore.

| Documento | Anno | Cosa norma |
|---|---|---|
| `ISO-17034` | 2017 | Requisiti generali per la competenza dei produttori di materiali di riferimento |
| `ISO-33401` | 2024 | Contenuti minimi di certificati, etichette e documentazione di accompagnamento |
| `ISO-33403` | 2024 | Requisiti e raccomandazioni per l'uso degli RM |
| `ISO-33405` | 2024 | Approcci per la caratterizzazione e la valutazione di omogeneità e stabilità |

> Nota: ISO 33401/33403/33405 sono **citabili come buona pratica**. Non sostituiscono il § 6.5.2-3 ISO 17025 dal punto di vista dell'accreditamento del lab di prova.

---

## 2. Come un ispettore lo verifica (🟡 LIVELLO 3 — prassi ispettiva)

> Pattern: meccanismo → estensione → efficacia.

### 2.1 Domanda 1 — Meccanismo

**"Mostratemi un CRM in uso. Voglio: certificato del materiale, evidenza che il produttore è accreditato ISO 17034 (o equivalente riconosciuto), condizioni di conservazione documentate, registro di carico/scarico, evidenza di apertura e ricondizionamento, scadenza."**

Analisi ispettiva:
- Cerco il **certificato originale** (PDF o cartaceo): controllo che riporti valore + U + riferibilità + scadenza + condizioni d'uso + condizioni di stoccaggio.
- Verifico se il produttore è dichiaratamente accreditato `ISO 17034` o è un riferimento internazionale riconosciuto (NMI, JRC-IRMM, NIST, ecc.). Un materiale "di un fornitore generico" venduto come "standard di riferimento" non è automaticamente CRM.
- Controllo le **condizioni di conservazione**: temperatura, esposizione luce, umidità, scadenza dopo apertura. Un CRM conservato in modo non conforme alle condizioni del certificato non garantisce più il valore certificato.

### 2.2 Domanda 2 — Estensione

**"Per ogni metodo accreditato che usa CRM, esiste un piano dei CRM (quali, per quale parametro prestazionale, frequenza d'uso, fornitore, scadenza)? Quanti CRM sono stati sostituiti da RM non certificati per ragioni di costo o disponibilità? Come è gestita la transizione?"**

Analisi ispettiva:
- Cerco la **lista CRM** e la sua coerenza con i dossier di validazione (`[[A1_validazione_metodi]]`) e di incertezza (`[[A2_incertezza_misura]]`).
- **Errore frequente**: un metodo è stato validato con un CRM specifico, poi quel CRM è scaduto o non più disponibile e si usa un RM sostitutivo senza ri-validazione/verifica del recupero. NC su § 7.2 + § 6.5.
- Verifico se il laboratorio usa CRM secondari **preparati internamente** (es. soluzioni intermedie diluite dal CRM madre): la riferibilità si trasmette, ma occorre dossier (procedura, incertezza della diluizione, etichettatura, scadenza ri-calcolata).

### 2.3 Domanda 3 — Efficacia

**"Quando un controllo su CRM ha dato un valore fuori range degli ultimi 12 mesi, cosa è stato fatto? Mi mostrate il caso e la sequenza che ne è seguita."**

Analisi ispettiva:
- Il CRM è uno dei pilastri dell'assicurazione di validità (§ 7.7). Cerco evidenza che gli esiti dei CRM finiscano nei **registri di controllo qualità** (carte di Shewhart, registri tabulari, log).
- Una serie di valori "sempre OK" può essere un buon controllo o una verifica cieca: indaga il **criterio di accettazione** e se è stato mai violato e con quale conseguenza.

### 2.4 Pattern di errori comuni

1. **Confusione RM/CRM**: il laboratorio chiama "CRM" un materiale che non ha certificato di riferibilità con incertezza. NC su § 6.5.2.
2. **Fornitore non qualificato**: il certificato proviene da un produttore non accreditato ISO 17034 né riconosciuto come NMI. NC su § 6.5.2 + RT-08.
3. **CRM conservato male**: condizioni d'uso non rispettate (temperatura, luce, scadenza post-apertura). NC tecnica su § 7.7 + potenziale impatto su risultati.
4. **CRM scaduto ancora in uso**: errore fatale, vedi `[[04_ERRORI_FATALI]]`. Tutti i risultati che hanno usato il CRM scaduto come riferimento di taratura/QC sono potenzialmente compromessi.
5. **Sostituzione silenziosa**: il dossier di validazione e di incertezza menziona CRM A, in operativa si usa CRM B (diverso fornitore o lotto) senza aggiornamento. NC su § 7.2 + § 6.5.
6. **CRM secondario preparato internamente senza dossier**: catena di riferibilità interrotta. NC su § 6.5.
7. **Esiti dei CRM non confluiscono nel CQ**: il § 7.7 chiede uso effettivo dei CRM per monitorare la validità; se gli esiti non sono registrati né riesaminati, il controllo è formale.

---

## 3. Best practice tecnica (🟣 LIVELLO 4 — best practice)

### 3.1 Cosa cercare in un certificato CRM (ISO 33401)

Si raccomanda (`ISO-33401:2024`) che il certificato di un CRM contenga almeno:
- Identificazione univoca del materiale e del lotto
- Proprietà certificata/e con valore + incertezza + livello di confidenza/k
- Dichiarazione di riferibilità metrologica (al SI o a riferimento appropriato)
- Caratterizzazione: metodo/i usato/i per la determinazione del valore certificato
- Omogeneità e stabilità: contributi quantificati nell'incertezza certificata (cfr. `ISO-33405`)
- Condizioni di conservazione e di uso (range temperatura, esposizione, umidità, scadenza pre/post apertura)
- Quantità minima per l'analisi (sub-sampling) garantita per mantenere la rappresentatività
- Eventuali avvertenze (es. omogeneizzare prima dell'uso, atmosfera inerte, vortex, scongelamento controllato)

### 3.2 Uso operativo (ISO 33403)

`ISO-33403:2024` raccomanda:
- **Apertura tracciata**: data di apertura, operatore, lotto, condizioni ambientali al momento.
- **Sub-sampling**: quantità prelevata in coerenza con il minimum sample size dichiarato; se inferiore, il valore certificato può non essere garantito (omogeneità intra-bottiglia).
- **Conservazione post-apertura**: rispetto delle condizioni; ricondizionamento (es. essiccatore, freezer −20°C, atmosfera modificata) se previsto.
- **Scadenza post-apertura**: spesso più breve di quella su flacone chiuso; va calcolata e registrata.
- **Registro di carico/scarico** per CRM critici (numero residuo di prove possibili → trigger di riordino).
- **Smaltimento**: tracciato a fine vita; non si "recupera" un CRM scaduto.

### 3.3 RM secondari preparati internamente (working standards / QC samples)

Pratica comune e ammissibile, a condizione che:
- la **diluizione o preparazione** segua una procedura documentata,
- il **CRM madre** sia certificato e in validità,
- l'**incertezza del lavoro intermedio** sia stimata (incertezza CRM + incertezza diluizione/pesata),
- l'etichetta riporti: lotto, data preparazione, scadenza, condizioni di conservazione, operatore,
- l'uso sia tracciato (registro).

### 3.4 Fornitori CRM riconosciuti (esempi non esaustivi)

Fornitori storicamente riconosciuti (accreditati ISO 17034 o equivalenti per le proprie attività):
- **NMI**: NIST (USA), INRIM (IT), NMI (NL), NPL (UK), PTB (DE), LNE (FR)
- **JRC**: ex IRMM (Geel, BE) — materiali ERM
- **LGC**: Reference Materials
- **Sigma-Aldrich / Merck**, **Restek**, **Inorganic Ventures**, **AccuStandard**, **HPS** — per molti materiali rilevanti, verificare di volta in volta l'accreditamento ISO 17034 sullo specifico lotto/materiale

> L'agente, per uno specifico fornitore non in elenco, **non dichiara** la conformità ISO 17034: chiede al QM di verificare sul certificato del prodotto e sul sito del fornitore lo specifico scope di accreditamento.

### 3.5 Quando il CRM non esiste

Per matrici o analiti senza CRM disponibile, le opzioni sono:
- **RM interno caratterizzato** mediante studio collaborativo o partecipazione a PT specifico,
- **Spike** in matrice rappresentativa (ammesso ma con limiti: il recupero da spike non sempre rispecchia il recupero da analita nativo),
- **Metodo concordato** con il cliente (§ 6.5.3),
- **Confronto interlaboratorio** dedicato (vedi `[[A5_prove_valutative_PT_ILC]]`).

In tutti i casi serve dossier che giustifichi la scelta e ne valuti l'incertezza.

---

## 4. Template operativo per il QM

### 4.1 Registro CRM/RM (struttura minima)

| Cod. interno | Tipo (CRM/RM/secondario) | Materiale / matrice | Analita/i + valore certificato + U | Fornitore + accreditamento ISO 17034 (n.) | Lotto | Data ricezione | Data apertura | Scadenza chiuso | Scadenza aperto | Condizioni conservazione | Metodo/i in cui è usato | Frequenza d'uso | Saldo | Note (carte CQ, esiti) |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|

### 4.2 Checklist "Verifica idoneità CRM in ingresso"

- [ ] Certificato fornito con riferibilità + valore + U + k + scadenza
- [ ] Produttore accreditato ISO 17034 per quel tipo di materiale (verifica scope) oppure NMI/JRC riconosciuto
- [ ] Condizioni di conservazione note e rispettate (temperatura, esposizione, atmosfera)
- [ ] Scadenza pre-apertura > orizzonte di utilizzo previsto
- [ ] Quantità sufficiente per il piano di utilizzo (n. analisi previste × consumo per analisi)
- [ ] CRM registrato nel registro CRM/RM con cod. interno
- [ ] Aggiornati i dossier di validazione e incertezza dei metodi che lo usano (se è nuovo o sostituisce uno precedente)
- [ ] Etichettatura interna univoca

### 4.3 Decision tree "RM o CRM?"

```
Materiale ha certificato con valore + U + dichiarazione riferibilità?
    Sì → CRM (potenziale § 6.5.2 b)
        Produttore accreditato ISO 17034 per quel materiale o NMI riconosciuto?
            Sì → CRM utilizzabile per riferibilità ai fini accreditamento
            No → CRM tecnicamente, ma NON garantisce § 6.5.2 b. Trattare come RM avanzato.
    No → RM (utilizzabile per QC / sviluppo metodo / verifica, NON per riferibilità formale ex § 6.5.2)
```

### 4.4 Format etichetta RM secondario preparato internamente

```
[Codice interno]
Soluzione di lavoro [analita] in [matrice]
Concentrazione: [valore] ± [U] (k=2)
Preparato da: CRM [cod. madre, lotto] mediante [procedura, riferimento]
Data preparazione: [data]   Operatore: [nome]
Conservazione: [condizioni]
Scadenza: [data]
```

---

## 5. Quando l'agente attiva questa appendice

1. Il QM chiede *"questo è un CRM?"*, *"posso usare questo materiale come riferimento?"*, *"il fornitore è qualificato?"*.
2. Un CRM in cartella sta per scadere o è stato sostituito: serve istruttoria sull'impatto su validazione e incertezza.
3. Si prepara un working standard / RM secondario interno: serve dossier.
4. È in lavorazione il piano CQ (`§ 7.7`) e i CRM sono uno dei pilastri (vedi `[[A5_prove_valutative_PT_ILC]]`).
5. Un risultato di CRM è fuori range del criterio: serve indagine NC.
6. Visita Accredia in preparazione: serve verifica della coerenza registro CRM ↔ dossier validazione/incertezza ↔ uso operativo.

---

## 6. Cosa l'agente NON può chiudere su questo tema (handoff QM obbligatorio)

Vedi `[[03_PROTOCOLLO_HANDOFF_QM]] § 4`. In materia di CRM/RM, l'agente **non** può:

1. **Qualificare formalmente un fornitore** di CRM come "competente" ai fini § 6.5.2 b. Propone l'esito dell'istruttoria documentale; la qualifica è del RT/QM.
2. **Approvare un RM secondario preparato internamente** come riferimento. Serve firma del RT autorizzato.
3. **Decidere la sostituzione di un CRM** con uno nuovo o di fornitore diverso. Implica ricaduta su validazione e incertezza: decisione del RT.
4. **Dichiarare un CRM scaduto ancora utilizzabile** "in deroga". Non è ammissibile: la decisione di gestire l'evento è del RT (di norma: ritiro immediato dall'uso operativo + valutazione impatto sui risultati emessi).
5. **Valutare l'impatto sui risultati emessi** se un CRM fuori range rivela un problema sistemico. Resta valutazione del RT (vedi `[[03_RISPOSTA_NC]]`).
6. **Dichiarare chiusa una NC** che tocca § 6.5 / § 7.7 per uso scorretto di CRM. La chiusura è del QM.
