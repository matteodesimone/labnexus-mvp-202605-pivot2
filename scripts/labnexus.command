#!/usr/bin/env bash
# labnexus.command — minimal launcher per macOS (Sprint 1.5.A + 1.5.C).
# Doppio click dal Finder apre Terminal e lancia il binary labnexus presente
# nella stessa cartella. La TUI charmbracelet/huh parte automaticamente se
# stdin/stdout sono TTY (è sempre il caso dentro Terminal).
#
# Sprint 1.5.C FR-36: pre-check API key shell-side. Se labnexus.config.toml
# contiene `eurouter_api_key = ""` (default zip iniziale) E env EUROUTER_API_KEY
# non setted, blocca prima del binary fork con messaggio chiaro in italiano.

cd "$(dirname "$0")" || {
  echo ""
  echo "ERRORE: impossibile entrare nella cartella di labnexus.command."
  echo ""
  read -n 1 -s -r -p "Premi un tasto per chiudere la finestra..."
  echo ""
  exit 2
}

# FR-36 pre-check API key
if [ -f "labnexus.config.toml" ] && grep -qE '^eurouter_api_key\s*=\s*""\s*$' labnexus.config.toml; then
  if [ -z "$EUROUTER_API_KEY" ]; then
    echo ""
    echo "═══════════════════════════════════════════════════════════════════"
    echo "  EUROUTER_API_KEY mancante"
    echo "═══════════════════════════════════════════════════════════════════"
    echo ""
    echo "Apri 'labnexus.config.toml' con un editor di testo e inserisci la"
    echo "tua chiave EUROUTER alla riga:"
    echo ""
    echo "    eurouter_api_key = \"sk-...\""
    echo ""
    echo "Poi rilancia labnexus.command con doppio click."
    echo ""
    read -n 1 -s -r -p "Premi un tasto per chiudere la finestra..."
    echo ""
    exit 2
  fi
fi

# Il binary deve esistere accanto al launcher.
if [ ! -f "./labnexus" ]; then
  echo ""
  echo "ERRORE: binary 'labnexus' non trovato accanto a labnexus.command."
  echo "Verifica che labnexus.command e labnexus siano nella stessa cartella."
  echo ""
  read -n 1 -s -r -p "Premi un tasto per chiudere la finestra..."
  echo ""
  exit 2
fi

# Auto-heal del bit eseguibile: i servizi di cloud-sync (kDrive, Dropbox,
# OneDrive, Google Drive) spesso azzerano il permesso +x sui file sincronizzati,
# causando "Permission denied" al doppio click. Lo ripristiniamo silenziosamente.
if [ ! -x "./labnexus" ]; then
  chmod +x ./labnexus 2>/dev/null
fi
if [ ! -x "./labnexus" ]; then
  echo ""
  echo "ERRORE: 'labnexus' non è eseguibile e non sono riuscito a ripristinare il"
  echo "permesso (chmod fallito). Apri Terminal nella cartella e lancia:"
  echo ""
  echo "    chmod +x labnexus labnexus.command"
  echo ""
  read -n 1 -s -r -p "Premi un tasto per chiudere la finestra..."
  echo ""
  exit 2
fi

./labnexus || {
  echo ""
  echo "ERRORE: binary 'labnexus' uscito con errore (crash o exit non-zero)."
  echo "Se il problema persiste, contatta il referente tecnico."
  echo ""
  read -n 1 -s -r -p "Premi un tasto per chiudere la finestra..."
  echo ""
}
