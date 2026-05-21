#!/usr/bin/env bash
# labnexus.command — minimal launcher per macOS (Sprint 1.5.A, post-pivot-3).
# Doppio click dal Finder apre Terminal e lancia il binary labnexus presente
# nella stessa cartella. La TUI charmbracelet/huh parte automaticamente se
# stdin/stdout sono TTY (è sempre il caso dentro Terminal).
#
# Vedi README.md per il setup pre-volo (EUROUTER_API_KEY come env var in
# Sprint 1.5.A; futuro labnexus.config.toml in Sprint 1.5.B).
cd "$(dirname "$0")" && ./labnexus || {
  echo ""
  echo "ERRORE: binary 'labnexus' non trovato accanto a labnexus.command."
  echo "Posiziona labnexus.command e labnexus nella stessa cartella."
  echo ""
  read -n 1 -s -r -p "Premi un tasto per chiudere la finestra..."
  echo ""
}
