---
codice: A5_prove_valutative_PT_ILC
tipo: appendice_tecnica
livello: misto
titolo: "Appendice A5 — Prove valutative (PT) e confronti interlaboratorio (ILC)"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.7" }
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.10" }
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "8.7" }
  - { doc: "RT-08", rev: "05", paragrafo: "7.7" }
  - { doc: "RT-39", rev: "00", note: "Prescrizioni Accredia partecipazione PT/ILC" }
fonti_secondarie:
  - { doc: "ISO-17043-2010", note: "riferimento esterno PT providers" }
  - { doc: "ISO-17043-2024", note: "riferimento esterno PT providers (ISO/IEC 17043:2023)" }
  - { doc: "EURACHEM-PT", rev: "3rd ed. 2021" }
  - { doc: "EURACHEM-MicroAcc", rev: "2023", note: "microbiologia" }
  - { doc: "Nordtest-TR-569", rev: "ed. 5.1", note: "Trollbook" }
  - { doc: "ISTISAN-12/29", rev: "4ª ed. 2011", note: "trad. Nordtest" }
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[00_GLOSSARIO]]"
  - "[[03_PROTOCOLLO_HANDOFF_QM]]"
  - "[[02_TRE_LIVELLI_DI_CONTROLLO]]"
  - "[[03_RISPOSTA_NC]]"
  - "[[sez_7_7_assicurazione_validita]]"
  - "[[sez_7_10_attivita_non_conformi]]"
  - "[[sez_8_7_azioni_correttive]]"
  - "[[A1_validazione_metodi]]"
  - "[[A2_incertezza_misura]]"
  - "[[A4_materiali_riferimento_CRM]]"
  - "[[A8_gestione_NC_da_accreditatore]]"
tags: [PT, ILC, prove_valutative, RT-39, z-score, En, piano_PT, copertura_PT]
escalation_qm: true
---

> **Marcatore di livello: 🟢🟡🟣** — Appendice trasversale. Contenuti normativi (livello 1+2), prassi ispettiva (livello 3), best practice tecnica (livello 4).

# Appendice A5 — Prove valutative (PT) e confronti interlaboratorio (ILC)

> **Riassunto operativo.** Le **PT** (Proficiency Testing) e gli **ILC** (Interlaboratory Comparisons) sono uno dei pilastri obbligatori dell'assicurazione di validità (§ 7.7) e sono oggetto di una **prescrizione Accredia dedicata** (`RT-39`). Un piano PT vuoto o opaco è una NC quasi certa in visita Accredia. Un PT "unsatisfactory" non gestito è errore fatale. Questa appendice serve al QM/RT per costruire il piano di partecipazione, scegliere i provider, leggere gli esiti, governare le NC che ne derivano. L'agente la carica quando il QM parla di "piano PT", "z-score", "E_n", "PT scadente", "scelta provider", "copertura PT dello scopo", o quando un esito PT richiede risposta.

---

## 1. Principio normativo (🟢 LIVELLO 1+2)

### 1.1 ISO 17025:2018 § 7.7 — Assicurazione della validità dei risultati

> Testo a pagamento, non riprodotto. Sintesi operativa:

- **§ 7.7.1** — il laboratorio deve avere una procedura per monitorare la validità dei risultati, registrare i dati in modo da permettere l'identificazione di tendenze e applicare tecniche statistiche al riesame. Tra i mezzi elencati: CRM, RM, strumentazione di taratura, repliche, riesame critico, carte di controllo, intra-laboratory comparisons.
- **§ 7.7.2** — il laboratorio deve monitorare le proprie prestazioni mediante confronto con altri laboratori, ove disponibile e appropriato. Ciò include, in modo non esclusivo, una o entrambe le seguenti opzioni:
  - (a) partecipazione a **proficiency testing**;
  - (b) partecipazione a **interlaboratory comparisons** diversi dai PT.
- **§ 7.7.3** — i dati di monitoraggio devono essere analizzati, usati per controllare e dove applicabile migliorare le attività di laboratorio. Se i risultati indicano performance non conformi a criteri predefiniti, il laboratorio adotta azioni appropriate per prevenire la riportata di risultati errati.

### 1.2 ISO 17025:2018 § 7.10 e § 8.7 — Trattamento dell'esito non conforme

- **§ 7.10** definisce come trattare le attività non conformi (e quindi anche un PT con esito unsatisfactory che indica un problema sull'attività).
- **§ 8.7** chiede l'azione correttiva quando la non conformità si manifesta.

### 1.3 RT-08 rev.05 § 7.7

> **Fonte**: `[RT-08 rev.05 § 7.7 — EC 10-02-2022]`

RT-08 § 7.7 si applica integralmente con rinvio specifico a **`RT-39`** per le prescrizioni di partecipazione a PT/ILC. RT-08 enfatizza che:
- la partecipazione a PT/ILC è **obbligatoria** per ciascuna categoria/sottocategoria di prova nel campo di accreditamento, con frequenza e copertura definite da RT-39,
- gli esiti di PT/ILC devono essere riesaminati con criteri statistici espliciti (z-score, E_n, ζ-score per dati con incertezza dichiarata), e gli esiti **questionable** o **unsatisfactory** devono attivare il processo § 7.10 / § 8.7.

### 1.4 RT-39 rev.00 — Prescrizioni Accredia per la partecipazione a PT/ILC

> **Fonte**: `[RT-39 rev.00 (Accredia, free)]`

`RT-39` è il documento di riferimento operativo. Punti chiave (sintesi):
- Il laboratorio definisce un **piano di partecipazione pluriennale** (tipicamente quadriennale o coerente con il ciclo di accreditamento), che copre il proprio scopo per **categoria/sottocategoria** secondo la tassonomia Accredia.
- La **scelta del provider** privilegia provider accreditati **`ISO 17043`**. In assenza di provider accreditato per una specifica categoria, il laboratorio può ricorrere a ILC non accreditati o auto-organizzati, motivando la scelta.
- Per ciascuna categoria coperta, RT-39 indica **frequenze minime di partecipazione** e modalità di copertura per parametri/matrici. Il laboratorio deve dimostrare la coerenza tra piano PT e scope di accreditamento.
- Gli esiti vengono **valutati** con criteri statistici predefiniti. Gli esiti non soddisfacenti devono essere oggetto di analisi causa, azione correttiva, verifica di efficacia, comunicazione (se prevista da contratto/cliente) e documentazione.

> Il QM è tenuto a consultare la versione vigente di RT-39 sul sito Accredia per la corretta applicazione di frequenze e criteri specifici per categoria.

---

## 2. Come un ispettore lo verifica (🟡 LIVELLO 3 — prassi ispettiva)

> Pattern: meccanismo → estensione → efficacia.

### 2.1 Domanda 1 — Meccanismo

**"Mostratemi il vostro piano pluriennale di partecipazione a PT/ILC, con la mappa di copertura dello scopo di accreditamento per categoria/sottocategoria, le frequenze applicate, i provider scelti (con accreditamento ISO 17043) e la procedura di gestione degli esiti."**

Analisi ispettiva:
- Cerco un **documento formale** (procedura + piano) firmato dal QM/RT. Una sequenza di partecipazioni "che capitano" non è un piano.
- Verifico la **copertura** della tassonomia Accredia (RT-39): ogni categoria/sottocategoria nello scope deve essere coperta nel periodo pluriennale.
- Controllo se i **provider** sono accreditati ISO 17043 per gli schemi specifici. Quando non lo sono (ILC auto-organizzato), cerco la **motivazione documentata**.

### 2.2 Domanda 2 — Estensione

**"Per ogni categoria nello scope, quanti PT/ILC avete eseguito negli ultimi 4 anni e con quali esiti? Mostratemi la matrice 'scopo × periodi × esiti'. Esistono categorie 'scoperte' o con frequenza inferiore al minimo RT-39?"**

Analisi ispettiva:
- Campionamento sullo scopo: prendo 3-5 categorie/sottocategorie e chiedo il dettaglio.
- **Errore frequente**: categorie scoperte da anni perché "non ci sono PT disponibili" — senza evidenza di tentativi di ricerca provider o di organizzazione di ILC alternativo. NC su § 7.7 + RT-39.
- **Errore frequente 2**: PT eseguiti "due volte sullo stesso parametro" e "zero volte" su altro parametro della stessa categoria. Frammentazione della copertura. NC potenziale su RT-39.

### 2.3 Domanda 3 — Efficacia

**"Negli ultimi 24 mesi, è arrivato un esito PT 'questionable' (|z|>2) o 'unsatisfactory' (|z|>3, o E_n>1)? Cosa avete fatto? Mi mostrate il caso completo: analisi causa, estensione, impatto sui risultati emessi, azione correttiva, verifica di efficacia."**

Analisi ispettiva:
- Un PT scadente è un **dono dell'ispettore**: l'evento c'è già, manca solo la dimostrazione della reazione del laboratorio. Cerco la **sequenza completa** ex `[[03_RISPOSTA_NC]]`.
- **Errore tipico**: l'esito unsatisfactory è registrato e basta. Nessuna analisi causa, nessuna NC interna aperta, nessuna verifica di impatto sui risultati che il laboratorio ha emesso nello stesso periodo per quella matrice/parametro. NC grave: NON è un esito tecnico, è una mancata gestione del § 7.7 / § 7.10 / § 8.7.
- Se il laboratorio ha tutti PT "satisfactory" da anni con z-score sempre stretti, indaga: o la prestazione è effettivamente eccellente, oppure il provider è poco severo, oppure i dati sono allineati a posteriori (sospetto). L'ispettore chiede ricostruzione.

### 2.4 Pattern di errori comuni

1. **Piano PT inesistente** o coincidente con "elenco partecipazioni": NC su § 7.7 + RT-39.
2. **Copertura incompleta dello scopo**: categorie non coperte, parametri/matrici non rotati. NC su RT-39.
3. **Provider non accreditato ISO 17043** senza motivazione: NC potenziale su RT-39.
4. **PT scadente non gestito**: registrazione c'è, NC interna non aperta, analisi causa assente. NC grave su § 7.7.3 + § 8.7.
5. **PT scadente "chiuso" con sola formazione del personale**: causa = "errore umano", azione = "rifare il corso". Vedi `[[03_RISPOSTA_NC]] Errore 1`. La barriera di sistema non è stata indagata.
6. **Mancata verifica di impatto sui risultati emessi** nel periodo del PT scadente. Errore "domanda dimenticata" (vedi 03_RISPOSTA_NC § Errore 3).
7. **Confusione PT vs ILC**: l'agente e il QM devono distinguere. Un PT formale ha provider accreditato ISO 17043, materiale fornito, valutazione statistica codificata. Un ILC è un confronto interlaboratorio organizzato (anche dal laboratorio stesso): valido ai fini § 7.7 ma con minore riconoscimento "automatico".
8. **PT/ILC e incertezza disallineati**: il laboratorio dichiara U piccola nei rapporti, ma sui PT presenta z-score larghi. Vedi `[[A2_incertezza_misura]] § 2.3`.

---

## 3. Best practice tecnica (🟣 LIVELLO 4 — best practice)

### 3.1 Schemi di valutazione statistica (riferimenti)

| Indicatore | Calcolo (sintesi) | Soglie tipiche | Riferimento |
|---|---|---|---|
| **z-score** | `(x_lab − x_assigned) / σ_pt` | `|z| ≤ 2` satisfactory; `2 < |z| ≤ 3` questionable; `|z| > 3` unsatisfactory | ISO 13528 (citato in ISO 17043) |
| **z'-score** | come z, con σ modificato per includere incertezza del valore assegnato | stesse soglie | ISO 13528 |
| **E_n** | `(x_lab − x_ref) / sqrt(U_lab² + U_ref²)` | `|E_n| ≤ 1` accettabile; `|E_n| > 1` non accettabile | ISO 13528 (per ILC con incertezza dichiarata) |
| **ζ-score** (zeta) | `(x_lab − x_assigned) / sqrt(u_lab² + u_ref²)` | `|ζ| ≤ 2` satisfactory, ecc. | ISO 13528 (usato quando il lab dichiara la propria incertezza) |

> `ISO-17043` (2010 e 2024) recepisce questi criteri e descrive come il provider li applica. Il laboratorio partecipante li **legge** dal rapporto PT, non li ricalcola autonomamente.

### 3.2 Guide internazionali

| Documento | Quando consultarlo |
|---|---|
| `EURACHEM-PT 3rd ed. 2021` | Guida d'elezione su selezione, uso e interpretazione dei PT in chimica. Approfondisce: come scegliere lo schema, come leggere l'esito, come gestire gli unsatisfactory |
| `EURACHEM-MicroAcc 2023` | Specificità del PT in microbiologia (cariche, incertezza di conteggio, valutazione su scala log) |
| `Nordtest-TR-569 (Trollbook) ed. 5.1` | Manuale operativo di controllo qualità interno per laboratori chimici, integra PT con carte di controllo |
| `ISTISAN-12/29 4ª ed. 2011` | Traduzione italiana di Nordtest TR-569; pratica per laboratori IT |

### 3.3 PT vs ILC: come scegliere

| Aspetto | PT formale (provider ISO 17043) | ILC (auto-organizzato o non accreditato) |
|---|---|---|
| Riconoscimento | Pieno per RT-39 quando provider ISO 17043 | Da motivare; ammesso in assenza di PT disponibili |
| Materiale | Fornito dal provider, omogeneità/stabilità garantite | Da definire (chi prepara, come si dimostra omogeneità/stabilità) |
| Valore assegnato | Determinato dal provider (consenso, riferimento, formulazione) | Da definire (media partecipanti, lab di riferimento) |
| Valutazione statistica | Standard (z, E_n, ζ) | Da definire e dichiarare |
| Sforzo organizzativo | Basso per il lab | Alto |
| Costo | Variabile per schema | Spesso minore ma con costo organizzativo |

### 3.4 Frequenze tipiche di partecipazione

> RT-39 dà le frequenze ufficiali. Indicativamente, su base **pluriennale** (4 anni), molti laboratori organizzano:
> - 1 PT all'anno per categoria principale,
> - rotazione di matrici/parametri all'interno della categoria,
> - aumento di frequenza in categorie critiche (richieste cogente, esiti storici borderline, metodi nuovi in onboarding).

### 3.5 Lettura dell'esito — checklist tecnica

- [ ] Il valore assegnato (x_assigned) è stato determinato come? (consenso partecipanti, valore di riferimento, formulazione gravimetrica)
- [ ] Quale σ è stato usato? È coerente con quello atteso per il metodo (es. funzione Horwitz per chimica analitica, o σ "prescritto" dal provider)?
- [ ] L'incertezza dichiarata dal lab (se richiesta) era congruente con quella della popolazione partecipante?
- [ ] Il valore del lab è ragionevolmente coerente con la deviazione tipica del proprio metodo dichiarata in validazione (vedi `[[A1_validazione_metodi]]`)?
- [ ] Eventuali commenti del provider (outlier, esclusioni, basso numero partecipanti) sono stati letti e annotati?

### 3.6 Quando il PT non c'è (RT-39 § "ILC alternativi")

Se per una categoria specifica del proprio scope non esistono provider PT accreditati ISO 17043 (caso reale in settori molto specialistici), il laboratorio:
1. **Documenta la ricerca** (provider consultati, esiti negativi, date),
2. **Pianifica un'alternativa**: ILC con altri laboratori accreditati, confronto con CRM su matrice simile, partecipazione a PT generalisti pur fuori categoria stretta, confronto su campione split con altro laboratorio,
3. **Motiva** la scelta nel piano PT,
4. **Riconsidera periodicamente** la disponibilità di provider (annualmente).

---

## 4. Template operativo per il QM

### 4.1 Piano pluriennale PT/ILC — struttura minima

| Anno | Categoria/sottocategoria Accredia | Parametri / matrici coperti | Provider previsto + accred. ISO 17043 | Frequenza prevista (n. PT/anno) | Note (rotazione, nuovo metodo, esito precedente) |
|---|---|---|---|---|---|

### 4.2 Registro esiti PT/ILC — struttura minima

| Anno | Schema PT (provider, codice) | Categoria/sottocategoria coperta | Matrice/parametro | Valore lab | x_assigned + σ | z-score / E_n / ζ | Esito (S/Q/U) | NC interna aperta? (rif.) | Stato chiusura |
|---|---|---|---|---|---|---|---|---|---|

### 4.3 Workflow "Esito PT non soddisfacente"

```
1. Registrazione esito nel Registro PT con codice univoco
2. Apertura immediata NC interna (rif. § 7.10) con codice e responsabile
3. Analisi di causa (rif. § 8.7) NON fermarsi a "errore umano"
   Vedi [[03_RISPOSTA_NC]] Errore 1
4. Analisi di estensione:
   - lo stesso problema ricorre in altri operatori / strumenti / matrici?
   - i risultati emessi nel periodo del PT scadente sono difendibili?
5. Verifica di impatto sui risultati emessi (rif. § 7.10.2):
   - identificazione rapporti potenzialmente impattati
   - rivalutazione tecnica (RT)
   - eventuale comunicazione cliente
6. Azione correttiva:
   - tecnica (es. ri-validazione parametro, ri-taratura strumento, modifica procedura)
   - sistemica se la barriera CQ non ha intercettato (modifica del CQ stesso)
7. Verifica di efficacia su casi reali (NON sulla sola formazione):
   - successivo PT sul medesimo parametro
   - prove di ricontrollo su campioni storici / spike / CRM
   - audit mirato
8. Documentazione completa e chiusura formale della NC da parte del QM
```

### 4.4 Checklist "Verifica copertura PT vs scope" (audit interno)

- [ ] Esiste mappa scope di accreditamento × categorie RT-39
- [ ] Per ogni categoria nello scope esiste almeno una partecipazione PT/ILC nell'ultimo periodo previsto da RT-39
- [ ] Per categorie con più matrici/parametri, è documentata la rotazione
- [ ] I provider scelti sono accreditati ISO 17043 (o motivata l'eccezione)
- [ ] Tutti gli esiti questionable/unsatisfactory degli ultimi 24-36 mesi sono stati gestiti come NC (rif. univoco)
- [ ] Le azioni correttive su PT scadenti sono state verificate per efficacia
- [ ] Il piano pluriennale è stato riesaminato nell'ultimo riesame di direzione (vedi `[[A6_audit_interni_riesame_direzione]]`)

---

## 5. Quando l'agente attiva questa appendice

1. Il QM chiede *"come imposto il piano PT?"*, *"questo PT è sufficiente?"*, *"il provider è OK?"*.
2. Arriva un esito PT con z-score borderline o unsatisfactory: serve workflow ex § 7.10 (attività non conformi).
3. È in lavorazione il riesame del piano PT (annuale o quadriennale) e serve verifica copertura.
4. È in onboarding un nuovo metodo o si estende lo scope: il piano PT va aggiornato.
5. È in lavorazione una NC ricevuta da Accredia che cita § 7.7 / RT-39 (vedi `[[A8_gestione_NC_da_accreditatore]]`).
6. Il riesame di direzione (vedi `[[A6_audit_interni_riesame_direzione]]`) chiede l'analisi degli esiti PT.

---

## 6. Cosa l'agente NON può chiudere su questo tema (handoff QM obbligatorio)

Vedi `[[03_PROTOCOLLO_HANDOFF_QM]] § 4`. In materia di PT/ILC, l'agente **non** può:

1. **Approvare il piano pluriennale PT** come adeguato. Propone bozza/critica; approvazione del QM, validazione tecnica del RT.
2. **Dichiarare un esito PT "chiuso"** anche se la documentazione sembra completa. La chiusura della NC interna conseguente è del QM (rif. `[[03_PROTOCOLLO_HANDOFF_QM]] § 4.1-4.2`).
3. **Dichiarare un PT scadente come "accettabile per cause esterne"** (es. provider con problema riconosciuto). Anche se documentato dal provider, la decisione di non aprire NC interna è del QM.
4. **Decidere di non partecipare** a un PT per una categoria. Riduzione della frequenza/copertura va motivata e formalizzata dal QM/RT, mai dall'agente.
5. **Selezionare il provider**. L'agente può istruire la valutazione comparativa (accreditamento ISO 17043, esperienza, costo, lingua, tempi); la scelta è gestionale.
6. **Decidere comunicazione clienti** per risultati potenzialmente impattati da PT scadente. Decisione di Direzione + QM + RT.
7. **Dichiarare chiusa una NC su § 7.7 / RT-39** o azione correttiva su PT scadente. La chiusura è del QM, l'efficacia è giudizio del RT.
