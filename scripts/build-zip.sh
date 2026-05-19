#!/usr/bin/env bash
# build-zip.sh — Pipeline di consegna per Denis (SHIP_CMD di config.sh).
# Costruisce labnexus.app + zip con profili + KB-ispettore + README.
# NO goreleaser (Sprint 1, binario non firmato — vedi OOS-5).
set -euo pipefail

DIST_DIR="dist"
ZIP_NAME="labnexus-sprint1-darwin-arm64.zip"

bash scripts/build-mac.sh

rm -rf "${DIST_DIR}/${ZIP_NAME}"
mkdir -p "${DIST_DIR}"

# La KB-ispettore è un symlink → la zip risolve seguendo il link.
zip -ry "${DIST_DIR}/${ZIP_NAME}" \
  labnexus.app \
  profili \
  KB-ispettore \
  README.md

echo "[build-zip] pacchetto pronto: ${DIST_DIR}/${ZIP_NAME}"
unzip -l "${DIST_DIR}/${ZIP_NAME}" | tail -10
