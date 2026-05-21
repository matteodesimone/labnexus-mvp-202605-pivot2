# language: it

Funzionalità: Sprint 1.5.A — Pivot 3 (eurouter default) + CLI puro multipiattaforma
  Come CTO che consegna il deliverable a Denis (Mac senza hardware Ollama),
  voglio che labnexus sia un binary CLI puro multipiattaforma con eurouter di default
  e niente bundle .app macOS, così che Denis possa usarlo senza configurazioni complicate
  e senza il bundle che introduceva verbosità inutile.

  Nota di approccio (test-scaffold): le verifiche sono STATIC ANALYSIS del codebase
  + filesystem checks. Non lanciamo bash scripts/build-mac.sh in CI (richiederebbe
  Go toolchain cross-compile environment + side-effects). Quel test si fa manualmente
  al pre-flight di deploy (pattern Sprint 1+2). Qui verifichiamo che il SORGENTE
  degli script + del codice Go riflette le rimozioni post-pivot-3.

  # --- FR-1, FR-2: .app bundle eliminato da scripts/build-mac.sh ---

  Scenario: scripts/build-mac.sh è semplificato a un singolo go build standalone
    Dato che il file di sorgente "scripts/build-mac.sh" esiste
    Allora il file di sorgente "scripts/build-mac.sh" contiene la stringa "GOOS=darwin GOARCH=arm64"
    E il file di sorgente "scripts/build-mac.sh" non contiene la stringa "labnexus.app"
    E il file di sorgente "scripts/build-mac.sh" non contiene la stringa "Info.plist"
    E il file di sorgente "scripts/build-mac.sh" non contiene la stringa "osascript"
    E il file di sorgente "scripts/build-mac.sh" non contiene la stringa "CFBundleExecutable"
    E il file di sorgente "scripts/build-mac.sh" non contiene la stringa "LSHandlerRank"
    E il file di sorgente "scripts/build-mac.sh" non contiene la stringa "labnexus-bin"

  Scenario: scripts/build-zip.sh non include più il bundle .app
    Dato che il file di sorgente "scripts/build-zip.sh" esiste
    Allora il file di sorgente "scripts/build-zip.sh" non contiene la stringa "labnexus.app"
    E il file di sorgente "scripts/build-zip.sh" non contiene la stringa "Contents/MacOS"
    E il file di sorgente "scripts/build-zip.sh" contiene la stringa "labnexus.command"

  # --- FR-3: labnexus.command minimal launcher ---

  Scenario: scripts/labnexus.command esiste come launcher minimal
    Dato che il file di sorgente "scripts/labnexus.command" esiste
    Allora il file di sorgente "scripts/labnexus.command" contiene la stringa "cd"
    E il file di sorgente "scripts/labnexus.command" contiene la stringa "./labnexus"
    E il file di sorgente "scripts/labnexus.command" non contiene la stringa "osascript"
    E il file di sorgente "scripts/labnexus.command" non contiene la stringa "Info.plist"
    E il file di sorgente "scripts/labnexus.command" ha il bit eseguibile attivo

  # --- FR-4: rimosso ErrEurouterGateMissing + requireLoopbackOrCloudGate da provider.go ---

  Scenario: internal/provider/provider.go non contiene più i gate cloud Fetta 2
    Dato che il file di sorgente "internal/provider/provider.go" esiste
    Allora il file di sorgente "internal/provider/provider.go" non contiene la stringa "ErrEurouterGateMissing"
    E il file di sorgente "internal/provider/provider.go" non contiene la stringa "requireLoopbackOrCloudGate"
    E il file di sorgente "internal/provider/provider.go" non contiene la stringa "LABNEXUS_ALLOW_CLOUD_PROVIDER"

  # --- FR-5: rimossi 10 test gate in select_test.go ---

  Scenario: internal/provider/select_test.go non contiene più i test gate
    Dato che il file di sorgente "internal/provider/select_test.go" esiste
    Allora il file di sorgente "internal/provider/select_test.go" non contiene la stringa "TestSelect_EurouterRejectedWithoutApprovalGate"
    E il file di sorgente "internal/provider/select_test.go" non contiene la stringa "TestSelect_EurouterRejectedViaEnvWithoutApprovalGate"
    E il file di sorgente "internal/provider/select_test.go" non contiene la stringa "TestSelect_OllamaRejectsRemoteEndpoint"
    E il file di sorgente "internal/provider/select_test.go" non contiene la stringa "TestSelect_OllamaRemoteAcceptedWithCloudGate"
    E il file di sorgente "internal/provider/select_test.go" non contiene la stringa "allowEurouterEnv"

  # --- FR-6: helpers_test.go non setta più LABNEXUS_ALLOW_CLOUD_PROVIDER ---

  Scenario: features/helpers_test.go non setta più il gate cloud
    Dato che il file di sorgente "features/helpers_test.go" esiste
    Allora il file di sorgente "features/helpers_test.go" non contiene la stringa "LABNEXUS_ALLOW_CLOUD_PROVIDER"

  # --- FR-7: 7 profili shippati con provider eurouter (ereditato dal master 1.5.B+) ---
  # Nota Sprint 1.5.B: i 7 profili non dichiarano più `provider` esplicito (vedi
  # spec FR-13) — è ereditato dal master `labnexus.config.toml`. L'invariante
  # rilevante è: nessun profilo deve dichiarare provider = "ollama" (override
  # locale a ollama sarebbe regressione del pivot 3).

  Scenario: nessun profilo shippato dichiara override locale a ollama
    Allora nessun profilo shippato dichiara la riga "provider = \"ollama\""
    E nessun profilo shippato dichiara la riga "provider: ollama"

  # --- FR-7b: audit-trail truthfulness (review finding HIGH-2, 2026-05-21) ---
  # Il frontmatter_default.locale era hardcoded a "true" dai 7 profili Sprint 1
  # quando il provider era ollama (locale). Post-pivot-3, provider default è
  # eurouter (cloud) e quindi "locale: true" diventerebbe audit trail menzognero
  # sull'output prodotto. Rimosso dai 7 profili in sub-fetta 1.5.A (review loop 2).

  Scenario: nessun profilo shippato dichiara locale: true nel frontmatter
    Allora nessun profilo shippato dichiara la riga "locale: true"
    E nessun profilo shippato dichiara la riga "locale: false"

  # --- FR-8: 2 bug file aggiornati a intentional_deviation ---

  Scenario: privacy-eurouter-gate-mancante è marcato intentional_deviation
    Dato che il file di sorgente ".pipeline/bugs/privacy-eurouter-gate-mancante.md" esiste
    Allora il file di sorgente ".pipeline/bugs/privacy-eurouter-gate-mancante.md" contiene la stringa "intentional_deviation_post_pivot_3"

  Scenario: ollama-endpoint-env-override-leak è marcato intentional_deviation
    Dato che il file di sorgente ".pipeline/bugs/ollama-endpoint-env-override-leak.md" esiste
    Allora il file di sorgente ".pipeline/bugs/ollama-endpoint-env-override-leak.md" contiene la stringa "intentional_deviation_post_pivot_3"

  # --- FR-9: BDD scenari @manual su bundle eliminati ---

  Scenario: features/motore-modo-interattivo.feature non contiene più scenari bundle .app
    Dato che il file di sorgente "features/motore-modo-interattivo.feature" esiste
    Allora il file di sorgente "features/motore-modo-interattivo.feature" non contiene la stringa "Info.plist"
    E il file di sorgente "features/motore-modo-interattivo.feature" non contiene la stringa "CFBundleExecutable"
    E il file di sorgente "features/motore-modo-interattivo.feature" non contiene la stringa "LSHandlerRank"
    E il file di sorgente "features/motore-modo-interattivo.feature" non contiene la stringa "drag&drop"
    E il file di sorgente "features/motore-modo-interattivo.feature" non contiene la stringa "labnexus.app"
    E il file di sorgente "features/motore-modo-interattivo.feature" non contiene la stringa "Contents/MacOS"
    E il file di sorgente "features/motore-modo-interattivo.feature" non contiene la stringa "wrapper bash"

  # --- FR-10: README aggiornato post-pivot-3 ---

  Scenario: README.md non menziona più .app, indica eurouter default e labnexus.command
    Dato che il file di sorgente "README.md" esiste
    Allora il file di sorgente "README.md" non contiene la stringa "labnexus.app"
    E il file di sorgente "README.md" non contiene la stringa "Contents/MacOS"
    E il file di sorgente "README.md" non contiene la stringa "Info.plist"
    E il file di sorgente "README.md" contiene la stringa "labnexus.command"
    E il file di sorgente "README.md" contiene la stringa "eurouter"
