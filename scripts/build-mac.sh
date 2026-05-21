#!/usr/bin/env bash
# build-mac.sh — Cross-compile darwin/arm64 standalone binary (Sprint 1.5.A).
#
# Sprint 1.5.A (post-pivot-3): nessun bundle, niente wrapper bash, niente
# AppleScript Terminal launcher, niente file plist di bundle. Solo un singolo
# binary Mach-O arm64.
# Lo zip di consegna include il binary + scripts/labnexus.command (launcher
# bash 2 righe per il doppio click dal Finder).
#
# Vedi .pipeline/archive/2026-05-21-fetta3-*.md per il razionale del refactor.
set -euo pipefail

OUTDIR="bin"
OUTFILE="${OUTDIR}/labnexus-darwin-arm64"

mkdir -p "${OUTDIR}"

echo "[build-mac] cross-compile darwin/arm64 → ${OUTFILE}..."
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 \
  go build -trimpath -ldflags "-s -w" \
  -o "${OUTFILE}" ./cmd/labnexus

echo "[build-mac] binary pronto: $(ls -lh "${OUTFILE}" | awk '{print $5}')"
file "${OUTFILE}"
