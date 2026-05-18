# Skill: BDD tag-based scenario scoping (godog/Cucumber)

Pattern per applicare setup/teardown condizionale in Before/After hooks di BDD framework (godog, Cucumber-Ruby, behave, ecc.) usando Gherkin tags invece di ispezionare URI o nome scenario.

## Trigger

Uno o più scenari richiedono setup particolare (tmp dir, env var, fixtures, mocks) mentre altri **non** lo richiedono. Un Before hook deve discriminare.

## Anti-trigger

Il setup è universale per tutti gli scenari di un feature file o di tutto il test run. In quel caso usa `Before(all)` / `BeforeFeature` direttamente — niente tag.

## Come applicarlo

1. Aggiungi tag a livello feature (o scenario) nei file `.feature`:

   ```gherkin
   @needs-install-dir @slow
   Feature: Installation flow
     Scenario: Installs binary to custom dir
       ...
   ```

2. Nel Before hook, controlla `sc.Tags`:

   ```go
   ctx.Before(func(c context.Context, sc *godog.Scenario) (context.Context, error) {
       hasTag := func(name string) bool {
           for _, t := range sc.Tags {
               if t.Name == name {
                   return true
               }
           }
           return false
       }
       if hasTag("@needs-install-dir") {
           // setup specifico
       }
       return c, nil
   })
   ```

3. Per tag multipli cumulativi, itera una volta e imposta flag distinti.

## Anti-pattern

```go
// DON'T: string-match su Uri o scenario name
if strings.Contains(sc.Uri, "install") { ... }
```

Silenziosamente si rompe al primo rename del file o dello scenario. Il tag è un contratto esplicito parte della feature; il nome file è un dettaglio organizzativo.

## Gotcha

- In godog, `sc.Tags` include sia tag a livello feature che a livello scenario (automaticamente uniti). Non serve append manuale.
- Tag iniziano con `@`. Il confronto è esatto, case-sensitive.
- Se rimuovi un tag dalla feature, rimuovi anche i check corrispondenti nel Before hook. Tag orfani nel codice non rompono nulla ma confondono il prossimo lettore.
- Tag come documentazione: se un tag descrive un requisito di ambiente (`@needs-docker`, `@needs-network`, `@needs-install-dir`), il nome serve come auto-documentazione — self-explanatory senza cercare nell'hook.

## Origine

Pattern estratto da fetta 1 "vibbly dev workflow loop" (vibbly project, 2026-04-19). In quella fetta lo string-matching su `sc.Uri` è stato sostituito con tag `@needs-install-dir` dopo aver identificato la brittleness al review.
