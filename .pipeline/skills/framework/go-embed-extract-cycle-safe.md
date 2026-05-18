# Go embed-extract per binary self-provisioning (cycle-safe)

Pattern per embeddare un binario Go A dentro un binario B (stesso repo,
stesso modulo) quando A importa un package F che farebbe naturalmente
da "home" per la go:embed directive.

## Trigger

- Binario A deve essere distribuito dentro binario B
- A importa package F (per qualunque ragione — logica condivisa, version,
  config embedded, ecc.)
- Vogliamo usare `//go:embed` invece di download runtime o tool esterno

## Anti-trigger

- A e B sono indipendenti (no import di F da parte di A) — metti embed in F
- A è binario esterno (non compilato da stesso repo) — usa GoReleaser +
  download from releases, non go:embed

## Cycle

Naiv approach:

```go
// framework/embed.go:
//go:embed embed/bin/vobbly
var VobblyBinary []byte
```

Build vobbly (`go build ./cmd/vobbly`) → compiler deve compilare dependency
tree incluso `framework` → `framework/embed.go` ha embed directive che
richiede `framework/embed/bin/vobbly` esistente → ma quello è l'artifact
che stiamo building! **Error: pattern embed/bin/vobbly: no matching files**.

## Fix

Separa embed in un package dedicato E che A NON importa:

```go
// internal/embed/vobblybin/embed.go:
package vobblybin
import _ "embed"

//go:embed vobbly
var Binary []byte
```

B (vibbly) importa `internal/embed/vobblybin` direttamente o transitivamente.
A (vobbly) NON importa `vobblybin` (verify con `go list -deps ./cmd/A`).

## Makefile

```make
# Build A first, place binary in E's dir, then build B (which embeds A via E)
build-A-embed:
	go build -o internal/embed/vobblybin/A ./cmd/A

build: build-A-embed
	go build -o bin/B ./cmd/B
```

Aggiungi `internal/embed/vobblybin/A` a `.gitignore` (build artifact).

## Verify no cycle

```sh
go list -deps ./cmd/A | grep internal/embed/vobblybin
# Should be empty — A does NOT transitively depend on vobblybin
```

## Gotcha

- **Package name != path segment**: il package si chiama `vobblybin`, ma il path
  è `internal/embed/vobblybin/`. Il binary file dentro la dir può chiamarsi
  come vuoi (es. `vobbly`, non `vobblybin`).
- **`_ "embed"` blank import**: il package necessita almeno un import di `embed`
  per attivare la direttiva. Se il package ha altri import, non serve il blank.
- **Cross-compile**: go:embed cattura il binario al momento del build di B.
  Se vuoi multi-platform, devi replicare embed per ogni platform — ma valuta
  se ne hai bisogno (memory `no_cross_compile_host_provisioned_tools`).
- **Git tracking**: tipicamente binary artifact è gitignorato. Il commit di
  vibbly contiene lo SOURCE che genera il binary, non il binary stesso.

## Origine

Pattern estratto da fetta 2026-04-19 "vobbly local install per project"
(vibbly project). Il cycle fu scoperto a implement-time, non a plan-time —
anti-pattern tipico di "dependency analysis mancata in plan".
