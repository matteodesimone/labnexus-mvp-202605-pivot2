#!/usr/bin/env bash
# Esegui.command — launcher per-capability (Sprint 1.5.C concierge mode).
#
# Doppio click dal Finder → esegue automaticamente labnexus run su QUESTA
# cartella di lavoro. Niente domande, niente drag&drop. Il nome del job è
# derivato dal nome della cartella che contiene questo file (substring match
# via labnexus.jobs.Resolve case-insensitive).
#
# Per usare un job diverso: edita la riga JOB= sotto, o (raccomandato) duplica
# l'intera cartella di lavoro e lasciaglielo derivare dal nome.

set -e

# Risolvi i path: questo file sta in lavori/<X>/Esegui.command;
# il binary labnexus è 2 livelli sopra.
DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT="$(cd "$DIR/../.." && pwd)"
JOB_NAME="$(basename "$DIR")"

cd "$ROOT" || {
  echo ""
  echo "ERRORE: impossibile entrare nella cartella di labnexus."
  read -n 1 -s -r -p "Premi un tasto per chiudere la finestra..."
  exit 2
}

# Pre-check API key (FR-36, stesso check del labnexus.command root)
if [ -f "labnexus.config.toml" ] && grep -qE '^eurouter_api_key\s*=\s*""\s*$' labnexus.config.toml; then
  if [ -z "$EUROUTER_API_KEY" ]; then
    echo ""
    echo "═══════════════════════════════════════════════════════════════════"
    echo "  EUROUTER_API_KEY mancante"
    echo "═══════════════════════════════════════════════════════════════════"
    echo ""
    echo "Apri 'labnexus.config.toml' (cartella principale, 2 livelli sopra"
    echo "questa cartella) con un editor di testo e inserisci la tua chiave"
    echo "EUROUTER alla riga:"
    echo ""
    echo "    eurouter_api_key = \"sk-...\""
    echo ""
    echo "Poi rilancia questo file con doppio click."
    echo ""
    read -n 1 -s -r -p "Premi un tasto per chiudere la finestra..."
    exit 2
  fi
fi

echo ""
echo "═══════════════════════════════════════════════════════════════════"
echo "  Lancio: $JOB_NAME"
echo "═══════════════════════════════════════════════════════════════════"
echo ""

# Auto-heal del bit eseguibile: cloud-sync (kDrive/Dropbox/OneDrive) azzera +x.
if [ -f "./labnexus" ] && [ ! -x "./labnexus" ]; then chmod +x ./labnexus 2>/dev/null; fi

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
