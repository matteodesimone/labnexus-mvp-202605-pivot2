---
codice: sez_7_11_controllo_dati
tipo: requisito_normativo
livello: misto
titolo: "§ 7.11 — Controllo dei dati e gestione delle informazioni"
sezione_norma: "7.11"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.11" }
  - { doc: "RT-08", rev: "05", paragrafo: "7.11" }
fonti_secondarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "7.5 (registrazioni tecniche)" }
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "4.2 (riservatezza)" }
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[02_TRE_LIVELLI_DI_CONTROLLO]]"
  - "[[sez_7_5_registrazioni_tecniche]]"
  - "[[sez_7_8_presentazione_risultati]]"
tags: [LIMS, software_validato, integrita_dati, fogli_elettronici, backup, cloud, audit_trail, cybersecurity]
escalation_qm: true
---

# § 7.11 — Controllo dei dati e gestione delle informazioni

> Il § 7.11 disciplina il "sistema nervoso" del laboratorio: LIMS, software dedicati, fogli elettronici, database, cloud, backup, riservatezza. Tre vincoli base: (1) i sistemi devono essere **validati** prima dell'uso e dopo ogni modifica significativa, (2) devono **proteggere l'integrità** dei dati (audit trail, autorizzazioni, backup), (3) deve essere garantita **riservatezza** e accessibilità. Per Accredia il punto più caldo è il **foglio elettronico Excel**: spessissimo è usato per calcoli critici senza essere documentato, validato, protetto.

## A. Testo della norma (🟢 LIVELLO 1 — fatto normativo)

> **Fonte**: `[ISO-17025-2018 § 7.11.1 - 7.11.6]` — testo a pagamento, non riprodotto. Riassunto operativo:

- **7.11.1** Il laboratorio deve avere **accesso a tutti i dati e le informazioni** necessarie per eseguire le attività di laboratorio.
- **7.11.2** I **sistemi informativi per la gestione delle informazioni** (LIMS / software / hardware) usati per raccolta, elaborazione, registrazione, reporting, immagazzinamento o recupero dei dati devono essere **validati per funzionalità** dal laboratorio (inclusi corretto funzionamento delle interfacce con altri sistemi) prima dell'introduzione. Ogni modifica, incluse modifiche al software dell'apparecchiatura o off-the-shelf, deve essere **autorizzata, documentata e validata** prima dell'attuazione. **Nota 1**: software off-the-shelf usato in modalità d'uso prevista (es. word processor, fogli di calcolo, programmi statistici) può essere considerato validato dal laboratorio. **Nota 2**: configurazioni e modifiche del software del laboratorio come fogli di calcolo per calcoli devono essere validate (vedi 7.11.6).
- **7.11.3** I sistemi informativi devono essere:
  - **a)** **protetti da accesso non autorizzato**;
  - **b)** **protetti da manomissione e perdita**;
  - **c)** **gestiti in un ambiente conforme alle specifiche** del fornitore o del lab (o in casi di sistemi non computerizzati, fornire condizioni che salvaguardano accuratezza della registrazione e trascrizione manuale);
  - **d)** **mantenuti in modo da assicurare integrità** dei dati e delle informazioni;
  - **e)** **inclusivi di registrazione di guasti del sistema** e di appropriate azioni immediate e correttive.
- **7.11.4** Quando il sistema informativo per la gestione delle informazioni è **gestito e mantenuto off-site** o tramite **fornitore esterno**, il laboratorio deve assicurare che il fornitore o l'operatore del sistema soddisfi tutti i requisiti applicabili della norma.
- **7.11.5** Il laboratorio deve assicurare che istruzioni, manuali e dati di riferimento rilevanti per il sistema informativo siano **prontamente disponibili al personale**.
- **7.11.6** **Calcoli e trasferimenti dati** devono essere controllati in modo appropriato e sistematico.

## B. Prescrizioni Accredia (🟢 LIVELLO 2 — prescrizione Accredia)

> **Fonte**: `[RT-08 rev.05 § 7.11]`

### 7.11.2 — Validazione di sistemi forniti esternamente

> **Citazione integrale RT-08**:
> *"L'estensione e la profondità delle validazioni da eseguire su un sistema fornito dall'esterno è commisurata alle evidenze che il produttore è in grado di fornire circa la conformità del sistema."*

Cioè: se il fornitore (es. produttore LIMS) fornisce robusta documentazione di validazione del prodotto (IQ/OQ/PQ standard), il lab può ridurre il proprio sforzo di validazione. Se documentazione fornitore scarsa, il lab deve fare validazione più approfondita in proprio.

### 7.11.4 — Off-site / cloud / fornitore esterno

> **Citazione integrale RT-08**:
> *"Il Laboratorio deve comunicare al fornitore esterno o all'operatore del sistema, nel caso di sistemi gestiti dall'organizzazione di cui il Laboratorio fa parte, i requisiti applicabili, inclusi quelli relativi a confidenzialità, integrità e accessibilità delle informazioni mantenute off-site."*

Cioè per cloud / outsourcing IT / gruppi multinazionali con datacenter centralizzati: il laboratorio deve **comunicare formalmente** al gestore i requisiti su:
- Confidenzialità (collega § 4.2);
- Integrità (es. audit trail, protezione modifica, backup);
- Accessibilità (es. SLA, tempi di ripristino).

Tipicamente questo si traduce in **clausole contrattuali** specifiche o **SLA** documentati col fornitore.

### 7.11.6 — Fogli elettronici e altri programmi di calcolo

> **Citazione integrale RT-08** (CRITICA):
> *"Con riferimento all'utilizzo di fogli elettronici o di altri programmi di calcolo commerciali, le applicazioni sviluppate dal Laboratorio (formule, macro) devono essere documentate, validate e protette per impedirne l'involontaria alterazione."*

I 3 requisiti per ogni Excel / foglio di calcolo / script / macro usata in laboratorio:
1. **Documentato**: descritto, identificato univocamente, mantenuto sotto controllo di versione.
2. **Validato**: testato per casi noti, edge case, dati anomali.
3. **Protetto**: contro alterazione involontaria (es. celle bloccate, password sulla struttura, file in cartella di sola lettura per gli operatori).

Punti 7.11.1, 7.11.3, 7.11.5: "Si applica il requisito di norma" senza integrazioni Accredia esplicite.

## C. Documenti applicabili (LS-04 rev.20)

- `ISO-17025-2018` § 7.11.
- `RT-08 rev.05` § 7.11.
- Riferimenti esterni non in LS-04 ma utili: GAMP 5 (validazione software in ambiente regolamentato), ISO 27001 (sicurezza informazioni), EU GDPR (per dati personali, es. campioni umani).
- Collegamenti interni: § 7.5 (registrazioni — audit trail), § 4.2 (riservatezza), § 6.5 (riferibilità — per software che gestisce taratura).

## D. Domande ispettive (🟡 LIVELLO 3 — prassi ispettiva)

> Pattern: meccanismo → estensione → efficacia.

### 7.11.2 — Validazione del LIMS e dei software

#### Domanda 1 — Meccanismo
**"Qual è il LIMS / sistema informativo principale del laboratorio? Mostratemi la documentazione di validazione iniziale (IQ/OQ/PQ o equivalente) e la documentazione di ciò che il fornitore ha già fornito come "pre-validato". L'estensione della vostra validazione è commisurata a quanto fornito dal produttore (RT-08 § 7.11.2)?"**

Analisi ispettiva: documentazione di validazione del LIMS è la base. Tipici problemi:
- LIMS installato senza alcuna validazione documentata, perché "lo abbiamo comprato e funziona".
- IQ/OQ/PQ assenti o solo cartacei firmati dal fornitore senza verifica indipendente.
- Modifiche successive (es. nuova prova aggiunta, nuovo modulo) non validate.

#### Domanda 2 — Estensione
**"Per altre apparecchiature con software dedicato (es. GC-MS, ICP, autoanalizzatori): il software è stato validato? Le interfacce (es. trasferimento dati da GC a LIMS) sono validate (norma esplicita)?"**

Analisi ispettiva: § 7.11.2 esplicito su interfacce. Tipico problema: il LIMS è validato, lo strumento è validato, ma il **trasferimento dati** strumento → LIMS non è mai stato verificato. Risultato: dati che cambiano in trasferimento senza che nessuno se ne accorga.

#### Domanda 3 — Efficacia
**"Per una modifica recente al LIMS (es. nuovo metodo aggiunto, nuova formula di calcolo, integrazione con apparecchio nuovo): mostrate l'autorizzazione, la documentazione e la validazione (norma esplicita)."**

Analisi ispettiva: § 7.11.2 esplicito su "ogni modifica deve essere autorizzata, documentata e validata". Tipico problema: l'amministratore IT modifica configurazioni senza tracciabilità tecnica/qualità.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Validazione LIMS + software strumentali + interfacce, con documentazione fornitore + estensione validazione lab + validazione modifiche. ISO § 7.11.2 + RT-08 § 7.11.2."*
- Urgenza: alta.

### 7.11.3 — Protezione, integrità, gestione

#### Domanda 1 — Meccanismo
**"Per il LIMS: come gestite (a) accessi e autorizzazioni (chi vede cosa, chi modifica cosa), (b) backup (frequenza, retention, restore test), (c) audit trail (è attivato su tutti i campi tecnici? è disattivabile?), (d) registrazione guasti e azioni?"**

Analisi ispettiva: i 5 elementi del § 7.11.3 (a-e). Cerco:
- **Accessi**: principio del minimo privilegio. Operatore vede solo le proprie prove, RT/QM vede tutto, amministratore IT è separato dagli operatori tecnici.
- **Backup**: frequenza (quotidiana per dati operativi è il minimo), retention (allineata al tempo di conservazione delle registrazioni — collega § 7.5), **restore test periodico** (almeno annuale — molti lab fanno backup ma non hanno mai testato il restore).
- **Audit trail**: ATTIVATO su tutti i campi che riguardano dato tecnico. NON disattivabile (o disattivabile solo dall'amministratore con tracciabilità).
- **Registrazione guasti**: ticket / log incidenti con azione correttiva.

#### Domanda 2 — Estensione
**"Mostratemi il log di accessi al LIMS dell'ultima settimana. Qualcuno ha modificato dati tecnici al di fuori dell'orario operativo? Qualcuno ha effettuato modifiche bulk?"**

Analisi ispettiva: campagna anomalie. Anomalie non spiegate → indagine.

**"Mostratemi il report di un restore test recente."**

Analisi ispettiva: se "non l'abbiamo mai fatto" → NC sostanziale: backup non testato = backup non garantito.

#### Domanda 3 — Efficacia
**"Se domani il LIMS va in crash, in quanto tempo siete operativi? Avete BCP (business continuity)?"**

Analisi ispettiva: criterio di efficacia. SLA con fornitore IT? Procedura manuale fallback?

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Sistema informativo: accessi, backup con restore test, audit trail attivo non disattivabile, registrazione guasti. ISO § 7.11.3."*

### 7.11.4 — Off-site / cloud / outsourcing

#### Domanda 1 — Meccanismo
**"Il LIMS / il datacenter è on-premise o in cloud? Se cloud o gestito da fornitore esterno o da società del gruppo: come avete formalizzato i requisiti su confidenzialità, integrità e accessibilità (RT-08 § 7.11.4)?"**

Analisi ispettiva: RT-08 esplicito. Deve esistere comunicazione formale dei requisiti (clausole contrattuali, SLA, audit periodici al fornitore). Tipico problema: contratto IT firmato dall'ufficio acquisti senza coinvolgimento Qualità, senza requisiti specifici di lab.

#### Domanda 2 — Estensione
**"Mostratemi il contratto / SLA col fornitore IT. Dove sono dichiarati i requisiti su riservatezza (collega § 4.2), integrità, accessibilità, retention, restore time?"**

Analisi ispettiva: lettura contrattuale. Cerco clausole specifiche, non generiche.

#### Domanda 3 — Efficacia
**"Se il fornitore cloud avesse un data breach, come ne sareste avvisati e cosa fareste verso i vostri clienti?"**

Analisi ispettiva: catena di notifica. Importante anche per GDPR (se ci sono dati personali).

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Off-site / cloud: comunicazione formale requisiti, SLA, gestione incidenti. RT-08 § 7.11.4 + § 4.2."*

### 7.11.6 — Fogli elettronici Excel (il punto caldo)

#### Domanda 1 — Meccanismo
**"Quali fogli elettronici Excel (o altri programmi di calcolo) usate per calcoli che entrano in un risultato di prova accreditata? Avete un inventario? Per ciascuno, mostrate: (a) documentazione (cos'è, chi l'ha sviluppato, quando, per cosa), (b) validazione (test case noti, edge case), (c) protezione (celle bloccate, password sulla struttura, accesso in sola lettura per operatori)."**

Analisi ispettiva: RT-08 esplicito sui 3 requisiti. Tipici problemi cumulativi:
- Inventario inesistente — i fogli sono "sparsi" nelle cartelle dei tecnici.
- "Validazione" = "abbiamo provato e funziona", senza casi test documentati.
- Nessuna protezione: ogni operatore può modificare formule "per migliorare".
- Versioning inesistente: il file è "Calcolo_v3_finale_OK_definitivo.xlsx".

#### Domanda 2 — Estensione (test concreto)
**"Aprite un foglio Excel usato per un calcolo recente. Mostrate: la formula nelle celle critiche, l'evidenza di protezione (struttura bloccata), la versione, chi può modificarlo. Mostratemi il dossier di validazione: casi test, valori attesi, valori ottenuti, esito."**

Analisi ispettiva: prova pratica. Se la formula è in chiaro e modificabile da chiunque → NC RT-08.

#### Domanda 3 — Efficacia
**"Se domani modificate una formula per migliorarla: qual è il processo? Test → autorizzazione → validazione → release in produzione? Le versioni precedenti sono mantenute per ricostruzione tecnica (collega § 7.5)?"**

Analisi ispettiva: change management sui fogli. Se "modifichiamo e basta" → NC.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Fogli elettronici Excel: inventario + documentazione + validazione + protezione + versioning + change management. RT-08 § 7.11.6 (i 3 requisiti) + collegamento § 7.5 per ricostruzione tecnica."*
- Urgenza: alta — è uno dei findings più frequenti di Accredia.
- Cosa NON può chiudere l'agente: NON dichiara il § 7.11.6 conforme senza vedere inventario fogli + 1 dossier di validazione completo + verifica protezione.

### 7.11.5 — Documentazione tecnica al personale

#### Domanda 1 — Meccanismo
**"Manuali del LIMS, istruzioni operative dei moduli, manuali strumenti con software: dove sono e come sono accessibili al personale che li deve usare?"**

Analisi ispettiva: norma esplicita "prontamente disponibili". Tipico problema: manuali in stanza del responsabile o solo digitali su account amministratore.

#### Domanda 2 — Estensione
**"L'operatore X può accedere alla documentazione necessaria per la sua attività? Verifico."**

Analisi ispettiva: test pratico.

#### Domanda 3 — Efficacia
**"Quando un manuale viene aggiornato (es. nuova release LIMS), come comunicate l'aggiornamento al personale?"**

Analisi ispettiva: change management documentazione.

**Escalation al QM se la risposta non regge**
- Bozza per il QM: *"Accessibilità documentazione tecnica al personale operativo. ISO § 7.11.5."*

## E. Errori fatali correlati

- **EF-7.11.2-A**: LIMS o software dedicato in uso senza validazione documentata — NC esplicita ISO.
- **EF-7.11.2-B**: interfacce strumento → LIMS non validate, dati che cambiano in trasferimento — NC grave (impatta integrità).
- **EF-7.11.2-C**: modifica al LIMS / configurazione senza autorizzazione/documentazione/validazione — NC.
- **EF-7.11.3-A**: audit trail disattivato o disattivabile per utenti tecnici — NC grave su integrità dati.
- **EF-7.11.3-B**: backup esistente ma mai testato (no restore test) — backup illusorio, NC sostanziale.
- **EF-7.11.3-C**: accessi LIMS condivisi (es. login generico "labtech") che impedisce tracciabilità individuale — NC, collega § 7.5.
- **EF-7.11.4-A**: fornitore esterno / cloud / gruppo IT senza comunicazione formale dei requisiti del laboratorio (confidenzialità, integrità, accessibilità) — NC esplicita RT-08.
- **EF-7.11.6-A**: fogli Excel non documentati / non validati / non protetti — NC esplicita RT-08.
- **EF-7.11.6-B**: foglio Excel modificato dall'operatore senza change management — NC + impatto su § 7.5 (registrazioni precedenti potrebbero essere irriproducibili).
- **EF-7.11.6-C**: foglio Excel usato per calcolo critico in più versioni "sparse" senza versioning — NC.
- **EF-7.11-D**: assenza inventario sistemi/fogli — non c'è scope chiaro di cosa validare.

Vedi anche [[04_ERRORI_FATALI]].

## F. Quando l'agente DEVE chiedere prima di rispondere

L'agente **non risponde** e **chiede al QM** se:

1. Si chiede di introdurre un nuovo software / modulo LIMS / foglio Excel per calcolo critico: prima va validato (§ 7.11.2 + § 7.11.6).
2. Si chiede di modificare una formula in un foglio già in uso: change management con tracciabilità.
3. Audit trail del LIMS è disattivato per "questioni di performance": NC immediata, decisione QM + IT.
4. Backup non testato da > 12 mesi: priorità di sicurezza.
5. Migrazione del LIMS / aggiornamento maggiore: piano di validazione e migrazione dati.
6. Outsourcing IT / spostamento a cloud / cambio fornitore: revisione contratto + clausole RT-08 § 7.11.4.
7. Data breach / incidente di sicurezza: profilo legale (GDPR se dati personali) + riservatezza clienti § 4.2.
8. Operatore segnala difficoltà di accesso a manuali / documentazione: § 7.11.5.
9. Stoccaggio dati in archivi storici in formato proprietario di software dismesso: rischio illeggibilità futura (collega § 7.5).

L'agente NON dichiara il § 7.11 conforme senza vedere: (a) inventario sistemi + fogli Excel, (b) dossier validazione LIMS + interfacce, (c) policy accessi + audit trail attivo, (d) test restore backup, (e) contratto/SLA fornitore IT con clausole RT-08 § 7.11.4, (f) almeno 1 foglio Excel con i 3 requisiti RT-08 § 7.11.6 soddisfatti.
