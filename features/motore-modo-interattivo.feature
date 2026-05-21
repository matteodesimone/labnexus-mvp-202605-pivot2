# language: it

Funzionalità: Motore LabNexus — modalità interattiva TUI cross-platform
  In quanto Denis (macOS) che lavora da Terminal con doppio click su labnexus.command,
  e in quanto CTO (Linux/WSL2 o macOS) che lancia da terminale,
  voglio una modalità interattiva sequenziale che giri uguale ovunque.

  Nota Sprint 1.5.A (post-pivot-3): il bundle .app è stato eliminato a favore di un
  binary CLI puro multipiattaforma + un piccolo launcher labnexus.command per macOS.
  Vedi features/sprint1.5-A-pivot3-cli-puro.feature per le verifiche del refactor.

  # --- FR-10: TUI cross-platform sequenziale ---

  @manual
  Scenario: lancio senza argomenti apre la TUI con selezione profilo
    Dato che lancio "labnexus" da terminale senza argomenti
    Allora la TUI mostra una lista selezionabile dei profili disponibili
    E accanto a ogni profilo è visibile la descrizione di una riga
    E posso navigare con frecce e confermare con invio

  @manual
  Scenario: dopo aver scelto il profilo, la TUI chiede la cartella di input
    Dato che ho selezionato "revisione" nella prima schermata della TUI
    Allora la TUI chiede "Cartella di input"
    E accetta un path digitato a mano o autocompletato

  @manual
  Scenario: dopo l'input, la TUI chiede la cartella di output
    Dato che ho confermato un path di input valido
    Allora la TUI chiede "Cartella di output"
    E posso confermare un path esistente o crearne uno nuovo

  @manual
  Scenario: un argomento posizionale = cartella esistente è usato come input
    Dato che lancio "labnexus /Users/denis/Test1"
    E "/Users/denis/Test1" è una cartella esistente
    Allora la TUI parte direttamente dalla schermata di selezione profilo
    E dopo la scelta del profilo chiede SOLO la cartella di output
    E l'input usato è "/Users/denis/Test1"

  Scenario: lancio in non-TTY rifiuta la modalità interattiva
    Dato che lancio "labnexus | cat" da uno script (stdout/stderr non sono TTY)
    Allora exit code è 2
    E stderr contiene "modalità interattiva non disponibile"
    E stderr suggerisce di "usare labnexus run con i flag"

  @manual
  Scenario: TUI funziona identicamente su darwin/arm64 e su linux/amd64
    Dato che lo stesso codice è compilato per "darwin/arm64" e "linux/amd64"
    Quando lancio "labnexus" su entrambi i target
    Allora il flusso TUI (selezione profilo → input → output) è identico
    E nessuna funzionalità della TUI dipende da chiamate osascript o API macOS

  # --- FR-11 amendment Sprint 1.5.A: binari standalone su entrambi i target ---

  @manual
  Scenario: binario Linux standalone
    Dato che il deliverable per linux/amd64 è il binario "labnexus" (senza bundle)
    Quando un utente Linux scarica e fa "chmod +x labnexus && ./labnexus"
    Allora la TUI parte uguale a macOS
    E non sono necessarie operazioni Gatekeeper-equivalenti

  @manual
  Scenario: lancio CLI puro funziona identicamente su entrambi i target
    Quando lancio "labnexus run --profile revisione --input /tmp/in --output /tmp/out"
    Allora l'esecuzione parte direttamente, nessuna TUI viene aperta
    E exit code e log sono gli stessi su darwin/arm64 e linux/amd64
