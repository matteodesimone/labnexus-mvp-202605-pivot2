---
codice: sez_7_5_registrazioni_tecniche
tipo: requisito_normativo
livello: misto
titolo: "§ 7.5 — Registrazioni tecniche"
sezione_norma: "7.5"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.5" }
  - { doc: "RT-08", rev: "05", paragrafo: "7.5" }
fonti_secondarie:
  - { doc: "RT-08", rev: "05", paragrafo: "8.4 (tempi conservazione registrazioni)" }
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[02_TRE_LIVELLI_DI_CONTROLLO]]"
  - "[[sez_7_4_manipolazione_oggetti]]"
  - "[[sez_7_8_presentazione_risultati]]"
  - "[[sez_7_11_controllo_dati]]"
tags: [registrazioni_tecniche, ricostruzione_tecnica, integrita_dati, correzioni, audit_trail]
escalation_qm: true
---

# § 7.5 — Registrazioni tecniche

> Le registrazioni tecniche sono ciò che permette di **ricostruire** una prova a distanza di mesi/anni. La regola guida: chiunque, con le registrazioni in mano, deve poter riprodurre la prova e arrivare allo stesso risultato del RdP. È il punto in cui l'ispettore mette alla prova la sostanza della 17025: si pesca un RdP a caso e si chiede "ricostruitemi tutto, dal ricevimento al risultato".

## A. Testo della norma (🟢 LIVELLO 1 — fatto normativo)

> **Fonte**: `[ISO-17025-2018 § 7.5.1 - 7.5.2]` — testo a pagamento, non riprodotto. Riassunto operativo:

- **7.5.1** Il laboratorio deve assicurare che le registrazioni tecniche di ciascuna attività di laboratorio contengano i **risultati, il rapporto e informazioni sufficienti a facilitare**, se possibile, l'**identificazione dei fattori che influenzano il risultato e la sua incertezza** e che consentano la **ripetizione dell'attività in condizioni il più possibile vicine all'originale**. Le registrazioni devono includere data e identità del personale responsabile di ciascuna attività e del controllo dei dati e dei risultati. Le osservazioni originali, i dati e i calcoli devono essere registrati al momento in cui sono effettuati e devono essere identificabili rispetto alla specifica attività.
- **7.5.2** Il laboratorio deve garantire che le **modifiche** alle registrazioni tecniche possano essere **tracciate alle versioni precedenti** o alle osservazioni originali. Sia i dati originali sia i dati modificati devono essere conservati, includendo data, indicazione degli aspetti alterati e personale responsabile delle modifiche.

## B. Prescrizioni Accredia (🟢 LIVELLO 2 — prescrizione Accredia)

> **Fonte**: `[RT-08 rev.05 § 7.5]`

### 7.5.1 — Tempi di conservazione

> **Citazione integrale RT-08**:
> *"Per i tempi di conservazione delle registrazioni tecniche vedere il § 8.4 del presente documento."*

Quindi il tempo di conservazione è quello generale delle registrazioni di SGQ del laboratorio (vedi file dedicato sez_8_4). In genere il laboratorio fissa contrattualmente o in procedura un tempo coerente con prescrizioni cogenti (es. controllo ufficiale alimenti, accise, ambientale).

### 7.5.2 — Correzioni alle registrazioni

> **Citazione integrale RT-08**:
> *"In caso di correzione di dati, ove dalle registrazioni non fosse desumibile la spiegazione, deve essere annotato il motivo della correzione."*

In pratica: oltre alla tracciabilità chi/cosa/quando della modifica (norma), Accredia chiede anche il **perché** se non è desumibile dal contesto.

## C. Documenti applicabili (LS-04 rev.20)

- `ISO-17025-2018` § 7.5 — norma base.
- `RT-08 rev.05` § 7.5 e § 8.4 (tempi conservazione registrazioni SGQ).
- Norme settoriali con prescrizioni specifiche su tempi di conservazione (es. residui veterinari, ambientale, accise).

## D. Domande ispettive (🟡 LIVELLO 3 — prassi ispettiva)

> Pattern: meccanismo → estensione → efficacia.
> Il test ispettivo principale qui è la **ricostruzione tecnica**: si prende un RdP e si chiede di riprodurre tutto il processo dai documenti tecnici.

### 7.5.1 — Contenuto delle registrazioni tecniche

#### Domanda 1 — Meccanismo
**"Quali registrazioni tecniche generate per ogni prova? Le elencatemi: ricevimento, preparazione, lettura strumentale, calcolo, controlli di qualità, validazione del dato, approvazione. Per ognuna, mi dite chi la firma, su che supporto (cartaceo/elettronico/LIMS), e con quale identificazione del responsabile."**

Analisi ispettiva: la lista deve coprire l'intero ciclo. Per i lab digitalizzati (LIMS), serve audit trail su ogni passaggio. Tipici problemi: registrazioni intermedie (es. integrazione cromatogrammi) non firmate o solo "by initials" senza tracciabilità identità.

#### Domanda 2 — Estensione (test di ricostruzione tecnica — IL TEST CENTRALE)
**"Scelgo io un RdP dell'ultimo trimestre. Ricostruite tutte le registrazioni tecniche corrispondenti: ricevimento, sub-campioni, condizioni ambientali, taratura/verifica strumenti usati, lotti dei reagenti, controlli di qualità (deve avere accettato), calcoli, audit trail eventuali modifiche, approvazione finale. Voglio poter ripetere mentalmente la prova."**

Analisi ispettiva: questo è IL test del § 7.5.1. Se la ricostruzione si blocca a un punto qualsiasi (es. "non troviamo il certificato di taratura dello strumento usato quel giorno"), c'è un buco nelle registrazioni. Cose da verificare:
- Identità di chi ha eseguito ogni fase e di chi ha approvato;
- Data/ora delle operazioni effettuate "al momento";
- Lotti reattivi e materiali di riferimento usati;
- Tarature/verifiche strumenti applicabili nella data della prova;
- Controlli di QC interni (collega § 7.7) applicati nella sessione e loro esito;
- Calcoli (formula, sostituzioni, risultato);
- Stima incertezza (collega § 7.6);
- Eventuali modifiche alle registrazioni (cosa, perché, chi, quando — collega § 7.5.2).

Questo test è il vero "esame del lab". Se passa, gran parte del § 7.5 è ok. Se fallisce, c'è NC e si rivede l'intera gestione delle registrazioni.

#### Domanda 3 — Efficacia
**"Per un metodo che usa fogli elettronici (Excel) per i calcoli: il foglio è documentato, validato e protetto come prescritto dal RT-08 § 7.11.6? E i risultati intermedi vengono salvati / sono ricostruibili?"**

Analisi ispettiva: collegamento esplicito con § 7.11.6 — fogli elettronici devono essere documentati, validati, protetti contro alterazione involontaria. Spesso questo è un buco enorme: file Excel su cartelle condivise, formule senza protezione, "salvataggio sopra" che cancella la versione precedente.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Test di ricostruzione tecnica su RdP recente: completezza registrazioni dal ricevimento al risultato, tracciabilità responsabili e data/ora, gestione fogli di calcolo. ISO § 7.5.1 + RT-08 § 7.11.6."*
- Urgenza: massima — è il test sostanziale della 17025.
- Cosa NON può chiudere l'agente: NON può dichiarare la registrazione conforme se anche un solo anello della catena di ricostruzione si rompe. Espone al QM.

### 7.5.2 — Modifiche e correzioni alle registrazioni

#### Domanda 1 — Meccanismo
**"Come gestite le correzioni alle registrazioni tecniche? Per cartaceo: come si correggono i dati senza coprirli? Per elettronico: avete audit trail che mantiene il dato originale e mostra chi/quando/cosa/perché?"**

Analisi ispettiva: cartaceo → riga singola, nuovo dato accanto, firma, data, motivo. Elettronico → audit trail nativo, NON sovrascrittura. Il RT-08 chiede esplicitamente il **motivo** quando non desumibile dal contesto.

#### Domanda 2 — Estensione
**"Mostratemi 3 esempi reali di correzioni recenti: voglio vedere cosa è stato corretto, da chi, quando, con quale motivazione."**

Analisi ispettiva: cerco esempi concreti. Se il lab dice "non capita mai" su grandi volumi, poco credibile. Se ci sono correzioni ma manca la motivazione (RT-08 esplicito), NC.

#### Domanda 3 — Efficacia
**"Per il LIMS / software di gestione dati: l'audit trail è attivato di default su tutti i campi tecnici? Chi può disattivarlo? Mostratemi un report di audit trail su un dato modificato recentemente."**

Analisi ispettiva: l'audit trail attivato è la base. Se può essere disattivato (anche solo dall'amministratore IT), c'è rischio di alterazione non tracciata. Collegamento con § 7.11 (controllo dati).

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Gestione correzioni registrazioni: cartaceo + elettronico, audit trail LIMS, motivazione delle correzioni (RT-08 esplicito). Esempi reali. ISO § 7.5.2 + RT-08 § 7.5.2 e § 7.11."*
- Urgenza: alta — l'alterazione non tracciata è un problema di integrità dati che invalida tutto il sistema.

### Tempi di conservazione

#### Domanda 1 — Meccanismo
**"Per quanto conservate le registrazioni tecniche? Quale base normativa o contrattuale usate? La regola è uniforme o varia per cliente / metodo / cogenza?"**

Analisi ispettiva: per molti settori la cogenza impone tempi minimi (es. controllo ufficiale alimenti, ambientale). Per gli altri vale il contratto (collega § 7.1.1). Tipico problema: regola interna "5 anni" ma cliente in capitolato richiede 10 anni e nessuno si è accorto.

#### Domanda 2 — Estensione
**"Mostratemi le registrazioni di una prova eseguita 4-5 anni fa: sono ancora reperibili e leggibili?"**

Analisi ispettiva: testo della conservazione. Per supporti elettronici verificare leggibilità (es. file in formato proprietario di software dismesso). Per cartaceo verificare condizioni ambientali archivio.

#### Domanda 3 — Efficacia
**"Cosa fate quando un cliente vi chiede una registrazione 3 anni dopo l'emissione del RdP?"**

Analisi ispettiva: tempo di risposta, processo di recupero documentato, sicurezza dell'invio (riservatezza § 4.2).

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Tempi conservazione registrazioni: base normativa + contrattuale, coerenza tra procedura e contratti, reperibilità di registrazioni vecchie. ISO § 7.5 + RT-08 § 7.5 e § 8.4."*

## E. Errori fatali correlati

- **EF-7.5.1-A**: ricostruzione tecnica fallita — anche un solo anello mancante (taratura applicabile, lotto reagente, QC, audit trail modifica) è NC.
- **EF-7.5.1-B**: registrazioni "by initials" senza tracciabilità univoca dell'identità (chi è "MR"? in 3 lab diversi del gruppo ci sono 4 MR).
- **EF-7.5.1-C**: dati non registrati al momento dell'operazione (es. cromatogrammi integrati il giorno dopo e firmati con data successiva senza spiegazione).
- **EF-7.5.1-D**: fogli Excel non documentati, non validati, non protetti — collega § 7.11.6.
- **EF-7.5.2-A**: correzioni senza motivazione (quando non desumibile dal contesto) — NC esplicita RT-08.
- **EF-7.5.2-B**: dati cartacei corretti con "white" o cancellati, originale non più leggibile — NC.
- **EF-7.5.2-C**: LIMS senza audit trail o con audit trail disattivabile a livello utente — NC grave, è alterabilità del dato.
- **EF-7.5-A**: registrazioni perse o non più reperibili entro il tempo contrattuale o cogente — NC grave (può anche impattare obblighi giuridici esterni).

Vedi anche [[04_ERRORI_FATALI]].

## F. Quando l'agente DEVE chiedere prima di rispondere

L'agente **non risponde** e **chiede al QM** se:

1. Il QM chiede se una specifica registrazione (es. integrazione manuale cromatogramma, foto, tabella raw data) è "tecnica" o "amministrativa": la classificazione cambia il regime di conservazione e di tracciabilità.
2. Si chiede di correggere una registrazione tecnica dopo l'emissione del RdP: collega § 7.5.2 + § 7.10 (NC) + § 7.8.8 (correzioni rapporti).
3. Si chiede di "ricostruire" una registrazione mancante (es. ricomporre da memoria): vietato. NON si ricostruisce, si dichiara la NC e si valuta impatto sul RdP.
4. Audit trail del LIMS o software non completo / non attivato in passato: serve valutazione retroattiva di impatto sui RdP già emessi.
5. Migrazione di sistema (es. nuovo LIMS): la conservazione delle registrazioni nello storico va validata (collega § 7.11.2).
6. Backup persi o corrotti: NC grave, valutazione impatto, comunicazione clienti.
7. Cliente chiede accesso a registrazioni tecniche oltre il RdP: profilo riservatezza (§ 4.2) + contratto + diritto del cliente.

L'agente NON dichiara la ricostruzione "completa" se anche un solo passaggio non è ricostruibile. NON accetta correzioni senza motivazione esplicita. NON dichiara conforme un foglio elettronico non validato/protetto.
