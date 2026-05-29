---
codice: 04_ERRORI_FATALI
tipo: logica_ispettiva
livello: 3
titolo: "Errori fatali da evitare in un laboratorio accreditato"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
autore_originale: "Denis Brazzo — Ispettore Tecnico Accredia DL0871"
file_collegati:
  - "[[01_COME_PENSA_ISPETTORE]]"
  - "[[03_RISPOSTA_NC]]"
  - "[[sez_7_2_selezione_verifica_validazione_metodi]]"
  - "[[sez_6_5_riferibilita]]"
  - "[[sez_7_5_registrazioni_tecniche]]"
nota_livello: "LIVELLO 3 — prassi ispettiva, alert. Elenca errori che, se osservati da un ispettore Accredia, generano tipicamente NC grave o sospensione dell'accreditamento. L'agente li usa per identificare situazioni critiche e proporre alert prioritari al QM con urgenza alta o critica (vedi [[03_PROTOCOLLO_HANDOFF_QM]])."
---

> **Marcatore di livello: 🟡 PRASSI ISPETTIVA (livello 3) con ALERT** — Gli errori descritti qui non hanno attenuanti in audit. Quando l'agente li rileva nel SGQ del laboratorio, alza immediatamente l'urgenza dell'handoff QM ad **alta** o **critica** secondo `[[03_PROTOCOLLO_HANDOFF_QM]]` § 3. Ogni errore è collegato al/ai § di norma e RT-08 che vengono violati.


# Errori fatali da evitare

> L'accreditamento non certifica la bellezza del SGQ. Dimostra la competenza tecnica del laboratorio a produrre risultati validi nello scopo accreditato. Quando questa fiducia viene meno, la risposta ispettiva diventa inevitabilmente severa.

Gli errori descritti in questo documento non sono carenze documentali. Sono situazioni in cui il laboratorio non è più in grado di dimostrare che il dato prodotto sia tecnicamente valido, riferibile, ricostruibile e rilasciato entro il perimetro autorizzato.

---

## 1. Errori che portano quasi certamente a NC maggiore o sospensione

### Attività fuori scopo presentata come accreditata

Emissione di rapporti di prova con riferimento all'accreditamento per attività non coperte dallo scopo, o eseguite con metodo, matrice, misurando, campo di prova o sede non coerenti con l'accreditamento. Il cliente riceve un risultato "coperto" da una fiducia che il laboratorio non aveva titolo a spendere.

> [!ALERT] Questo errore non ammette attenuanti. Genera immediatamente rilievo grave.

### Metodo non verificato, non validato o modificato senza valutazione

Uso di un metodo non verificato, non validato o modificato in modo sostanziale senza adeguata valutazione delle prestazioni, con risultati già emessi. Il problema non è la carta della validazione: è che non si sa se il metodo, così applicato, sia idoneo allo scopo.

### Riferibilità metrologica interrotta o dotazioni critiche fuori stato

Strumenti non tarati, tarature scadute, tarature fuori criterio non valutate retrospettivamente, campioni di riferimento non riferibili, materiali di riferimento o standard scaduti usati in prove accreditate senza rivalutazione tecnica. Se una dotazione condiziona direttamente il risultato e il laboratorio non sa dimostrare lo stato valido al momento della prova, il risultato diventa vulnerabile in modo sostanziale.

> [!ATTENZIONE] Non basta che la taratura "ci sia". Bisogna dimostrare che era valida al momento della prova e che lo strumento era idoneo all'uso previsto.

### Controlli qualità ignorati o gestiti senza razionale tecnico documentato

Sedute analitiche accettate con bianchi contaminati, recuperi fuori criterio, curve non valide, duplicati non conformi, materiali di controllo fuori limite, PT negativi ripetuti senza analisi, carte di controllo con derive evidenti non valutate. Se il risultato è stato emesso come se nulla fosse, la criticità diventa sistemica.

La barriera del controllo qualità deve funzionare. Se intercetta l'anomalia ma non produce una decisione tecnica documentata, non ha protetto il sistema.

### Registrazioni tecniche non ricostruibili

Dati grezzi mancanti, cromatogrammi non collegabili al campione, fogli di lavoro compilati a posteriori, identificazione incompleta di campione/operatore/strumento/materiali, correzioni non tracciate, file sovrascritti, audit trail assente. Quando le registrazioni impediscono di ricostruire come il risultato è stato generato, chi lo ha generato, con quali strumenti e con quali controlli — la criticità non è più minore.

### Personale non competente o non autorizzato su attività critiche

Prove eseguite, riesaminate o autorizzate da personale non qualificato, soprattutto su validazione/verifica metodi, analisi risultati, dichiarazioni di conformità, autorizzazione rapporti. La qualifica formale perde valore se il tecnico non sa valutare uno scostamento o non conosce i criteri di accettabilità del metodo.

### Gestione impropria del campione con perdita di identità o integrità

Campioni scambiati, non identificabili, conservati fuori condizioni critiche, accettati pur non idonei senza valutazione tecnica, aliquote non rintracciabili, tempi di conservazione non controllati. Se il laboratorio non può dimostrare che il risultato appartiene a quel campione e che quel campione era idoneo alla prova — il sistema può essere perfetto, ma il risultato non è difendibile.

### Mancata valutazione dell'impatto sui risultati già emessi

Scoperta di un problema tecnico (taratura fuori criterio, metodo applicato male, controllo non conforme, errore di calcolo) seguito solo da correzione "da oggi in avanti" senza analisi retrospettiva. Questo è spesso il passaggio che trasforma una carenza tecnica in NC maggiore: non solo è accaduto un errore, ma il laboratorio non ha governato le conseguenze.

> [!ALERT] La prima domanda non è "come correggo la procedura?". È "quali risultati potrebbero essere stati influenzati?".

---

## 2. Errori che il laboratorio commette credendo di essere nel giusto

### La mentalità del "risultato plausibile"

Sanare a posteriori una criticità tecnica perché il risultato finale sembra comunque plausibile o perché i controlli generali della seduta sono rientrati. Frasi tipiche: "non ha mai dato problemi", "i PT sono sempre andati bene", "il risultato era coerente", "il responsabile tecnico lo sapeva", "abbiamo esperienza su questa matrice".

Queste frasi, da sole, non dimostrano nulla. Fare una cosa da anni non la rende automaticamente corretta. Una prassi è difendibile solo se è coerente con il metodo, formalizzata quando necessario, tecnicamente giustificata, registrata in modo ricostruibile e monitorata attraverso evidenze.

> [!ISPETTORE] L'esperienza del laboratorio va rispettata. Ma va anche trasformata in controllo tecnico documentato. Se una prassi non è formalizzata, verificata e collegata alle prestazioni del metodo, può diventare un problema grave proprio perché sistematico.

### Deviazioni operative non riconosciute come tali

Modifica "minore" del metodo normalizzato senza valutazione formale: tempo di estrazione diverso, temperatura non rispettata, diluizione gestita fuori schema, curva accettata con punto anomalo, campione conservato oltre i tempi, standard usato oltre validità, taratura non eseguita nei tempi. Il laboratorio pensa di fare un adattamento intelligente; l'ispettore vede una modifica non governata del metodo.

### Spiegazioni retrospettive per registrazioni insufficienti

"Lo sappiamo perché era sicuramente quella bilancia", "quel lotto era l'unico in uso", "l'operatore quel giorno era lui". Il laboratorio pensa di chiarire; in realtà sta mostrando che la ricostruzione dipende dalla memoria delle persone, non dalla tenuta delle registrazioni.

> [!ATTENZIONE] In ISO/IEC 17025 la registrazione tecnica deve consentire la ricostruzione dell'attivita'. Se quel collegamento deve essere "raccontato" a posteriori, la catena delle evidenze è debole.

### "Non l'abbiamo scritto perché è ovvio"

In un laboratorio accreditato, ciò che è critico per il risultato non può essere affidato all'ovvio. Deve essere governato.

---

## 3. Errori nella gestione delle NC che amplificano la criticità

### Chiudere la NC senza capire la causa reale

"Abbiamo corretto il modulo", "abbiamo rifatto la prova", "abbiamo informato il personale", "abbiamo aggiornato la procedura". Queste possono essere azioni utili, ma non sono automaticamente azioni correttive efficaci. Se il laboratorio non dimostra di aver capito perché il problema è accaduto, dove poteva essere intercettato e se può ripresentarsi, la criticità non è governata — è solo archiviata.

### Causa scritta come "errore umano"

"Errore umano", "disattenzione", "mancata applicazione della procedura", "carenza di formazione". Queste formule possono descrivere l'evento, non la causa vera. La domanda giusta non è solo "perché il tecnico si è dimenticato": è perché il sistema glielo ha consentito, perché il riesame non lo ha intercettato, perché la registrazione non guida il dato critico, perché quella informazione non viene controllata prima dell'emissione.

Quando la causa resta superficiale, l'azione correttiva diventa debole. Si fa una riunione, si firma un addestramento, ma il meccanismo che ha generato l'errore rimane intatto.

### Non fare la valutazione di estensione

Trovare un problema su una prova e trattarlo come episodio isolato senza verificare se la stessa condizione riguardi altri metodi, matrici, strumenti, periodi, clienti o rapporti già emessi. Se una formula Excel era sbagliata — da quando era in uso? Su quali prove? Con quale impatto numerico? Se una taratura era fuori criterio — quali risultati coprono quel periodo?

### Verifica di efficacia puramente formale

"La procedura è stata aggiornata", "il personale è stato formato" dimostrano che l'azione è stata fatta, non che ha funzionato. L'efficacia si dimostra con evidenze successive e pertinenti: nuovi rapporti campionati senza ricorrenza, sedute gestite correttamente, registrazioni complete su casi reali, trend migliorati.

Una verifica di efficacia fatta il giorno dopo la formazione, senza campionare attività reali, è quasi sempre debole.

### Ricorrenza mascherata

Prima mancavano i lotti degli standard, poi le identificazioni delle dotazioni, poi le condizioni ambientali, poi i raw data non erano collegati. Formalmente NC diverse; sostanzialmente la stessa debolezza nelle registrazioni tecniche. Il laboratorio corregge episodi ma non presidia processi.

> [!ISPETTORE] Una NC non si chiude quando il modulo è completo. Si chiude quando il laboratorio può dimostrare che il processo è tornato sotto controllo e che i risultati coinvolti sono stati protetti.

---

## 4. Errori di data integrity — Subito gravissimi anche se presentati come sviste

### Perdita o indisponibilità del dato grezzo primario

Se il laboratorio emette un rapporto ma non riesce a mostrare il dato originario da cui deriva — file strumentale, pesata, lettura, cromatogramma, spettro, sequenza, quaderno, log, output originario — l'ispettore non può verificare la catena tecnica. Nessuna dichiarazione verbale compensa.

### Non attribuibilità delle operazioni

Utenze generiche sugli strumenti, password condivise, firme elettroniche non personali, operatori che lavorano con l'account del collega, file salvati da postazioni comuni senza audit trail. Non è più dimostrabile chi ha fatto cosa e quando. Il riesame tecnico perde valore.

### Modifiche a dati o calcoli senza traccia

Foglio di calcolo validato sovrascritto, formula cambiata senza blocco o verifica, file strumentale rielaborato cancellando la versione precedente, risultato nel LIMS modificato senza motivazione. Il laboratorio non sta solo rischiando un errore: sta mostrando che le barriere di integrità del dato non sono adeguate.

### Correzioni cartacee non tracciate

Bianchetto, cancellature, riscritture, fogli sostituiti, registrazioni rifatte "in bella", date aggiunte dopo, firme apposte in blocco. Anche quando il valore finale è corretto, la modalità di correzione compromette la fiducia nel processo: non consente di distinguere errore originario, correzione, motivazione e responsabilità.

### Ricostruzioni a posteriori di dati non registrati in tempo reale

Se manca un lotto, una pesata, una temperatura, un operatore, una sequenza, e il laboratorio la ricostruisce giorni dopo sulla base di memoria o abitudine — sta trasformando una carenza di registrazione in un potenziale problema di integrità. L'integrazione tardiva è possibile solo se chiaramente identificata con data, autore, motivo e base oggettiva.

### LIMS e fogli Excel non verificati, audit trail non governati

Molti laboratori pensano che "se è nel LIMS allora è controllato". In realtà il LIMS è affidabile solo se configurazioni, profili utente, calcoli, autorizzazioni, audit trail, modifiche e backup sono governati. Audit trail disattivati, non riesaminati o non conservati sono un segnale serio.

> [!ISPETTORE] La domanda non è "e' successo una volta?". La domanda è "se fosse successo altre dieci volte, il sistema se ne sarebbe accorto?". Se la risposta è no — non è una svista. È una perdita di controllo sul dato.

---

## 5. I cinque errori più ricorrenti con conseguenze serie sull'accreditamento

Questi cinque errori hanno un elemento comune: non compromettono necessariamente ogni risultato prodotto, ma compromettono la capacità del laboratorio di dimostrare che quel risultato è valido, tracciabile e difendibile. Ed è lì che l'accreditamento diventa vulnerabile.

### ❶ Registrazioni tecniche non robuste

Fogli incompleti, dati grezzi non collegabili al campione, correzioni non tracciate, lotti mancanti, identificazione non chiara di strumenti/operatori/condizioni, ricostruzioni affidate alla memoria. Spesso percepita come carenza formale; in visita diventa seria appena impedisce di ricostruire con continuità come il risultato è stato prodotto.

### ❷ Controlli qualità non realmente interpretati

I CQ vengono registrati, ma non "letti". Sedute accettate con segnali borderline, trend evidenti, recuperi anomali trattati come eventi isolati. Il problema non è avere un controllo fuori criterio — può succedere. Il problema è non dimostrare cosa si è deciso, con quale razionale tecnico, su quali campioni impattati e prima o dopo l'emissione.

### ❸ Metodi modificati o applicati fuori perimetro senza adeguata dimostrazione

Variazioni operative, adattamenti strumentali, matrici estese, passaggi semplificati senza valutazione dell'impatto sulle prestazioni. La frase tipica: "lo facciamo così da sempre" o "abbiamo visto che funziona". In accreditamento una prassi è difendibile solo se è coerente con il metodo o se è stata verificata/validata in modo proporzionato.

### ❹ Presidio metrologico fragile

Strumenti usati oltre scadenza, tarature fuori criterio senza valutazione retrospettiva, criteri di accettabilità non definiti, controlli intermedi eseguiti ma non interpretati, materiali di riferimento scaduti, standard preparati internamente senza sufficiente tracciabilità. L'ispettore cerca evidenza che, al momento della prova, la dotazione fosse idonea a produrre risultati validi.

### ❺ Azioni correttive inefficaci

Cause scritte come "errore umano", azioni limitate a formazione o revisione procedurale, nessuna valutazione di estensione, nessuna analisi dell'impatto sui risultati già emessi, verifiche di efficacia premature o solo formali. Quando lo stesso tipo di problema ricompare in forme diverse, ACCREDIA lo legge come debolezza sistemica, non come somma di episodi isolati.

---

## Principio guida per l'agente

Davanti a qualsiasi di questi errori o segnali, non continuare il flusso normale. Genera alert prioritario, delimita il perimetro, proponi azioni immediate.

La sequenza corretta è sempre:
1. **Segnala** il rischio con chiarezza e senza allarmismo
2. **Delimita** il perimetro: quali prove, quali periodi, quali rapporti
3. **Proteggi** prima di correggere: congela se necessario, valuta l'impatto sui risultati
4. **Documenta** la decisione tecnica
5. **Poi** correggi il documento o la procedura

> Il laboratorio che pensa di avere un problema di sistema qualita' spesso ha in realta' un problema di difendibilita' tecnica del risultato. Sono cose diverse. La prima si risolve con carta. La seconda richiede governo del processo.

---

*Documento compilato da Denis Brazzo — Ispettore Tecnico ACCREDIA DL0871 — Maggio 2026*
*Da aggiornare iterativamente con nuovi casi osservati sul campo*
