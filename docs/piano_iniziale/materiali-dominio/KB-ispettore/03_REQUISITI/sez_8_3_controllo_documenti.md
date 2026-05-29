---
codice: sez_8_3_controllo_documenti
tipo: requisito_normativo
livello: misto
titolo: "§ 8.3 — Controllo dei documenti del sistema di gestione (Opzione A)"
sezione_norma: "8.3"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "8.3" }
  - { doc: "RT-08", rev: "05", paragrafo: "8.3" }
fonti_secondarie:
  - { doc: "LS-04", rev: "20" }
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[02_TRE_LIVELLI_DI_CONTROLLO]]"
  - "[[sez_8_2_documentazione_sgq]]"
  - "[[sez_8_4_controllo_registrazioni]]"
tags: [controllo_documenti, documenti_esterni, versionamento, distribuzione, 3_mesi]
escalation_qm: true
---

# § 8.3 — Controllo dei documenti del sistema di gestione (Opzione A)

> Il § 8.3 disciplina come il laboratorio gestisce i documenti del SGQ (interni ed esterni): approvazione, revisione, identificazione delle modifiche, accessibilità, prevenzione dell'uso non intenzionale di documenti obsoleti.

> ⚠️ **Disclaimer di ricostruzione KB**
> Questo file fa parte della ricostruzione della **sezione 8** sull'indice ufficiale **ISO/IEC 17025:2018**. La vecchia KB aveva numerazione errata mutuata da ISO 9001.

---

## A. Testo della norma (🟢 LIVELLO 1 — fatto normativo)

> **Fonte**: `[ISO-17025-2018 § 8.3]` — testo a pagamento, non riprodotto integralmente.

**Sintesi operativa per sotto-paragrafo**:

- **8.3.1** — Il laboratorio deve controllare i documenti (interni ed esterni) relativi al rispetto della presente norma.
  - **Nota**: i documenti possono essere in qualsiasi formato/supporto (cartaceo, elettronico).
- **8.3.2** — Il laboratorio deve assicurare che:
  - a) i documenti siano approvati per adeguatezza prima dell'emissione da personale autorizzato;
  - b) i documenti siano periodicamente riesaminati e aggiornati quando necessario;
  - c) le modifiche e lo stato di revisione corrente siano identificati;
  - d) le versioni pertinenti dei documenti applicabili siano disponibili nei luoghi di utilizzo;
  - e) i documenti siano identificati in modo univoco;
  - f) sia impedito l'uso non intenzionale di documenti obsoleti e che essi siano opportunamente identificati se conservati per scopi qualsiasi.

---

## B. Prescrizioni Accredia (🟢 LIVELLO 2 — prescrizione Accredia)

> **Fonte**: `[RT-08 rev.05 § 8.3]`

### 8.3.1 (RT-08)
> *"Si applica il requisito di norma. Nel caso di aggiornamenti di documenti di origine esterna (es. norme, metodi, leggi, regolamenti), ove non diversamente indicato, il Laboratorio è tenuto ad applicare le nuove versioni entro tre mesi dall'emissione."*

### 8.3.2 (RT-08)
> *"Si applica il requisito di norma. Correzioni a mano sono permesse sui documenti a solo uso interno, qualora se ne ravvisi l'urgenza. Tali documenti modificati andranno tempestivamente aggiornati."*

---

### 🔑🔑🔑 LA REGOLA DEI 3 MESI — prescrizione FORTE Accredia

> **TESTO RT-08**: *"Nel caso di aggiornamenti di documenti di origine esterna (es. norme, metodi, leggi, regolamenti), ove non diversamente indicato, il Laboratorio è tenuto ad applicare le nuove versioni entro tre mesi dall'emissione."*

**Implicazioni pratiche**:

- Il countdown parte dalla **data di emissione** della nuova versione (es. data di pubblicazione UNI, ISO, Gazzetta Ufficiale), NON dalla data in cui il laboratorio se ne accorge.
- "Applicare la nuova versione" significa:
  1. **Approvare** internamente la nuova versione del documento esterno;
  2. **Aggiornare procedure/istruzioni** interne che ne derivano;
  3. **Formare** il personale sulla nuova versione;
  4. **Modificare** moduli/registrazioni se richiesto;
  5. **Riconvalidare** il metodo se la nuova versione comporta cambi tecnici significativi (vedi § 7.2.2).
- Se la nuova versione comporta una **rivalidazione**, 3 mesi possono non essere sufficienti: il laboratorio deve **documentare il piano transitorio** e comunicarlo eventualmente al cliente.
- Eccezione: *"ove non diversamente indicato"* — alcune norme/leggi prevedono periodi transitori più lunghi (es. coesistenza vecchia/nuova versione per 6-12 mesi). In tal caso il laboratorio segue il transitorio normativo.

**Esempio pratico**:
- UNI EN ISO XXX rev.2026 pubblicata il 15 marzo 2026 → il laboratorio deve avere applicato la nuova versione entro il **15 giugno 2026**.
- Se applica vecchia versione il 20 giugno 2026 senza giustificazione → **NC**.

---

### 🔑 Correzioni a mano (RT-08 § 8.3.2)

- **Solo su documenti a uso interno** (no rapporti di prova, no documenti distribuiti a clienti).
- **Solo se urgenza** dimostrabile (es. errore di battitura su istruzione operativa scoperto in corso d'opera).
- La correzione deve essere:
  - sbarrata sul vecchio testo (mai cancellata);
  - controsiglata (firma + data) da chi corregge;
  - se rilevante, controfirmata dal QM o dal responsabile autorizzato.
- **Va aggiornata tempestivamente** una versione formale del documento (non si lascia la correzione a mano in eterno).

---

### 🔑 Documenti esterni tipici sotto controllo § 8.3

| Tipo | Esempi | Frequenza riesame |
|---|---|---|
| Norme tecniche | UNI, ISO, EN, ASTM, EPA, OECD | Annuale + alert pubblicazione |
| Metodi ufficiali | DM, Allegati GU, Farmacopea, AOAC, EPA | Annuale + alert pubblicazione |
| Leggi e regolamenti | D.Lgs., Direttive UE, Regolamenti UE | Continuo (alert normativo) |
| Documenti Accredia | RT, RG, LS, Circolari | Alla pubblicazione (sito Accredia) |
| Documenti EA/ILAC | Guide EA, ILAC P-series, G-series | Alla pubblicazione |
| Schede tecniche fornitori | SDS reagenti, manuali strumenti, certificati MR | Alla ricezione + ad ogni nuovo lotto |

---

### 🔑 Identificazione univoca dei documenti

Codice tipico: `TIPO-AREA-NN rev.YY del DD/MM/AAAA`
- TIPO: PG (procedura gestionale), PO (procedura operativa), IO (istruzione operativa), MO (modulo), MQ (manuale qualità)
- AREA: 00 (gestione), 01 (chimica), 02 (microbiologia), ecc.
- NN: progressivo
- rev.YY: revisione
- Data emissione

---

## C. Documenti applicabili (LS-04 rev.20)

- **RG-02** — Regolamento generale accreditamento laboratori
- **RT-08 rev.05** — Prescrizioni Accredia
- **PG-XX** — Procedura interna gestione documenti (deve esistere nel SGQ del lab)
- Elenco documenti SGQ (master list)
- Elenco documenti esterni controllati

---

## D. Domande ispettive (🟡 LIVELLO 3 — prassi ispettiva)

### 8.3.1 — Documenti esterni e regola dei 3 mesi

#### Domanda 1 — Meccanismo
*"Come venite a conoscenza dell'emissione di una nuova versione di una norma che usate?"*
- Risposte attese: abbonamento UNI (Edilex, UNIstore), alert ISO, mail Accredia, mailing list settoriale, audit normativo periodico.

#### Domanda 2 — Estensione
*"Mostratemi l'elenco dei documenti esterni controllati. Per ognuno, qual è la versione corrente in uso?"*
- Verificare a campione 3-5 norme: confrontare data di pubblicazione vs data di adozione interna. Se >3 mesi senza giustificazione → **NC**.

#### Domanda 3 — Efficacia
*"Mi mostrate l'ultimo caso in cui una norma è stata aggiornata e voi avete dovuto recepire? Quale norma? Quando emessa? Quando recepita? Quale impatto sulle vostre procedure? Quale formazione del personale?"*
- Trabocchetto: chiedere PROVA documentale (verbale formazione, riemissione procedura, evidenza rivalidazione).

### 8.3.2 — Controllo documenti interni

#### Domanda 1 — Meccanismo
*"Chi approva un nuovo documento del SGQ? Chi può revisionarlo? Chi lo distribuisce? Dove è documentato il flusso?"*

#### Domanda 2 — Estensione
*"Esistono correzioni a mano sui documenti in uso? Mostratemi i casi recenti. Quando saranno riemessi?"*

#### Domanda 3 — Efficacia
*[Test sul campo]: vado in laboratorio, prendo una procedura dallo scaffale, verifico che la versione affissa = versione master list. Se obsoleta → **NC**.*

---

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"L'agente ha rilevato: [a] documento esterno non recepito entro 3 mesi (norma XXX rev.YYYY, emessa GG/MM, in uso ancora versione precedente); [b] correzione a mano non aggiornata dopo X settimane; [c] documento obsoleto reperibile nei punti di utilizzo."*
- Riferimento: ISO § 8.3 + RT-08 § 8.3
- Urgenza: ALTA se documento esterno obsoleto impatta su metodi accreditati; MEDIA altrimenti
- Cosa NON può chiudere agente: emissione/approvazione documenti, decisione su rivalidazione metodo

---

## E. Errori fatali correlati

1. **Norma esterna aggiornata da >3 mesi e non recepita** senza giustificazione → **NC maggiore** se impatta metodo accreditato.
2. **Documento obsoleto in uso** nei punti operativi (es. istruzione vecchia affissa accanto allo strumento) → **NC**.
3. **Correzioni a mano sul rapporto di prova** o su documenti destinati al cliente → **NC grave**.
4. **Mancanza di master list documentale** o master list non aggiornata → **NC**.
5. **Documenti non identificati univocamente** (no codice, no revisione, no data) → **NC**.
6. **Approvazione documento da personale non autorizzato** (es. tecnico che approva una procedura gestionale) → **NC**.
7. **Recepimento "formale" della norma ma senza modifica delle procedure derivate** o senza formazione → recepimento solo apparente, **NC** in sostanza.

---

## F. Quando l'agente DEVE chiedere prima di rispondere

- Se il laboratorio chiede di interpretare l'applicabilità del termine "3 mesi" a un caso specifico (es. norma con transitorio) → **escalation a QM**.
- Se chiede di decidere se una nuova versione di norma richieda rivalidazione del metodo → **escalation a QM + responsabile tecnico**.
- Se chiede di approvare/distribuire un documento al posto del personale autorizzato → **escalation a QM**, l'agente non ha autorità di approvazione.
- Se chiede di gestire una correzione a mano su documento destinato al cliente → **rifiutare**, va emessa una nuova versione formale.
