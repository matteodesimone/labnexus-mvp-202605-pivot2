# language: it

Funzionalità: Motore LabNexus — modalità interattiva TUI cross-platform e bundle .app
  In quanto Denis (macOS) che lavora dal Finder, e in quanto CTO (Linux/WSL2 o macOS) che lancia da terminale,
  voglio una modalità interattiva sequenziale che giri uguale ovunque,
  e su macOS un bundle .app che apra Terminal con la TUI in modo trasparente.

  # --- FR-10: TUI cross-platform sequenziale ---

  Scenario: lancio senza argomenti apre la TUI con selezione profilo
    Dato che lancio "labnexus" da terminale senza argomenti
    Allora la TUI mostra una lista selezionabile dei profili disponibili
    E accanto a ogni profilo è visibile la descrizione di una riga
    E posso navigare con frecce e confermare con invio

  Scenario: dopo aver scelto il profilo, la TUI chiede la cartella di input
    Dato che ho selezionato "revisione" nella prima schermata della TUI
    Allora la TUI chiede "Cartella di input"
    E accetta un path digitato a mano o autocompletato

  Scenario: dopo l'input, la TUI chiede la cartella di output
    Dato che ho confermato un path di input valido
    Allora la TUI chiede "Cartella di output"
    E posso confermare un path esistente o crearne uno nuovo

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

  Scenario: TUI funziona identicamente su darwin/arm64 e su linux/amd64
    Dato che lo stesso codice è compilato per "darwin/arm64" e "linux/amd64"
    Quando lancio "labnexus" su entrambi i target
    Allora il flusso TUI (selezione profilo → input → output) è identico
    E nessuna funzionalità della TUI dipende da chiamate osascript o API macOS

  # --- FR-11: bundle .app macOS + binario Linux ---

  Scenario: bundle .app contiene l'eseguibile arm64 e Info.plist valido
    Quando ispeziono "labnexus.app"
    Allora esiste "labnexus.app/Contents/MacOS/labnexus" eseguibile arm64
    E "labnexus.app/Contents/Info.plist" dichiara CFBundleExecutable = "labnexus"
    E "Info.plist" dichiara LSHandlerRank e supporto a "Folder" (per drag&drop)

  Scenario: doppio click sul .app apre Terminal e avvia la TUI
    Dato che faccio doppio click su "labnexus.app" dal Finder
    Allora si apre una finestra Terminal
    E nella finestra parte la TUI sequenziale di FR-10
    E la finestra resta aperta a fine esecuzione per permettere di leggere log/output

  Scenario: drag&drop di una cartella sull'icona .app pre-seleziona l'input
    Quando trascino la cartella "/Users/denis/Test1" sull'icona "labnexus.app"
    Allora si apre Terminal con la TUI
    E la TUI parte dalla selezione del profilo (non chiede l'input)
    E "/Users/denis/Test1" è l'input usato

  Scenario: bundle non firmato → procedura Gatekeeper documentata
    Dato che "labnexus.app" non è firmato con certificato Apple Developer
    Quando un utente lancia per la prima volta
    Allora macOS blocca con il messaggio Gatekeeper standard
    E il README documenta la procedura "control-clic → Apri" come passo una tantum

  Scenario: binario Linux standalone
    Dato che il deliverable per linux/amd64 è il binario "labnexus" (no .app)
    Quando un utente Linux scarica e fa "chmod +x labnexus && ./labnexus"
    Allora la TUI parte uguale a macOS
    E non sono necessari .app bundle né operazioni Gatekeeper-equivalenti

  Scenario: lancio CLI puro funziona identicamente su entrambi i target
    Quando lancio "labnexus run --profile revisione --input /tmp/in --output /tmp/out"
    Allora l'esecuzione parte direttamente, nessuna TUI viene aperta
    E exit code e log sono gli stessi su darwin/arm64 e linux/amd64
