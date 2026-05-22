#!/usr/bin/env bash
# build-zip.sh — Pipeline di consegna per Denis (SHIP_CMD di config.sh).
# Sprint 1.5.A (post-pivot-3): zip layout senza .app bundle.
# Sprint 1.5.B: include labnexus.config.toml master.
# Sprint 1.5.C: include lavori/ (rinominata da data/) per concierge mode.
# Sprint 1.5.D: include bin/pandoc + bin/typst per rendering PDF accoppiato.
#               Inoltre popola direttamente dist/labnexus-sprint1-darwin-arm64/
#               (oltre allo zip) così sviluppo locale e smoke usano la cartella
#               già pronta all'uso senza dover estrarre il zip.
#
# Contenuto del deliverable:
#   labnexus                          (binary Mach-O arm64 standalone)
#   labnexus.command                  (launcher bash per doppio click Finder)
#   labnexus.config.toml              (config master TOML — Denis edita eurouter_api_key, Sprint 1.5.B)
#   GUIDA.md                          (Sprint 1.5.C: micro guida essenziale Denis zero-CLI)
#   profili/                          (7 profili shippable Sprint 1, formato .toml Sprint 1.5.B)
#   KB-ispettore/                     (symlink-resolved KB Denis)
#   lavori/                           (Sprint 1.5.C: 7 cartelle CAPABILITY A-G con _labnexus.toml metadata
#                                      per auto-discovery `labnexus jobs`)
#   bin/pandoc                        (Sprint 1.5.D: arm64 native ~180MB)
#   bin/typst                         (Sprint 1.5.D: arm64 native ~40MB)
#   README.md                         (italiano, concierge tone Sprint 1.5.C)
#   meta-prompt-genera-profilo.md     (FR-20, deliverable feat-meta Sprint 1; FR-22 update TOML schema Sprint 1.5.B)
#   guida-meta-prompt-denis.md        (FR-21)
#
# NO goreleaser (Sprint 1, binario non firmato — Gatekeeper-binary workaround
# documentato nel README: `xattr -d com.apple.quarantine labnexus labnexus.command`).
set -euo pipefail

DIST_DIR="dist"
PKG_NAME="labnexus-sprint1-darwin-arm64"
DIST_PKG="${DIST_DIR}/${PKG_NAME}"  # cartella deliverable popolata
ZIP_PATH="${DIST_DIR}/${PKG_NAME}.zip"

bash scripts/build-mac.sh

# Verifica che pdf-tools sia stato eseguito (i binari devono già essere in
# ${DIST_PKG}/bin/). Caller deve aver fatto `make pdf-tools` prima.
if [[ ! -x "${DIST_PKG}/bin/pandoc" || ! -x "${DIST_PKG}/bin/typst" ]]; then
  echo "[build-zip] errore: pandoc/typst mancanti in ${DIST_PKG}/bin/." >&2
  echo "            Esegui 'make pdf-tools' prima di 'make ship/package'." >&2
  exit 1
fi

# Cleanup tutto in DIST_PKG TRANNE bin/ (preserva pandoc/typst scaricati da pdf-tools)
find "${DIST_PKG}" -mindepth 1 -maxdepth 1 ! -name "bin" -exec rm -rf {} +

# Popola DIST_PKG con il contenuto del deliverable
cp ./bin/labnexus-darwin-arm64 "${DIST_PKG}/labnexus"
chmod +x "${DIST_PKG}/labnexus"
cp ./scripts/labnexus.command "${DIST_PKG}/labnexus.command"
chmod +x "${DIST_PKG}/labnexus.command"
cp ./labnexus.config.toml "${DIST_PKG}/labnexus.config.toml"
cp ./GUIDA.md "${DIST_PKG}/GUIDA.md"
cp ./README.md "${DIST_PKG}/README.md"
cp ./docs/meta-prompt-genera-profilo.md "${DIST_PKG}/"
cp ./docs/guida-meta-prompt-denis.md "${DIST_PKG}/"

# Cartelle: copia ricorsiva con dereferenziazione dei symlink (KB-ispettore può
# essere un symlink al repo principale; Denis deve ricevere i file reali).
cp -RL ./profili "${DIST_PKG}/profili"
cp -RL ./KB-ispettore "${DIST_PKG}/KB-ispettore"
cp -RL ./lavori "${DIST_PKG}/lavori"

# Crea lo zip dalla cartella popolata. `cd` nel DIST_DIR per preservare la
# struttura "labnexus-sprint1-darwin-arm64/..." dentro al zip.
rm -f "${ZIP_PATH}"
(cd "${DIST_DIR}" && zip -rq "${PKG_NAME}.zip" "${PKG_NAME}")

echo "[build-zip] cartella pronta all'uso: ${DIST_PKG}/"
echo "[build-zip] zip di consegna:         ${ZIP_PATH} ($(du -h "${ZIP_PATH}" | cut -f1))"
