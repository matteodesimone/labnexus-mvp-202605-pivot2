---
codice: 00_FONTI_NORMATIVE
tipo: registro_fonti
livello: 0
titolo: "Registro delle fonti normative citabili"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
sempre_in_contesto: true
base_estrazione: "LS-04 rev.20 del 30-09-2025, sezione 3 (Laboratori di prova) — entrata in vigore 01-10-2025"
nota: "Registro CHIUSO. L'agente cita solo documenti elencati qui. Se un documento non è qui, l'agente dichiara 'fonte non disponibile, chiedo al QM' (vedi protocollo ASK-BEFORE-ANSWER)."
---

# Registro delle fonti normative citabili

## 0. Come si usa questo registro

- Ogni claim normativo o tecnico nella KB **deve** essere accompagnato dal codice del documento qui elencato.
- Il campo **Accesso** dice se l'agente può estrarre/citare testo integrale oppure solo richiamare il riferimento (per le norme UNI/ISO a pagamento l'agente NON copia testo, rinvia alla copia ufficiale del laboratorio).
- Il campo **Livello** indica la gerarchia di citabilità: 1 norma · 2 Accredia · 3 prassi · 4 best practice.
- Il campo **Scope KB** dice in quali file della KB il documento è richiamato (back-link).

## 1. Norma base (livello 1)

| Codice | Documento | Rev / Data | Accesso | Scope KB |
|---|---|---|---|---|
| `ISO-17025-2018` | UNI CEI EN ISO/IEC 17025:2018 — Requisiti generali per la competenza dei laboratori di prova e di taratura | 2018 + EC 1-2018 del 19/06/2018 | A pagamento (UNI) | tutti i file `sez_*` |

## 2. Documenti Accredia per laboratori di prova (livello 2)

| Codice | Documento | Rev | Data | Accesso | Scope KB |
|---|---|---|---|---|---|
| `RT-08` | Prescrizioni per l'accreditamento dei laboratori di prova (rif. UNI CEI EN ISO/IEC 17025:2018) | 05 | 15-12-2021 + EC 10-02-2022; in vigore 01-04-2022 | Free (accredia.it) | tutti i file `sez_*` |
| `RG-02` | Regolamento per l'accreditamento dei laboratori di prova e dei laboratori medici | 08 | — | Free (accredia.it) | `sez_4_*`, `sez_8_*`, `01_SYSTEM` |
| `RG-02-01` | Regolamento per l'accreditamento dei laboratori multisito | 04 | — | Free (accredia.it) | `sez_5_*` |
| `RG-09` | Regolamento per l'utilizzo del logo e del marchio ACCREDIA | 12 + EC 28-01-2025 | — | Free (accredia.it) | `sez_7_8_*` |
| `RT-23` | Prescrizioni per la definizione del campo di accreditamento per i laboratori di prova | 05 | 15-04-2026; **in vigore 01-11-2026** (rev. generale; rev.04 valida fino al 31-10-2026) | Free (accredia.it) | `sez_5_*`, `sez_7_2_*`, `00_GLOSSARIO`, `04_ERRORI_FATALI`, `A8_*` |
| `RT-26` | Prescrizioni per l'accreditamento con campo di accreditamento flessibile | 07 | — | Free (accredia.it) | `sez_7_2_*`, `A1_*` |
| `RT-39` | Prescrizioni per la partecipazione a PT/ILC | 00 | — | Free (accredia.it) | `sez_7_7_*`, `A5_*` |
| `LS-04` | Elenco norme e documenti di riferimento per l'accreditamento (lab prova, medici, PT provider) | 20 | 30-09-2025; in vigore 01-10-2025 | Free (accredia.it) | base di questo registro |
| `LS-09` | Elenco documenti per le attività di taratura applicabili ai laboratori di prova (cfr. LS-04 § 2) | — | — | Free (accredia.it + euramet.org) | `sez_6_4_*`, `sez_6_5_*`, `A3_*` |
| `Pol-RemAss` | Politica ACCREDIA per l'esecuzione delle verifiche da remoto | 00 | — | Free (accredia.it) | `01_SYSTEM`, `02_LOGICA_ISPETTIVA` |
| `Racc-CdIG` | Raccomandazioni del Comitato di Indirizzo e Garanzia ACCREDIA — criteri omogenei verifica 17025 | 00 | — | Free (accredia.it) | `02_LOGICA_ISPETTIVA`, vari `sez_*` |
| `DT-0002` | Guida per la valutazione e l'espressione dell'incertezza nelle misurazioni | 01 | — | Free (accredia.it) | `sez_7_6_*`, `A2_*` |
| `DT-0002/1` | Esempi applicativi di valutazione dell'incertezza nelle misurazioni elettriche | 01 | — | Free | `A2_*` |
| `DT-0002/2` | Esempi applicativi di valutazione dell'incertezza nelle misurazioni meccaniche | 00 | — | Free | `A2_*` |
| `DT-0002/3` | Avvertenze per la valutazione dell'incertezza nel campo dell'analisi chimica | 00 | — | Free | `A2_*` |
| `DT-0002/4` | Esempi applicativi di valutazione dell'incertezza nelle misurazioni chimiche | 00 | — | Free | `A2_*` |
| `DT-0002/5` | Esempio applicativo per misurazioni su materiali strutturali | 01 | — | Free | `A2_*` |
| `DT-0002/6` | Guida al calcolo della ripetibilità di un metodo di prova e alla sua verifica nel tempo | 00 | — | Free | `A1_*`, `A2_*` |
| `DT-08-DL` | Guida per la taratura di strumenti nel settore della compatibilità elettromagnetica | 00 | — | Free | `A3_*` |

## 3. Vocabolari e norme orizzontali (livello 1)

| Codice | Documento | Rev / Data | Accesso | Scope KB |
|---|---|---|---|---|
| `ISO-17000` | UNI CEI EN ISO/IEC 17000:2020 — Valutazione della conformità, vocabolario | 2020 | A pagamento | `00_GLOSSARIO`, vari `sez_*` |
| `ISO-17011` | UNI CEI EN ISO/IEC 17011:2018 — Requisiti per gli organismi di accreditamento (richiamata da RT-23 rev.05 § 1, § 7.8.3 per la struttura del campo di accreditamento) | 2018 | A pagamento | `00_GLOSSARIO`, `sez_5_*`, `sez_7_2_*` |
| `UNI-EN-45020` | Normazione e attività connesse, vocabolario generale | 2007 | A pagamento | `00_GLOSSARIO` |
| `UNI-CEI-70099` / `VIM` | Vocabolario internazionale di metrologia (VIM) = JCGM 200:2012 | 2008 / 2012 | A pagamento (UNI) / Free (BIPM) | `00_GLOSSARIO`, `sez_6_5_*`, `A2_*` |
| `JCGM-100` / `GUM` | Guida all'espressione dell'incertezza di misura | 2008 | Free (BIPM) | `sez_7_6_*`, `A2_*` |
| `UNI-CEI-70098-3` | Incertezza di misura, parte 3 — guida espressione | 2016 + EC 03-05-2018 | A pagamento | `A2_*` |
| `ISO-9000` | UNI EN ISO 9000:2015 — SGQ, fondamenti e vocabolario | 2015 | A pagamento | `00_GLOSSARIO` |
| `ISO-9001` | UNI EN ISO 9001:2015 — SGQ, requisiti | 2015 | A pagamento | `sez_8_*` (analogie/differenze 17025 ↔ 9001) |
| `ISO-10012` | UNI EN ISO 10012:2004 / ISO 10012:2003 — Sistemi di gestione della misurazione | 2003-2004 | A pagamento | `sez_6_4_*` |

## 4. Norme per PT e CRM (livello 1, citate come riferimento esterno)

| Codice | Documento | Rev / Data | Accesso | Scope KB |
|---|---|---|---|---|
| `ISO-17043-2010` | UNI CEI EN ISO/IEC 17043:2010 — Requisiti per provider PT | 2010 | A pagamento | `A5_*` |
| `ISO-17043-2024` | UNI CEI EN ISO/IEC 17043:2024 (ISO/IEC 17043:2023) | 2024 | A pagamento | `A5_*` |
| `ISO-17034` | UNI CEI EN ISO 17034:2017 — Requisiti per produttori di materiali di riferimento | 2017 | A pagamento | `A4_*` (riferimento esterno, fuori scope KB base) |
| `ISO-33401` | Reference materials — Contents of certificates, labels and accompanying documentation | 2024 | A pagamento | `A4_*` |
| `ISO-33403` | Reference materials — Requirements and recommendations for use | 2024 | A pagamento | `A4_*` |
| `ISO-33405` | Reference materials — Approaches for characterization and assessment of homogeneity and stability | 2024 | A pagamento | `A4_*` |

## 5. Best practice tecniche (livello 4)

| Codice | Documento | Rev / Data | Accesso | Scope KB |
|---|---|---|---|---|
| `EURACHEM-FFP` | The Fitness for Purpose of Analytical Methods — Method Validation | 3rd ed. 2025 | Free (eurachem.org) | `A1_*` |
| `EURACHEM-PT` | Selection, Use and Interpretation of PT Schemes | 3rd ed. 2021 | Free | `A5_*` |
| `EURACHEM-VIM3-Intro` | Terminology in Analytical Measurement — Introduction to VIM 3 | 2nd ed. 2023 | Free | `00_GLOSSARIO`, `A2_*` |
| `EURACHEM-MicroAcc` | Accreditation for Microbiological Laboratories | 2023 | Free | `A1_*`, `A5_*` |
| `EURACHEM-Sampling` | Measurement uncertainty arising from sampling | 2nd ed. 2019 | Free | `A7_*` |
| `EURACHEM-CITAC-QUAM` | Quantifying Uncertainty in Analytical Measurement (CG4) | 3rd ed. 2012 | Free | `A2_*` |
| `EURACHEM-CITAC-Quality` | Guide to Quality in Analytical Chemistry | 3rd ed. 2016 | Free | `A1_*` |
| `Nordtest-TR-569` | Internal Quality Control — Handbook for Chemical Laboratories (Trollbook) | ed. 5.1 | Free (nordtest.info) | `sez_7_7_*` |
| `Nordtest-TR-537` | Handbook for calculation of measurement uncertainty in environmental laboratories | ed. 4, 2017 | Free | `A2_*`, `A7_*` |
| `EUROLAB-1/2007` | Measurement uncertainty revisited — alternative approaches | 2007 | Free | `A2_*` |
| `ISTISAN-22/39` | Eurachem/CITAC — Incertezza di misura dovuta al campionamento (trad. it. 17-05-2023) | 2023 | Free (iss.it) | `A7_*` |
| `ISTISAN-12/29` | Controllo qualità interno — manuale per i laboratori di analisi chimiche (trad. Nordtest 569) | 4ª ed. 2011 | Free (iss.it) | `sez_7_7_*` |
| `ISTISAN-03/30` | Quantificazione dell'incertezza nelle misure analitiche (trad. QUAM CG4) | 2ª ed. 2000 | Free (iss.it) | `A2_*` |
| `JCGM-106` | The role of measurement uncertainty in conformity assessment | 2012 | Free (BIPM) | `sez_7_8_*`, `A2_*` |
| `SANTE-11312` | Analytical QC and method validation procedures for pesticide residues analysis in food and feed | 2021 | Free (EU) | `A1_*` (settore agroalimentare) |
| `JRC-GMO-Flex` | European technical guidance document for flexible scope accreditation of GMO labs | 2nd ed. 2014 | Free (JRC) | `A1_*` (settore GMO) |
| `ISO-21748` | Repeatability/reproducibility/trueness estimates in measurement uncertainty estimation | 2017 | A pagamento | `A1_*`, `A2_*` |
| `ISO-19036` | Microbiologia alimenti — incertezza di misura per determinazioni quantitative | 2020 | A pagamento | `A2_*` (microbio) |
| `ISO-7218` | Microbiologia alimenti — requisiti generali e guida | 2024 | A pagamento | `A1_*` (microbio) |
| `IEC-Guide-115` | Uncertainty of measurement in conformity assessment, electrotechnical sector | 2023 | A pagamento | `A2_*` (elettrico) |
| `ISO-IEC-TS-23532-1` | Requirements for the competence of IT security testing labs — Part 1 (ISO/IEC 15408) | 2021 / UNI 2025 | A pagamento | `A1_*` (cybersecurity) |

## 6. Regolamenti europei e nazionali (livello 2 — obbligatori dove applicabili)

| Codice | Documento | Accesso | Scope KB |
|---|---|---|---|
| `Reg-UE-765-2008` | Regolamento (CE) n. 765/2008 — accreditamento e vigilanza del mercato | Free (eur-lex) | `01_SYSTEM`, `02_LOGICA_ISPETTIVA` |
| `Reg-UE-2019/1020` | Regolamento (UE) 2019/1020 — vigilanza mercato e conformità prodotti | Free | `01_SYSTEM` |
| `Reg-UE-2017/625` | Regolamento (UE) 2017/625 — controlli ufficiali alimenti/mangimi | Free | settore agroalimentare |
| `Reg-UE-2021/808` | Regolamento (UE) 2021/808 — metodi analitici residui farmacologici | Free | settore agroalimentare |
| `Reg-UE-305/2011` | Regolamento (UE) n. 305/2011 — prodotti da costruzione | Free | settore costruzioni |
| `Reg-UE-2024/3110` | Regolamento (UE) 2024/3110 — nuove norme armonizzate prodotti da costruzione | Free | settore costruzioni |
| `D-Lgs-106-2017` | Decreto Legislativo 16 giugno 2017 n. 106 — adeguamento nazionale Reg. 305/2011 | Free (gazzetta) | settore costruzioni |

## 7. Documenti EA/ILAC (livello 2, citati come applicabili dove richiamati)

`LS-04 § 6` rinvia ai cataloghi documentali ufficiali EA e ILAC, costantemente aggiornati. L'agente non duplica qui l'elenco completo: cita i documenti EA/ILAC solo quando un file di requisito o appendice esplicitamente li richiama. Risorse:
- European Accreditation: https://european-accreditation.org
- ILAC: https://ilac.org

**Documenti EA/ILAC esplicitamente richiamati nella KB:**

| Codice | Documento | Rev / Data | Accesso | Scope KB |
|---|---|---|---|---|
| `ILAC-G17` | Guidelines for Measurement Uncertainty in Testing | 01/2021 | Free (ilac.org) | `sez_7_2_*`, `sez_7_6_*`, `A2_*` |
| `ILAC-G18` | Guideline for the Formulation of Scopes of Accreditation | corrente | Free (ilac.org) | `sez_5_*`, `sez_7_2_*` (formulazione/valutazione del campo di accreditamento; richiamata da RT-23 rev.05 § 5) |

> **Nota vigenza EA/ILAC**: la storica `EA-4/16` ("Expression of Uncertainty in Quantitative testing") è stata **ritirata da EA (marzo 2021)** e NON va citata come fonte vigente. Per l'incertezza nelle prove il riferimento corrente è `ILAC-G17:01/2021`. Per i documenti ILAC, ai sensi di RT-23 rev.05 (nota 3), si intende equivalente anche la successiva versione "Global", salvo diverse indicazioni. La lista ufficiale aggiornata delle pubblicazioni EA è `EA-INF/01` (european-accreditation.org); le guide ILAC sono su ilac.org. Verificare sempre la revisione corrente prima di citare (i cataloghi EA/ILAC sono dinamici).

## 8. Documenti riconoscimenti speciali (WADA, FCC, EPA, ISED)

`LS-04 § 7` elenca i documenti applicabili per i riconoscimenti speciali. **Fuori scope della KB base** (specifici per laboratori che lavorano in questi schemi). Vengono integrati nella KB verticale del laboratorio in onboarding solo se applicabili.

## 9. Regole operative dell'agente sul registro

1. **Citare solo da qui**. Se un documento serve all'agente per rispondere e non è nel registro, l'agente dichiara la mancanza e chiede al QM (vedi `02_PROTOCOLLO_ASK_BEFORE_ANSWER.md`).
2. **Rispettare il copyright**. Per documenti contrassegnati "A pagamento", l'agente NON copia testo integrale. Cita solo codice, § e rinvia alla copia ufficiale del laboratorio.
3. **Verificare la vigenza**. Il registro fotografa lo stato a `2026-05-29`. Se l'agente sa o sospetta che una revisione più recente sia stata emessa da Accredia, lo dichiara e chiede al QM di verificare su https://www.accredia.it. **Promemoria attivo**: RT-23 rev.05 entra in vigore il **01-11-2026**; fino al 31-10-2026 resta vigente la rev.04. Dopo tale data aggiornare i riferimenti residui a rev.04.
4. **Citare la revisione esplicita**. Es. "RT-08 rev.05 § 6.4.7" — mai "RT-08 § 6.4.7" senza la rev.
5. **Per LS-04 stesso**: la revisione corrente nel registro è la rev.20 del 30-09-2025 in vigore dal 01-10-2025. Se è stata emessa una rev.21, il QM lo segnala e si aggiorna il registro.

> **Tu QM/RT che leggi**: questo file è la "biblioteca di sistema" della KB. Tienilo allineato. Quando Accredia emette nuove revisioni, aggiorni qui le date e nelle fonti dei file requisito che lo richiamano, oppure chiedi all'agente di proporti il diff.
