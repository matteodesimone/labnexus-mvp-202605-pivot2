# Skill: Idempotent BDD Given steps for composable setup

Pattern per Given step di setup in BDD (godog/Cucumber/behave) che bootstrappano prerequisiti mancanti invece di richiederli esplicitamente nella feature, garantendo che gli step siano riusabili sia standalone sia composti.

## Trigger

Più scenari richiedono lo stesso prerequisito composito (es. git repo + branch + commit baseline), ma la combinazione varia: alcuni scenari lo vogliono "clean", altri "dirty", altri "tagged", ecc.

## Anti-trigger

Un prerequisito unico globale per tutti gli scenari — in quel caso usa `BeforeAll` / `BeforeSuite` hook, non step individuali.

## Come applicarlo

1. Identifica il prerequisito transitorio (es. "git repo inizializzato").
2. Nel setup step, check presenza:

   ```go
   if _, err := os.Stat(filepath.Join(dir, ".git")); os.IsNotExist(err) {
       // bootstrap
   }
   ```

3. Bootstrap minimale: solo quello che serve per questo step.
4. Poi applica il comportamento specifico dello step (il SUO valore aggiunto, non quello del prerequisito).

## Anti-pattern

```go
// DON'T: fail hard se manca il prerequisito
func (tc *testContext) uncommittedChangesExist() error {
    // assume git is already initialized by prior step
    return os.WriteFile(...)
}
```

Se un altro scenario usa lo step standalone, fallisce enigmaticamente ("git: not a repo"). La feature file DEVE elencare il prerequisito mancante — overhead cognitivo per ogni scenario.

## Vantaggi

- Feature file dichiarativi: `Given uncommitted changes exist` — breve, auto-documenta lo scenario.
- Scenari indipendenti: zero ordering dependencies implicite.
- Refactor locale: cambiare il bootstrap tocca UN step, non N scenari.

## Gotcha

- **Idempotenza reale**: non solo check-and-skip, ma anche: quale stato assume "già setuppato"? Se uno scenario precedente ha inizializzato git con branch custom, lo step rispetta quello o lo sovrascrive? Preferenza standard: rispettalo — lo step fa SOLO quello che manca.
- **Test data leakage**: se uno scenario lascia file residui, scenari successivi possono ereditare. Soluzione ortogonale: `tmpDir` fresh per scenario (godog Before hook).
- **Nome dello step**: non mentire. `uncommittedChangesExist` è onesto. `aVibblyProjectWithHooksAndGit` sarebbe meglio di `aVibblyProjectWithHooks` se lo step imbarca git dentro.

## Origine

Pattern estratto dal bugfix `bdd-tests-vibbly-not-in-path` (vibbly project, 2026-04-19). Durante il fix di 2 scenari pre-esistenti (hook + vobbly_release) che non setuppavano git, si è scelto di rendere gli step idempotenti invece di modificare le feature file.
