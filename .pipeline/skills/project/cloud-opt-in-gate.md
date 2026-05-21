---
slug: cloud-opt-in-gate
created: 2026-05-21
source_cycle: bugfix pass post-Fetta-2 (#001 + #003)
applies_to: ogni nuovo provider o endpoint override che possa puntare fuori dalla macchina locale di Denis
---

# Cloud opt-in gate via singolo env var

## Problema

LabNexus post-pivot 2 ha promessa contrattuale "**tutta l'elaborazione AI gira locale**" sui dati reali del SGQ di Denis. Esistono però vari punti dove un override (env, flag, profilo) potrebbe silenziosamente deviare il traffico a un endpoint cloud:

- `--provider eurouter` o `LABNEXUS_PROVIDER=eurouter` → traffico a `api.eurouter.ai`
- `LABNEXUS_OLLAMA_ENDPOINT=https://ollama-cloud.example.com/api` → traffico remoto camuffato da "provider: ollama"

Senza gate, una distrazione di configurazione (env stantia, profilo errato in dev) può esfiltrare dati.

## Pattern

**Singolo env var di consenso esplicito**: `LABNEXUS_ALLOW_CLOUD_PROVIDER`. Senza valore impostato:
- `eurouter` rifiutato con `ErrEurouterGateMissing` (qualunque sia la provenienza: flag/env/profile)
- `LABNEXUS_OLLAMA_ENDPOINT` non-loopback rifiutato con `requireLoopbackOrCloudGate`

Con il gate impostato (qualunque valore non vuoto, raccomandato `approved-for-synthetic-data`):
- Entrambi gli override sono accettati.

Loopback hosts ammessi senza gate per `LABNEXUS_OLLAMA_ENDPOINT`:
- `127.0.0.0/8` (IPv4 loopback)
- `::1` (IPv6 loopback)
- `localhost` (alias)

## Implementazione canonica

`internal/provider/provider.go`:
- `case "eurouter"`: check `os.Getenv("LABNEXUS_ALLOW_CLOUD_PROVIDER") == ""` → return `ErrEurouterGateMissing`.
- `case "ollama"`: chiama `requireLoopbackOrCloudGate(endpoint)`.
- `requireLoopbackOrCloudGate`: parse URL, estrae host, se loopback → ok; altrimenti se gate set → ok; altrimenti errore esplicito che cita NFR-1 e i passi per autorizzare.

## Test infrastructure (BDD)

`features/helpers_test.go::buildEnv` setta `LABNEXUS_ALLOW_CLOUD_PROVIDER=approved-for-synthetic-data` **sempre** per i test BDD (ambiente test controllato; httptest server locale è effettivamente loopback). Il gate stesso è coperto da unit test in `internal/provider/select_test.go`.

## Quando applicarlo

Ogni volta che si aggiunge:
- Un nuovo provider con backend remoto
- Un override env/flag che può ridirezionare il traffico provider esistente
- Un nuovo endpoint override (es. `LABNEXUS_EUROUTER_ENDPOINT` — già coperto implicitamente perché eurouter passa già dal gate)

## Cosa NON copre

- Validation di HTTPS vs HTTP (sito remoto loopback OK; ma sito remoto HTTPS fuori loopback richiede gate — questo è coperto).
- Detection di tunnel SSH locali che proxy a host remoti. Out-of-scope per Sprint 1; rientra in threat model "system administrator hostile".
- Audit logging dei tentativi di override (potenziale Sprint 2 enhancement).
