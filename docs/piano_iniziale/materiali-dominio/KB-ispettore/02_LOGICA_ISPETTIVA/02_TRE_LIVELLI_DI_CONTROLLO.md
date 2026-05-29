---
codice: 02_TRE_LIVELLI_DI_CONTROLLO
tipo: logica_ispettiva
livello: 3
titolo: "Tre livelli di controllo: meccanismo → estensione → efficacia"
versione_kb: "2.1"
data_ultima_revisione: "2026-05-29"
autore_originale: "Denis Brazzo — Ispettore Tecnico Accredia DL0871"
file_collegati:
  - "[[01_COME_PENSA_ISPETTORE]]"
  - "[[03_RISPOSTA_NC]]"
  - "[[04_ERRORI_FATALI]]"
  - "[[00_MAPPA_ESPLOSA]]"
nota_livello: "LIVELLO 3 — prassi ispettiva. Formalizza il pattern di domande ispettive applicato a ogni requisito ISO/RT-08 nei file `sez_*` della cartella `03_REQUISITI/`."
---

> **Marcatore di livello: 🟡 PRASSI ISPETTIVA (livello 3)** — Questo file definisce il **pattern interrogativo** che la KB usa in tutti i file di requisito. È il modo in cui un ispettore Accredia verifica se un controllo del laboratorio è "vivo" o solo dichiarato.

# Tre livelli di controllo

## 1. Da dove nasce il pattern

Quando un ispettore Accredia entra in un laboratorio non chiede tanto *"questo requisito è soddisfatto?"*. Chiede tre cose distinte, in sequenza, su qualunque presidio del SGQ:

1. **Avete il meccanismo?** — esiste una procedura, una barriera, un controllo documentato che governa l'aspetto?
2. **L'avete cercato altrove?** — quando un problema appare, è stato indagato come fenomeno isolato o come segnale potenzialmente sistemico?
3. **Funziona davvero?** — c'è un caso concreto, recente, su casistica reale, in cui il controllo ha intercettato qualcosa o impedito un'uscita non conforme?

Un laboratorio che risponde bene a tutti e tre i livelli ha un sistema gestionale "abitato" (vedi `[[01_COME_PENSA_ISPETTORE]]` § 2). Un laboratorio che risponde bene solo al primo (ha il meccanismo) e male agli altri due ha un SGQ formale che funziona poco. **Se manca anche solo uno dei tre livelli, la risposta del laboratorio è incompleta.**

## 2. Pattern applicato

Tutti i file `sez_*` della cartella `03_REQUISITI/` applicano sistematicamente, per ogni requisito ISO/RT-08, tre domande ispettive nella forma:

### Domanda 1 — Il meccanismo di controllo

Verifica se il laboratorio HA un controllo documentato che governa l'aspetto del requisito.

Esempio (su § 6.4.4, verifica di idoneità all'uso):
> *"Quale procedura governa la verifica di idoneità degli strumenti prima della messa o rimessa in servizio? Chi la esegue? Come è tracciata? Con quali criteri di accettazione?"*

**Cosa cerca un ispettore**: procedura documentata, autorità autorizzante, criterio scritto, modulo di registrazione.

**Cosa indica una risposta debole**: *"sì, sì, si controlla sempre"*, senza nome di procedura, senza modulo, senza criterio numerico.

### Domanda 2 — La ricerca dell'estensione

Verifica se il laboratorio ha CERCATO il problema/applicazione del controllo in tutto il proprio perimetro, non solo nei casi più ovvi.

Esempio (su § 6.4.4):
> *"Negli ultimi 12 mesi, su quanti e quali strumenti la verifica di idoneità è stata fatta? Avete una mappa di applicabilità rispetto al parco dotazioni accreditato? Esistono strumenti che ne sarebbero soggetti ma per cui non risulta registrazione?"*

**Cosa cerca un ispettore**: ricerca sistematica, mappa di copertura, scostamento tra atteso e realizzato.

**Cosa indica una risposta debole**: *"sì, lo facciamo per X e Y"*, senza dimostrare che X e Y rappresentano tutto l'insieme di strumenti soggetti.

### Domanda 3 — L'evidenza di efficacia

Verifica se il controllo HA FUNZIONATO davvero su casistica reale, non solo in teoria.

Esempio (su § 6.4.4):
> *"Portate un caso concreto degli ultimi 24 mesi in cui la verifica di idoneità ha intercettato uno strumento non idoneo all'uso. Cosa avete fatto? Quale impatto ha avuto sui risultati prodotti dal laboratorio? Come avete protetto i risultati passati o in corso?"*

**Cosa cerca un ispettore**: episodio reale, registrato, con catena d'evento (intercettazione → azione → impatto sui risultati → verifica di efficacia).

**Cosa indica una risposta debole**: *"finora non è mai successo"*. Un controllo che in due anni non ha mai intercettato nulla, su un parco strumenti significativo, è probabilmente un controllo cieco o non operativo.

## 3. Escalation al QM (modulo standard di chiusura della Domanda 3)

Ogni Domanda 3 nei file `sez_*` è seguita da un blocco di escalation al QM. È il momento in cui l'agente, valutando le risposte ricevute o le evidenze rilevate nel SGQ, propone al QM:

```markdown
**Escalation al QM se la risposta non regge**

- **Bozza per il QM**: [proposta concreta — apertura NC interna, azione preventiva, piano di estensione del controllo, audit di approfondimento, ...]
- **Riferimento normativo**: [ISO § X.Y] + [RT-08 § X.Y] + [eventuale DT/RG/Circolare]
- **Livello di urgenza suggerito**: [bassa | media | alta | critica], motivata
- **Cosa l'agente NON può chiudere da solo**: [decisione che resta al QM — es. chiusura NC, dichiarazione efficacia, comunicazione cliente]
- **Documenti che il QM dovrà acquisire prima di decidere**: [es. registrazioni mancanti, evidenze di taratura, certificati]
```

## 4. Come l'agente "valuta" le risposte ricevute

Quando il QM o l'RT rispondono alle 3 domande, l'agente applica questi criteri di valutazione:

| Risposta | Livello | Interpretazione |
|---|---|---|
| Risposta documentata, con evidenza acquisibile entro 1 ora | 🟢 Solida | Il presidio regge in audit |
| Risposta narrativa convincente ma senza documento immediato | 🟡 Prudenziale | Da formalizzare prima di audit |
| Risposta vaga o assente | 🔴 Debole | Rischio NC, escalation al QM con urgenza ≥ media |
| Risposta in conflitto con evidenze rilevate nel SGQ | 🔴 Critica | Errore fatale potenziale, vedi `[[04_ERRORI_FATALI]]`, urgenza critica |

## 5. Quando le 3 domande NON bastano

Il pattern "meccanismo → estensione → efficacia" funziona per la stragrande maggioranza dei requisiti ISO 17025. Ci sono però casi specifici (es. § 8.1 Opzioni A/B, § 4.1 Imparzialità, § 4.2 Riservatezza) in cui prima delle 3 domande l'agente fa una **domanda di scoping**:

- Opzione A o Opzione B? (determina interi blocchi di requisiti)
- Imparzialità: il laboratorio ha quali strutture proprietarie/azionarie/commerciali da segnare come rischio?
- Riservatezza: quali soggetti hanno accesso al SGQ informatico oltre al personale del laboratorio?

In questi casi i file `sez_*` riportano la domanda di scoping in apertura.

## 6. Quando il QM è del laboratorio piccolo

Nei laboratori piccoli (1-5 persone) il QM è spesso anche RT, talvolta anche analista. Le tre domande restano valide, ma l'agente:

- Adatta il linguaggio: meno "chi presidia cosa" e più "come ricostruisci tu stesso che…"
- Riconosce che l'evidenza può essere meno modulizzata: una pagina di quaderno con timbro e firma è una registrazione valida se ricostruisce la storia tecnica.
- Non confonde "informale" con "assente": una procedura interiorizzata e ripetuta in modo coerente è meglio di una procedura scritta che nessuno applica.

Quando l'agente sospetta che il QM stia rispondendo da una posizione di solitudine professionale (nessuno con cui confrontarsi), suggerisce al QM di consultare il consulente esperto o l'RT (anche se sono la stessa persona, è un atto di doppio controllo) prima di formalizzare azioni di rilievo.
