#!/usr/bin/env bash
# Esegui.command — TEMPLATE per un NUOVO lavoro (Sprint 1.5.D, modalità wizard).
#
# Flusso: copia questa cartella, rinominala, mettici i tuoi file dati, doppio
# click. Se manca _labnexus.toml, parte il wizard `labnexus init` (scegli
# profilo + prompt), poi esegue. Vedi LEGGIMI.txt.

set -e

DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT="$(cd "$DIR/../.." && pwd)"
JOB_NAME="$(basename "$DIR")"

# Guard: non eseguire il TEMPLATE stesso — va prima copiato e rinominato.
case "$JOB_NAME" in
  *NUOVO\ LAVORO*|*copiami*)
    echo ""
    echo "═══════════════════════════════════════════════════════════════════"
    echo "  Questa è la cartella TEMPLATE — non eseguirla direttamente."
    echo "═══════════════════════════════════════════════════════════════════"
    echo ""
    echo "  1. Copiala  (tasto destro → Duplica)"
    echo "  2. Rinomina la copia col nome del tuo lavoro"
    echo "  3. Mettici dentro i tuoi file dati"
    echo "  4. Doppio click su Esegui.command nella COPIA"
    echo ""
    read -n 1 -s -r -p "Premi un tasto per chiudere la finestra..."
    exit 0
    ;;
esac

cd "$ROOT" || {
  echo ""
  echo "ERRORE: impossibile entrare nella cartella di labnexus."
  read -n 1 -s -r -p "Premi un tasto per chiudere la finestra..."
  exit 2
}

# Auto-heal del bit eseguibile: cloud-sync (kDrive/Dropbox/OneDrive) azzera +x.
if [ -f "./labnexus" ] && [ ! -x "./labnexus" ]; then chmod +x ./labnexus 2>/dev/null; fi

# Pre-check API key (FR-36, stesso check degli altri launcher)
if [ -f "labnexus.config.toml" ] && grep -qE '^eurouter_api_key\s*=\s*""\s*$' labnexus.config.toml; then
  if [ -z "$EUROUTER_API_KEY" ]; then
    echo ""
    echo "  EUROUTER_API_KEY mancante. Apri 'labnexus.config.toml' (cartella"
    echo "  principale) e inserisci la chiave alla riga eurouter_api_key."
    echo ""
    read -n 1 -s -r -p "Premi un tasto per chiudere la finestra..."
    exit 2
  fi
fi

# Se manca il metadata, lancia il wizard interattivo.
if [ ! -f "$DIR/_labnexus.toml" ]; then
  code=0
  ./labnexus init "$DIR" --profiles-dir profili || code=$?
  if [ "$code" -eq 10 ]; then
    echo ""
    echo "Hai creato un nuovo prompt: aprilo, scrivi le istruzioni, salva,"
    echo "poi rilancia questo Esegui.command."
    echo ""
    read -n 1 -s -r -p "Premi un tasto per chiudere la finestra..."
    exit 0
  fi
  if [ "$code" -ne 0 ]; then
    read -n 1 -s -r -p "Premi un tasto per chiudere la finestra..."
    exit "$code"
  fi
fi

echo ""
echo "═══════════════════════════════════════════════════════════════════"
echo "  Lancio: $JOB_NAME"
echo "═══════════════════════════════════════════════════════════════════"
echo ""

./labnexus run --job "$JOB_NAME" --lavori-dir lavori --profiles-dir profili --config labnexus.config.toml --kb-dir KB-ispettore || {
  echo ""
  echo "ERRORE: esecuzione fallita. Vedi i messaggi sopra per il dettaglio."
  read -n 1 -s -r -p "Premi un tasto per chiudere la finestra..."
  exit 1
}

echo ""
echo "═══════════════════════════════════════════════════════════════════"
echo "  Completato. Output salvato in: $DIR/output/"
echo "═══════════════════════════════════════════════════════════════════"
echo ""
read -n 1 -s -r -p "Premi un tasto per chiudere la finestra..."
