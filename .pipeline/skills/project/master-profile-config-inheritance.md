# Skill: Master-Profile Config Inheritance

**Created**: 2026-05-21 in sub-fetta 1.5.B (compound)
**Trigger**: progetti CLI tool con ≥3 file di configurazione che condividono molti campi comuni e differiscono solo per pochi campi specifici (es. profili LabNexus, capability AICertus, ma anche project templates, build configs, ecc.).

## Il problema

Sprint 1 aveva 7 profili `.yml`, ognuno duplicava 5+ campi comuni (`provider`, `modello`, `temperature`, `max_tokens`, `context_window`). Per cambiare il default globale (es. switch provider ollama → eurouter), Denis doveva editare 7 file. Errore-prone, friction.

## Il pattern

Un **config master** al root del progetto fornisce i defaults globali. I profili dichiarano solo i campi specifici (es. `kb_files`, `trigger_prompt`) + opzionalmente override puntuali. Il runtime fa `Merge(master, profile) → profile merged` prima di Validate/usage.

### Schema (Go example)

```go
// internal/config/config.go
package config

type Master struct {
    Provider       string  `toml:"provider"`
    Modello        string  `toml:"modello"`
    Temperature    float64 `toml:"temperature"`
    MaxTokens      int     `toml:"max_tokens"`
    // Side-channel fields (es. API keys, endpoints) sono nel master ma propagati
    // via env vars al runtime (vedi applyMasterToEnv).
    EurouterAPIKey string  `toml:"eurouter_api_key"`
    OllamaEndpoint string  `toml:"ollama_endpoint"`
}

func Load(path string) (*Master, error) {
    data, err := os.ReadFile(path)
    if err != nil { return nil, err }
    var m Master
    if err := toml.Unmarshal(data, &m); err != nil { return nil, err }
    return &m, nil
}

func Merge(master *Master, p *profile.Profile) *profile.Profile {
    if master == nil || p == nil { return p }
    cp := *p
    if cp.Provider == "" { cp.Provider = master.Provider }
    if cp.Modello == "" { cp.Modello = master.Modello }
    if cp.Temperature == 0 { cp.Temperature = master.Temperature }
    if cp.MaxTokens == 0 { cp.MaxTokens = master.MaxTokens }
    return &cp
}
```

### Convenzioni critiche

1. **Precedenza**: profile (se setted) > master. Mai il contrario. L'utente che customizza un profile si aspetta che il suo valore vinca.
2. **Validate post-merge**: chiama `Validate` SUL PROFILO MERGED, non sull'isolato. Un profile parziale (provider non setted) deve passare validation se il master lo fornisce.
3. **Fallback CWD per portabilità**: il binary risolve il config relativo a sé stesso (paths.DeliveryRoot), ma fa fallback a `./labnexus.config.toml` nella CWD. Permette al binary di essere spostato senza perdere il config.
4. **Side-channel via env vars**: campi del master che vanno a librerie/HTTP client esterni (API key, endpoint) vengono propagati via `os.Setenv` con precedenza env > master. Funzione dedicata `applyMasterToEnv(master)`.
5. **Documentazione esplicita di precedenza**: nel commento del Merge + nei commenti del config TOML.

## Bug subtle: zero-value detection

Per i campi numerici (`float64`, `int`), il check `== 0` non distingue "campo omesso nel TOML" da "campo esplicitamente setted a 0". Caso problematico: Denis vuole `temperature = 0.0` per output deterministico, ma viene override col master.

Fix vero: pointer types (`*float64`, `*int`) — refactor invasivo.

Workaround pragmatico: documentato nel commento; utente che vuole zero usa un valore piccolo (`0.001`).

Vedi `.pipeline/ideas/2026-05-21-config-merge-pointer-types-deterministic.md` per il backlog.

## Anti-pattern

- **NON validare prima del merge**: profile parziale fallirebbe; merge poi non sa cosa fare.
- **NON propagare API keys via campi profile o struct passing**: usare env vars per le credenziali (best practice security).
- **NON hardcodare il path del master**: usare flag `--config` + fallback CWD per portability.

## Precedente concreto

Sub-fetta 1.5.B LabNexus — pattern emerso da pivot 3 (default eurouter): Denis edita `labnexus.config.toml` con `eurouter_api_key`, i 7 profili shippati ereditano provider/modello/parametri. Cambio di default globale = 1 riga di edit.

Vedi `.pipeline/solutions/2026-05-21-sub-fetta-1.5.B-config-toml-parser-nuovi.md` e `internal/config/config.go`.

## Quando NON applicare

- Progetti con 1-2 file di configurazione: l'overhead Merge non si ripaga
- Progetti dove ogni profilo è radicalmente diverso (no campi comuni): non c'è inheritance utile
- Progetti dove la duplicazione è semanticamente significativa (es. ogni profilo è una versione di reference da non collegare): l'inheritance maschera la separazione voluta

## Generalizzabilità (verso framework promotion)

Pattern universal (vale per CLI tools con multi-config, build systems con per-env override, ecc.). Candidato a `.pipeline/proposed-updates/` come `skills/framework/`.
