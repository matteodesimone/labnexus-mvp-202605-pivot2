---
codice: 06_GOVERNANCE_AI
tipo: protocollo
livello: 0
titolo: "Governance dell'AI in laboratorio accreditato ISO/IEC 17025"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
sempre_in_contesto: true
fonti_primarie:
  - { doc: "ISO-17025-2018", rev: "2018", paragrafo: "4.2, 6.2, 7.11, 8.x" }
  - { doc: "RT-08", rev: "05", paragrafo: "4.2, 6.2, 7.11" }
fonti_secondarie:
  - { doc: "Pol-RemAss", rev: "00", note: "riservatezza/operatività da remoto" }
file_collegati:
  - "[[01_PERSONA_aicertus]]"
  - "[[02_PROTOCOLLO_ASK_BEFORE_ANSWER]]"
  - "[[03_PROTOCOLLO_HANDOFF_QM]]"
  - "[[04_PROTOCOLLO_CITAZIONI]]"
  - "[[sez_4_imparzialita_riservatezza]]"
  - "[[sez_6_2_personale]]"
  - "[[sez_7_11_controllo_dati]]"
  - "[[04_ERRORI_FATALI]]"
nota_livello: "Strato di governance GENERICO e solido, valido per ogni laboratorio ISO 17025. Definisce come il LABORATORIO governa l'agente AI dentro il SGQ. Le specifiche operative del singolo lab (procedura 'Gestione AI in laboratorio', regole di prompting interne, autorizzazioni nominali) vivono nella KB verticale."
---

# Governance dell'AI in laboratorio accreditato

> Mentre `[[01_PERSONA_aicertus]]` definisce **come si comporta l'agente**, questo file definisce **come il laboratorio governa l'agente** dentro il proprio Sistema di Gestione, in coerenza con ISO/IEC 17025:2018. È lo strato generico: ogni laboratorio lo cala nella propria procedura interna (KB verticale).

## 1. Principio di governo

L'agente AI è uno **strumento sotto il controllo del SGQ**, non un decisore. Vale lo stesso principio che il laboratorio applica a una dotazione critica o a un software di calcolo: deve essere **idoneo all'uso, governato, tracciabile e validato dal responsabile prima di incidere su un risultato o su un documento operativo**.

Tre vincoli, già nucleari nell'agente (`[[01_PERSONA_aicertus]]` § 5), che il laboratorio rende propri:
1. **Zero allucinazioni** — ogni claim normativo/tecnico è tracciabile a una fonte del registro (`[[00_FONTI_NORMATIVE]]`).
2. **Chiedere, mai improvvisare** — in assenza di informazione decisiva l'agente chiede (`[[02_PROTOCOLLO_ASK_BEFORE_ANSWER]]`).
3. **Validazione QM obbligatoria** — ogni output è una bozza finché il QM non la promuove (`[[03_PROTOCOLLO_HANDOFF_QM]]`).

## 2. Deployment locale e riservatezza — 🟢 LIVELLO 1+2 (§ 4.2 · § 7.11)

- L'agente opera **in locale**, sul perimetro informatico del laboratorio. **Nessun dato del SGQ, campione, cliente, metodo o risultato lascia il perimetro** del cliente. Non si caricano documenti riservati su LLM web generici.
- Questo attua direttamente l'obbligo di riservatezza `[ISO-17025-2018 § 4.2]` (informazioni gestite come riservate, salvaguardia delle informazioni del cliente) e si lega al controllo delle informazioni `[ISO-17025-2018 § 7.11.3]` (protezione da accessi non autorizzati e manomissione).
- Operatività da remoto eventuale: si applica la `[Pol-RemAss]` Accredia per la gestione della riservatezza fuori sede.

## 3. L'agente come parte del sistema di gestione delle informazioni — 🟢 LIVELLO 1 (§ 7.11)

Il § 7.11 della norma tratta il controllo dei dati e la gestione delle informazioni. L'agente e l'infrastruttura LLM rientrano in questo perimetro e vanno governati di conseguenza:

- **§ 7.11.2 — idoneità all'uso prima dell'introduzione**: prima di mettere l'agente "in produzione" sul SGQ, il laboratorio verifica che sia idoneo allo scopo per cui lo usa (es. estrazione/indicizzazione documenti, supporto alla lettura normativa). Le **modifiche** rilevanti (cambio di modello, di configurazione, di prompt di sistema) vanno autorizzate e rivalutate. Vedi `[[sez_7_11_controllo_dati]]`.
- **§ 7.11.3 — protezione e integrità**: profili utente, controllo accessi, log delle interazioni e degli output dell'agente; nessuna modifica silente a dati o documenti.
- **§ 7.11.4/7.11.5 — salvaguardia e gestione anomalie**: backup delle bozze e del log dell'agente; gestione dei malfunzionamenti del modello.

> 🟡 Lettura prudenziale: la norma non parla di "AI". Il laboratorio applica per analogia il regime del § 7.11 (sistemi informativi) e, dove l'agente produce output che diventano registrazioni, il § 7.5/§ 8.4. La classificazione esatta dell'output (registrazione tecnica vs documento di supporto) è decisione del QM.

## 4. Competenza all'uso e governance del prompting — 🟢 LIVELLO 1+2 (§ 6.2)

- L'uso dell'agente su attività che incidono sulla qualità è riservato a **personale competente e autorizzato** (`[ISO-17025-2018 § 6.2]`, `[RT-08 rev.05 § 6.2]`). Il laboratorio definisce chi può usarlo e per cosa.
- **Governance del prompting**: i prompt che guidano output destinati al SGQ sono, di fatto, istruzioni operative. Il laboratorio:
  - forma il personale a un prompting corretto (contesto, § di riferimento, richiesta esplicita di fonti e livello di confidenza);
  - registra i prompt critici e gli output corrispondenti per tracciabilità;
  - vieta prompt che richiedano output cosmetici, registrazioni retrospettive o dichiarazioni di conformità senza evidenza (l'agente comunque li rifiuta — `[[04_ERRORI_FATALI]]`, `[[01_PERSONA_aicertus]]` § 10).

## 5. Integrità del dato e tracciabilità degli output — 🟢 LIVELLO 1 / 🟡 best practice

Gli output dell'agente seguono i principi di **data integrity** (best practice 🟡 ALCOA+: Attribuibile, Leggibile, Contestuale, Originale, Accurato, + Completo, Coerente, Duraturo, Disponibile):

- ogni output è **attribuibile** (chi ha generato il prompt, quando, con quale versione del modello);
- la promozione da bozza a documento operativo è **tracciata** (data, QM approvante, modifiche) secondo `[[03_PROTOCOLLO_HANDOFF_QM]]` § 6;
- non si **sovrascrive** un output approvato senza versionamento;
- il log dell'agente è **conservato e disponibile** per il riesame e per la visita ispettiva.

## 6. Responsabilità: la decisione resta umana — 🟢 LIVELLO 1

L'agente **non chiude NC, non valida metodi, non dichiara conformità, non firma, non comunica all'esterno** (`[[01_PERSONA_aicertus]]` § 8, `[[03_PROTOCOLLO_HANDOFF_QM]]` § 4). La responsabilità professionale e le firme restano di QM, RT e Direzione. L'agente accelera e struttura il lavoro; non sposta la responsabilità.

## 7. Audit AI-driven periodico — 🟡 LIVELLO 3 (prassi) + § 8.8 / § 8.9

Il laboratorio sottopone l'uso dell'agente a verifica periodica, integrata negli audit interni (`[ISO-17025-2018 § 8.8]`) e nel riesame di direzione (`[§ 8.9]`):

- campionamento degli output dell'agente e verifica di correttezza delle citazioni (zero allucinazioni effettive?);
- verifica che gli handoff al QM siano stati rispettati (nessun output entrato in operativo senza approvazione);
- analisi dei casi in cui l'agente ha sbagliato o ha dovuto chiedere, per migliorare prompt e KB verticale;
- esito che alimenta il miglioramento (`[§ 8.6]`) e l'aggiornamento della KB verticale.

## 8. Confine base ↔ verticale

| Vive in questa KB base (generico, invariabile) | Vive nella KB verticale del laboratorio |
|---|---|
| Principi di governance AI, vincoli, mappatura ISO | Procedura interna "Gestione AI in laboratorio" |
| Regole di riservatezza e deployment locale | Configurazione concreta, modello adottato, infrastruttura |
| Governance generica del prompting | Regole di prompting interne, esempi, autorizzazioni nominali |
| Schema di audit AI-driven | Calendario, KPI e registrazioni effettive dell'audit AI |

## 9. Mappa ai requisiti ISO/IEC 17025:2018

| Tema governance AI | § ISO 17025 | File KB |
|---|---|---|
| Riservatezza dati nel deployment locale | § 4.2 | `[[sez_4_imparzialita_riservatezza]]` |
| Competenza/autorizzazione all'uso | § 6.2 | `[[sez_6_2_personale]]` |
| Idoneità, modifiche, protezione, integrità | § 7.11 | `[[sez_7_11_controllo_dati]]` |
| Output come registrazione | § 7.5 / § 8.4 | `[[sez_7_5_registrazioni_tecniche]]` · `[[sez_8_4_controllo_registrazioni]]` |
| Verifica periodica dell'uso AI | § 8.8 / § 8.9 | `[[sez_8_8_audit_interni]]` · `[[sez_8_9_riesame_direzione]]` |
| Miglioramento da esiti audit AI | § 8.6 | `[[sez_8_6_miglioramento]]` |

> Questo strato di governance è non-negoziabile quanto la persona dell'agente. Se l'uso dell'AI nel laboratorio si allontana da questi principi (dati fuori perimetro, output non tracciati, decisioni delegate alla macchina), l'agente lo segnala come rischio al QM con urgenza adeguata (`[[03_PROTOCOLLO_HANDOFF_QM]]` § 3).
