# Makefile — LabNexus Sprint 1
# Convenzione: `make` (senza target) stampa l'help. `make ship` produce il
# pacchetto completo di consegna per Denis (build + test + zip).
# Tutti i comandi long-running sono safe-da-rilanciare (idempotenti).

SHELL       := /usr/bin/env bash
.SHELLFLAGS := -eu -o pipefail -c

# ─── Configurazione ──────────────────────────────────────────────────────────

BINARY      := labnexus
MAIN_PKG    := ./cmd/labnexus
DIST_DIR    := dist
ZIP_NAME    := labnexus-sprint1-darwin-arm64.zip
# Sprint 1.5.A (post-pivot-3): rimosso APP_BUNDLE (no più .app bundle).

# Versione: usa `git describe` se disponibile, altrimenti "dev"
VERSION     := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS     := -s -w -X main.version=$(VERSION)

GO          := go
GOFLAGS     := -trimpath

# Materiale per smoke test (Test 1 reale di Denis, presente nel repo)
SMOKE_INPUT := docs/piano_iniziale/materiali-dominio/test-precedenti/test-1-input

.DEFAULT_GOAL := help

# ─── Help (target di default) ────────────────────────────────────────────────

.PHONY: help
help: ## Mostra la lista dei target disponibili
	@echo ""
	@echo "LabNexus — Makefile (versione: $(VERSION))"
	@echo ""
	@echo "Target principali:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z][a-zA-Z0-9_-]*:.*?## / {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Workflow tipico:"
	@echo "  make test       # tutti i test (unit + integration + BDD)"
	@echo "  make ship       # pacchetto di consegna completo (clean + test + zip)"
	@echo ""

# ─── Build ───────────────────────────────────────────────────────────────────

.PHONY: build
build: ## Build per la macchina corrente (in ./bin/labnexus)
	@mkdir -p bin
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(BINARY) $(MAIN_PKG)
	@echo "✓ bin/$(BINARY) ($$(file bin/$(BINARY) | cut -d: -f2 | xargs))"

.PHONY: build-mac
build-mac: ## Cross-compile macOS Apple Silicon → bin/labnexus-darwin-arm64 (binary standalone, no bundle Sprint 1.5.A)
	bash scripts/build-mac.sh

.PHONY: build-linux
build-linux: ## Cross-compile Linux amd64 → dist/labnexus-linux-amd64
	bash scripts/build-linux.sh

.PHONY: build-all
build-all: build-mac build-linux ## Build per entrambi i target Sprint 1 (darwin/arm64 + linux/amd64)

.PHONY: pdf-tools
pdf-tools: ## Scarica pandoc + typst (arm64 native) in dist/.../bin/ per rendering PDF (Sprint 1.5.D)
	bash scripts/download-pdf-tools.sh

# ─── Test ────────────────────────────────────────────────────────────────────

.PHONY: test
test: ## Tutti i test (unit + integration + BDD acceptance)
	$(GO) test ./...

.PHONY: test-unit
test-unit: ## Solo unit/integration test dei package internal/
	$(GO) test ./internal/...

.PHONY: test-bdd
test-bdd: ## Solo BDD acceptance (godog) — esclude scenari @manual
	$(GO) test ./features/

.PHONY: test-verbose
test-verbose: ## Test con output verbose (utile per debug fallimenti)
	$(GO) test -v ./...

# ─── Lint / Sanity check ─────────────────────────────────────────────────────

.PHONY: vet
vet: ## go vet ./... (analisi statica built-in)
	$(GO) vet ./...

.PHONY: fmt
fmt: ## go fmt ./... (formattazione idiomatica)
	$(GO) fmt ./...

.PHONY: lint
lint: ## golangci-lint (richiede install separato, vedi DEV-GUIDE.md)
	@command -v golangci-lint >/dev/null 2>&1 || { echo "✗ golangci-lint non installato. Vedi docs/DEV-GUIDE.md → Linting"; exit 1; }
	golangci-lint run

.PHONY: tidy
tidy: ## go mod tidy (pulisce go.mod / go.sum)
	$(GO) mod tidy

# ─── Smoke test (binario reale su Test 1) ────────────────────────────────────

.PHONY: smoke
smoke: build ## Smoke test del binario: list + describe + check su Test 1 reale (dry-run)
	@# Nota: bin/labnexus risolve i default profili/KB relativi al binario (bin/profili
	@# non esiste); per il dev workflow puntiamo esplicitamente a quelli della repo.
	@echo ""
	@echo "→ labnexus list"
	@./bin/$(BINARY) --profiles-dir profili --kb-dir KB-ispettore list
	@echo ""
	@echo "→ labnexus describe revisione"
	@./bin/$(BINARY) --profiles-dir profili --kb-dir KB-ispettore describe revisione
	@echo ""
	@echo "→ labnexus validate revisione"
	@./bin/$(BINARY) --profiles-dir profili --kb-dir KB-ispettore validate revisione
	@echo ""
	@echo "→ labnexus check revisione --input <test1-input>  (dry-run, no LLM)"
	@./bin/$(BINARY) --profiles-dir profili --kb-dir KB-ispettore check revisione --input "$(SMOKE_INPUT)" 2>&1 | tail -10

# ─── Packaging / Ship ────────────────────────────────────────────────────────

.PHONY: package
package: build-mac ## Costruisce binary darwin/arm64 + zip di consegna (dist/labnexus-sprint1-darwin-arm64.zip)
	bash scripts/build-zip.sh

.PHONY: ship
ship: clean pdf-tools test vet package ## Pipeline completa di consegna: clean → pdf-tools → test → vet → build-mac → zip
	@echo ""
	@echo "✓ Pacchetto di consegna pronto:"
	@ls -lh $(DIST_DIR)/$(ZIP_NAME)
	@echo ""
	@echo "Contenuto:"
	@unzip -l $(DIST_DIR)/$(ZIP_NAME) | tail -15

# ─── Cleanup ─────────────────────────────────────────────────────────────────

.PHONY: clean
clean: ## Rimuove artefatti di build (bin/, dist/)
	rm -rf bin/
	rm -rf $(DIST_DIR)/
	rm -rf labnexus.app/   # cleanup legacy bundle se presente da pre-Sprint-1.5.A
	rm -f labnexus-smoke labnexus-final
	@echo "✓ artefatti rimossi"

.PHONY: distclean
distclean: clean ## clean + rimuove cache Go locale del modulo
	$(GO) clean -testcache
	@echo "✓ cache test pulita"

# ─── CI shortcut ─────────────────────────────────────────────────────────────

.PHONY: ci
ci: tidy vet test ## Shortcut da CI: tidy + vet + test (no build)
	@echo "✓ CI OK"
