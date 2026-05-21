#!/usr/bin/env bash
# build-zip.sh — Pipeline di consegna per Denis (SHIP_CMD di config.sh).
# Costruisce labnexus.app + zip con profili + KB-ispettore + README +
# 2 docs deliverable Sprint 1 Fetta 3 (meta-prompt + guida d'uso, FR-20/21).
# NO goreleaser (Sprint 1, binario non firmato — vedi OOS-5).
set -euo pipefail

DIST_DIR="dist"
ZIP_NAME="labnexus-sprint1-darwin-arm64.zip"

bash scripts/build-mac.sh

rm -rf "${DIST_DIR}/${ZIP_NAME}"
mkdir -p "${DIST_DIR}"

# Stage temporaneo per i due docs deliverable (Sprint 1 Fetta 3): vanno
# nel root dello zip così Denis li trova accanto al README, non sepolti
# in docs/. Pattern coerente con README già al root.
STAGE_DIR="$(mktemp -d)"
trap 'rm -rf "${STAGE_DIR}"' EXIT
cp docs/meta-prompt-genera-profilo.md "${STAGE_DIR}/"
cp docs/guida-meta-prompt-denis.md "${STAGE_DIR}/"

# IMPORTANTE: niente -y. `zip -r` (senza -y) FOLLOW i symlink e include i file
# reali della KB-ispettore. Con -y il symlink verrebbe preservato e Denis
# riceverebbe un link broken al primo extract.
zip -r "${DIST_DIR}/${ZIP_NAME}" \
  labnexus.app \
  profili \
  KB-ispettore \
  README.md

# Aggiunge i 2 docs deliverable al root dello zip (`-j` = no path).
zip -j "${DIST_DIR}/${ZIP_NAME}" \
  "${STAGE_DIR}/meta-prompt-genera-profilo.md" \
  "${STAGE_DIR}/guida-meta-prompt-denis.md"

echo "[build-zip] pacchetto pronto: ${DIST_DIR}/${ZIP_NAME}"
unzip -l "${DIST_DIR}/${ZIP_NAME}" | tail -12
