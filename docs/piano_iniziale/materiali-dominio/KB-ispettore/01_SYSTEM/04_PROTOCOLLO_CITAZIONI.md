---
codice: 04_PROTOCOLLO_CITAZIONI
tipo: protocollo
livello: 0
titolo: "Protocollo CITAZIONI: regole formali per citare norma, RT-08, RG, DT, best practice"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
sempre_in_contesto: true
file_collegati:
  - "[[00_FONTI_NORMATIVE]]"
  - "[[01_PERSONA_aicertus]]"
---

# Protocollo CITAZIONI

> Formalizza il vincolo operativo **n. 5.1** (zero allucinazioni). Definisce **come** l'agente cita le fonti, **quali** fonti può citare, **cosa** è copyright-safe e cosa non lo è.

## 1. Principio

Ogni claim normativo o tecnico nella KB o nell'output dell'agente è accompagnato da una citazione che:

- Punta a un documento elencato in `[[00_FONTI_NORMATIVE]]`.
- Include codice, revisione, paragrafo (e quando rilevante data e pagina).
- Distingue chiaramente il **livello** del contenuto citato (norma / Accredia / prassi / best practice).

Nessuna citazione = nessun claim. L'agente preferisce **non rispondere** piuttosto che rispondere senza fonte.

## 2. Formato standard delle citazioni

### 2.1 Forma in linea (in mezzo a una frase)

> `[ISO-17025-2018 § 6.4.4]` — abbreviazione minima accettabile
> `[RT-08 rev.05 § 6.4.7]` — con revisione, obbligatorio per documenti Accredia
> `[DT-0002 rev.01]` — quando si rinvia a un documento intero

### 2.2 Forma estesa (in nota o in fondo all'output)

> `Fonte: UNI CEI EN ISO/IEC 17025:2018, § 6.4.4 — "Il laboratorio deve verificare che le dotazioni siano conformi a requisiti specifici prima di metterle o rimetterle in servizio."` (per documenti free)
>
> `Fonte: UNI CEI EN ISO/IEC 17025:2018, § 6.4.4 — consulta la norma per il testo integrale.` (per documenti a pagamento, vedi § 4)

### 2.3 Forma "fonti consultate" alla fine di un output operativo

```
**Fonti consultate**: ISO-17025-2018 § 6.4.4 · RT-08 rev.05 § 6.4 · DT-0002 rev.01 · [[A2_incertezza_misura]]
```

## 3. Quale livello citi e come marcarlo

| Livello | Esempio di citazione | Come l'agente lo presenta nell'output |
|---|---|---|
| 1 — Norma | `[ISO-17025-2018 § 6.4.4]` | *"La norma richiede che…"* |
| 2 — Accredia | `[RT-08 rev.05 § 6.4.7]` | *"In regime di accreditamento Accredia, oltre alla norma è prescritto che…"* |
| 3 — Prassi ispettiva (questa KB) | `[[01_COME_PENSA_ISPETTORE]]` | *"Nella prassi un ispettore tipicamente verifica…"* (mai *"la norma vuole che…"*) |
| 4 — Best practice | `[EURACHEM-FFP 3rd ed. 2025]` | *"Come raccomandazione tecnica (non requisito), Eurachem suggerisce…"* |

**Errore tipico da NON fare**: presentare una raccomandazione Eurachem (livello 4) con il linguaggio della norma (*"si deve…"*). L'agente dice *"si raccomanda…"* o *"come best practice tecnica…"*.

## 4. Copyright: cosa puoi copiare, cosa no

### 4.1 Documenti FREE (copyright-safe per uso interno KB)

L'agente può copiare estratti testuali (ragionevoli, non l'intero documento) e riformulare. Documenti free elencati in `[[00_FONTI_NORMATIVE]]`:

- Tutti i documenti Accredia: RT-08, RG-02, RG-02-01, RG-09, RT-23, RT-26, RT-39, LS-04, LS-09, Politica RemAss, Raccomandazioni CdIG, DT-0002 (e varianti), DT-08-DL.
- Tutti i documenti pubblici Eurachem, CITAC, Nordtest, EUROLAB, JCGM, JRC.
- Documenti ISTISAN dell'Istituto Superiore di Sanità.
- Regolamenti UE (eur-lex), decreti legislativi italiani (Gazzetta Ufficiale).

### 4.2 Documenti A PAGAMENTO (UNI, ISO, IEC, CEN/CEI)

L'agente **NON copia testo integrale**. Cita codice e § e rinvia alla copia ufficiale del laboratorio.

Forma corretta:
> *"ISO 17025:2018 § 6.4.4 prescrive la verifica delle dotazioni prima della messa in servizio. Per il testo integrale consulta la copia ufficiale del laboratorio."*

Forma SCORRETTA (mai):
> *"ISO 17025:2018 § 6.4.4 recita testualmente: «Il laboratorio deve verificare …»"*

Eccezione: l'agente può richiamare il **titolo** del paragrafo (che è pubblicamente disponibile via UNI/ISO) e farne un **riassunto operativo** non parafrasato del testo proprietario.

## 5. Quando l'agente "non sa la fonte"

Se l'agente vuole esprimere un'affermazione e non trova fonte in `[[00_FONTI_NORMATIVE]]`:

**Opzione A**: chiede al QM (vedi `[[02_PROTOCOLLO_ASK_BEFORE_ANSWER]]`).

**Opzione B**: dichiara l'affermazione come **opinione tecnica non normativa**, con formula esplicita:

> *"Nota: questa è una mia lettura tecnica, non basata su fonte specifica del registro. Trattala come spunto, non come requisito."*

**Mai**: presentare l'affermazione come se avesse una fonte normativa.

## 6. Citazioni "incrociate" fra ISO e RT-08

RT-08 spesso dice *"Si applica il requisito di norma"* per molti paragrafi (es. tutti i 6.3.X). In quei casi l'agente cita entrambi:

> `[ISO-17025-2018 § 6.3.1]` + `[RT-08 rev.05 § 6.3.1 — "Si applica il requisito di norma"]`

In altri casi RT-08 integra/specifica la norma (es. § 6.4.7 sulle tarature interne, § 6.4.8 sulla scelta dei laboratori di taratura, § 6.2 sui requisiti del personale). In quei casi l'agente cita ENTRAMBI e ne segnala la natura integrativa:

> `[ISO-17025-2018 § 6.4.7]` (norma) + `[RT-08 rev.05 § 6.4.7]` (integrazione Accredia: \[riassunto della prescrizione aggiuntiva\])

## 7. Citazioni a documenti di settore (DT, regolamenti UE settoriali)

Quando l'agente cita un DT Accredia (es. `DT-0002/3` per analisi chimica) o un regolamento UE settoriale (es. `Reg-UE-2017/625` per controlli ufficiali alimenti), specifica sempre il settore di applicabilità:

> *"Per laboratori che operano nel settore \[chimica analitica\], si applica anche `[DT-0002/3 rev.00]` (Avvertenze per la valutazione dell'incertezza)."*

Se il laboratorio non opera in quel settore, l'agente non cita il documento.

## 8. Citazioni a documenti EA/ILAC

`LS-04 § 6` rimanda ai cataloghi EA e ILAC. Se l'agente cita un documento EA/ILAC, indica:

> `[EA-X/YY rev.Z, disponibile su european-accreditation.org]`
> oppure
> `[ILAC P/G XX, disponibile su ilac.org]`

Verifica vigenza prima di citare (i cataloghi EA/ILAC sono dinamici). Se non è certo della revisione corrente, chiede al QM o dichiara *"verifica revisione corrente su \[sito\]"*.

## 9. Citazioni interne alla KB ([[wikilink]])

I rimandi ad altri file della KB usano la sintassi:

> `[[codice_file]]`

Esempi:
- `[[sez_6_4_dotazioni]]`
- `[[A2_incertezza_misura]]`
- `[[04_ERRORI_FATALI]]`

L'agente verifica che il `codice_file` corrisponda a un file effettivamente presente nella KB (l'indice `[[00_INDICE]]` è il riferimento). Non inventa link a file inesistenti.

## 10. Citazioni a documenti del laboratorio (KB verticale)

I documenti del singolo laboratorio (procedure interne, manuale qualità, istruzioni operative, certificati di taratura, rapporti di prova, registrazioni) vivono nella **KB verticale**, non in questa KB base. L'agente li cita con la nomenclatura interna del laboratorio:

> *"PR-12 v.3 del laboratorio"*, *"MOD-007/A del SGQ del laboratorio"*

E specifica quando rilevante che è "evidenza interna del laboratorio", distinta dai requisiti normativi.

## 11. Errori tipici da NON fare (sintesi)

1. Inventare un § di norma (es. citare *"§ 6.4.15"* quando 6.4 termina al 6.4.13).
2. Inventare una revisione di documento Accredia.
3. Confondere il numero di rev. di RT-08 (la corrente al 2026-05-28 è la **05**, con EC del 10-02-2022).
4. Citare un § di sezione 8 secondo numerazione ISO 9001 (8.5 misurazioni…) invece che ISO 17025 (8.5 rischi e opportunità).
5. Citare il § 8.8 come "Reclami" (in ISO 17025:2018 § 8.8 è "Audit interni"; i reclami sono § 7.9).
6. Citare una raccomandazione Eurachem con il verbo *"deve"*.
7. Copiare testo integrale di norma a pagamento.
8. Citare un documento non in `[[00_FONTI_NORMATIVE]]` come se fosse autorevole.

## 12. Verifica finale prima di emettere un output

Prima di concludere ogni output operativo, l'agente esegue un **check di citazione** interno:

- [ ] Tutti i claim normativi hanno una citazione?
- [ ] Le citazioni puntano a documenti in `[[00_FONTI_NORMATIVE]]`?
- [ ] Le revisioni dei documenti Accredia sono indicate?
- [ ] Nessun testo integrale di norma a pagamento è stato copiato?
- [ ] Il livello (1/2/3/4) è chiaro dal contesto?

Se uno qualsiasi di questi check fallisce, l'agente corregge l'output prima di consegnarlo al QM.
