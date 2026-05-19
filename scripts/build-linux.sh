#!/usr/bin/env bash
# build-linux.sh — Costruisce il binario standalone labnexus per linux/amd64 (FR-11, OQ-6).
# Usato dal CTO su Windows WSL2 con NVIDIA RTX 4090 per Ollama accelerato via GPU.
set -euo pipefail

OUT_DIR="dist"
OUT_BIN="${OUT_DIR}/labnexus-linux-amd64"

mkdir -p "${OUT_DIR}"

echo "[build-linux] cross-compile linux/amd64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
  go build -trimpath -ldflags "-s -w" \
  -o "${OUT_BIN}" ./cmd/labnexus

echo "[build-linux] binario pronto: ${OUT_BIN}"
file "${OUT_BIN}" || true
