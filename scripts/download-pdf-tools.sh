#!/usr/bin/env bash
# download-pdf-tools.sh — Sprint 1.5.D
#
# Scarica pandoc + typst nei binari Mac arm64 nativi e li posiziona in
# dist/labnexus-sprint1-darwin-arm64/bin/. Idempotente: skip se già presenti.
#
# Stack PDF rendering (pivot da goldmark+fpdf):
#   labnexus → pandoc (MD → Typst markup) → typst (Typst → PDF nativo)
#
# Disponibilità binari Mac Apple Silicon native (no Rosetta):
#   - pandoc 3.9.0.2: arm64-macOS.zip (nativo)
#   - typst 0.14.2: typst-aarch64-apple-darwin.tar.xz (nativo Rust)

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DIST_BIN="${ROOT}/dist/labnexus-sprint1-darwin-arm64/bin"
CACHE="${ROOT}/.cache/pdf-tools"

PANDOC_VERSION="3.9.0.2"
PANDOC_URL="https://github.com/jgm/pandoc/releases/download/${PANDOC_VERSION}/pandoc-${PANDOC_VERSION}-arm64-macOS.zip"

TYPST_VERSION="0.14.2"
TYPST_URL="https://github.com/typst/typst/releases/download/v${TYPST_VERSION}/typst-aarch64-apple-darwin.tar.xz"

mkdir -p "${DIST_BIN}" "${CACHE}"

# --- pandoc ---
if [[ -x "${DIST_BIN}/pandoc" ]]; then
  echo "[pdf-tools] pandoc già presente in ${DIST_BIN}/pandoc — skip"
else
  echo "[pdf-tools] download pandoc ${PANDOC_VERSION}..."
  PANDOC_ZIP="${CACHE}/pandoc-${PANDOC_VERSION}.zip"
  if [[ ! -f "${PANDOC_ZIP}" ]]; then
    curl -L --fail "${PANDOC_URL}" -o "${PANDOC_ZIP}"
  fi
  PANDOC_EXTRACT="${CACHE}/pandoc-${PANDOC_VERSION}-extract"
  rm -rf "${PANDOC_EXTRACT}"
  unzip -q "${PANDOC_ZIP}" -d "${PANDOC_EXTRACT}"
  PANDOC_BIN="$(find "${PANDOC_EXTRACT}" -name pandoc -type f -perm +111 | head -1)"
  if [[ -z "${PANDOC_BIN}" ]]; then
    echo "[pdf-tools] errore: binario pandoc non trovato in ${PANDOC_EXTRACT}" >&2
    exit 1
  fi
  cp "${PANDOC_BIN}" "${DIST_BIN}/pandoc"
  chmod +x "${DIST_BIN}/pandoc"
  # Rimuovi quarantine flag così macOS Gatekeeper non blocca
  xattr -d com.apple.quarantine "${DIST_BIN}/pandoc" 2>/dev/null || true
  echo "[pdf-tools] pandoc installato in ${DIST_BIN}/pandoc ($(du -h "${DIST_BIN}/pandoc" | cut -f1))"
fi

# --- typst ---
if [[ -x "${DIST_BIN}/typst" ]]; then
  echo "[pdf-tools] typst già presente in ${DIST_BIN}/typst — skip"
else
  echo "[pdf-tools] download typst ${TYPST_VERSION}..."
  TYPST_TAR="${CACHE}/typst-${TYPST_VERSION}.tar.xz"
  if [[ ! -f "${TYPST_TAR}" ]]; then
    curl -L --fail "${TYPST_URL}" -o "${TYPST_TAR}"
  fi
  TYPST_EXTRACT="${CACHE}/typst-${TYPST_VERSION}-extract"
  rm -rf "${TYPST_EXTRACT}"
  mkdir -p "${TYPST_EXTRACT}"
  tar -xJf "${TYPST_TAR}" -C "${TYPST_EXTRACT}"
  TYPST_BIN="$(find "${TYPST_EXTRACT}" -name typst -type f -perm +111 | head -1)"
  if [[ -z "${TYPST_BIN}" ]]; then
    echo "[pdf-tools] errore: binario typst non trovato in ${TYPST_EXTRACT}" >&2
    exit 1
  fi
  cp "${TYPST_BIN}" "${DIST_BIN}/typst"
  chmod +x "${DIST_BIN}/typst"
  xattr -d com.apple.quarantine "${DIST_BIN}/typst" 2>/dev/null || true
  echo "[pdf-tools] typst installato in ${DIST_BIN}/typst ($(du -h "${DIST_BIN}/typst" | cut -f1))"
fi

echo "[pdf-tools] verifica versioni:"
"${DIST_BIN}/pandoc" --version | head -1
"${DIST_BIN}/typst" --version 2>&1 | head -1

echo "[pdf-tools] done"
