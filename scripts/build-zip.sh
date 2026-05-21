#!/usr/bin/env bash
# build-zip.sh — Pipeline di consegna per Denis (SHIP_CMD di config.sh).
# Sprint 1.5.A (post-pivot-3): zip layout senza .app bundle.
# Sprint 1.5.B: include labnexus.config.toml master.
# Sprint 1.5.C: include lavori/ (rinominata da data/) per concierge mode.
# Contenuto del deliverable:
#   labnexus                          (binary Mach-O arm64 standalone)
#   labnexus.command                  (launcher bash per doppio click Finder)
#   labnexus.config.toml              (config master TOML — Denis edita eurouter_api_key, Sprint 1.5.B)
#   GUIDA.md                          (Sprint 1.5.C: micro guida essenziale Denis zero-CLI)
#   profili/                          (7 profili shippable Sprint 1, formato .toml Sprint 1.5.B)
#   KB-ispettore/                     (symlink-resolved KB Denis)
#   lavori/                           (Sprint 1.5.C: 7 cartelle CAPABILITY A-G con _labnexus.toml metadata
#                                      per auto-discovery `labnexus jobs`)
#   README.md                         (italiano, concierge tone Sprint 1.5.C)
#   meta-prompt-genera-profilo.md     (FR-20, deliverable feat-meta Sprint 1; FR-22 update TOML schema Sprint 1.5.B)
#   guida-meta-prompt-denis.md        (FR-21)
#
# NO goreleaser (Sprint 1, binario non firmato — Gatekeeper-binary workaround
# documentato nel README: `xattr -d com.apple.quarantine labnexus labnexus.command`).
set -euo pipefail

DIST_DIR="dist"
ZIP_NAME="labnexus-sprint1-darwin-arm64.zip"

bash scripts/build-mac.sh

rm -rf "${DIST_DIR}/${ZIP_NAME}"
mkdir -p "${DIST_DIR}"

# Stage temporaneo per assemblare il layout target dello zip al root:
#   labnexus (binary rinominato da bin/labnexus-darwin-arm64)
#   labnexus.command (copia da scripts/)
#   profili/, KB-ispettore/, README.md (link diretti)
#   docs/meta-prompt-genera-profilo.md, docs/guida-meta-prompt-denis.md → root
STAGE_DIR="$(mktemp -d)"
trap 'rm -rf "${STAGE_DIR}"' EXIT

cp ./bin/labnexus-darwin-arm64 "${STAGE_DIR}/labnexus"
chmod +x "${STAGE_DIR}/labnexus"
cp ./scripts/labnexus.command "${STAGE_DIR}/labnexus.command"
chmod +x "${STAGE_DIR}/labnexus.command"
cp ./docs/meta-prompt-genera-profilo.md "${STAGE_DIR}/"
cp ./docs/guida-meta-prompt-denis.md "${STAGE_DIR}/"
# Sprint 1.5.B: config master TOML. Denis edita `eurouter_api_key` post-extract.
cp ./labnexus.config.toml "${STAGE_DIR}/labnexus.config.toml"
# Sprint 1.5.C: micro guida essenziale Denis (zero-CLI, drag&drop Finder, dove
# mettere la chiave EUROUTER spiegato esplicitamente). Vedi GUIDA.md al repo root.
cp ./GUIDA.md "${STAGE_DIR}/GUIDA.md"

# IMPORTANTE: niente -y. `zip -r` (senza -y) FOLLOW i symlink e include i file
# reali della KB-ispettore. Con -y il symlink verrebbe preservato e Denis
# riceverebbe un link broken al primo extract.
zip -j "${DIST_DIR}/${ZIP_NAME}" \
  "${STAGE_DIR}/labnexus" \
  "${STAGE_DIR}/labnexus.command" \
  "${STAGE_DIR}/labnexus.config.toml" \
  "${STAGE_DIR}/GUIDA.md" \
  "${STAGE_DIR}/meta-prompt-genera-profilo.md" \
  "${STAGE_DIR}/guida-meta-prompt-denis.md"

# Le cartelle si aggiungono con -r (no -j perché vogliamo preservare la
# struttura profili/, KB-ispettore/, lavori/).
# Sprint 1.5.C: include `lavori/` (rinominata da `data/`) con i 7 _labnexus.toml
# per auto-discovery (FR-25).
zip -r "${DIST_DIR}/${ZIP_NAME}" \
  profili \
  KB-ispettore \
  lavori \
  README.md

echo "[build-zip] pacchetto pronto: ${DIST_DIR}/${ZIP_NAME}"
unzip -l "${DIST_DIR}/${ZIP_NAME}" | tail -12
