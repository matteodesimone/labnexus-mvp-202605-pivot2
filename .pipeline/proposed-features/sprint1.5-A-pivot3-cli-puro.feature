# language: it

Funzionalità: Sprint 1.5.A — Pivot 3 (eurouter default) + CLI puro multipiattaforma
  Come CTO che consegna il deliverable a Denis (Mac senza hardware Ollama),
  voglio che labnexus sia un binary CLI puro multipiattaforma con eurouter di default
  e niente bundle .app macOS,
  così che Denis possa usarlo senza configurazioni complicate e senza il bundle che
  introduceva verbosità inutile.

  # --- FR-1, FR-2: eliminato .app bundle, CLI puro cross-platform ---

  Scenario: build-mac.sh produce un binary standalone, niente bundle
    Dato che lancio "scripts/build-mac.sh"
    Allora viene prodotto un file "labnexus" (Mach-O arm64) come singolo binary
    E NON viene creato nessun "labnexus.app" o "Contents/" directory
    E NON sono prodotti file "Info.plist" o "labnexus-bin"
    E nello stdout/stderr di build NON compare la parola "osascript"

  # --- FR-3: .command launcher minimal per macOS ---

  Scenario: labnexus.command lancia il binary in Terminal
    Dato che esiste "labnexus" (binary) accanto a "labnexus.command"
    E faccio doppio click su "labnexus.command" dal Finder
    Allora Terminal si apre
    E il binary "labnexus" viene eseguito nella directory dello script
    E la TUI charmbracelet/huh parte regolarmente (TTY rilevato)

  Scenario: labnexus.command contiene solo 2 righe bash essenziali
    Dato che leggo "scripts/labnexus.command"
    Allora il file contiene esattamente: cd "$(dirname "$0")" + invocazione "./labnexus"
    E NON contiene wrapper osascript
    E NON contiene logica Info.plist o LSHandlerRank

  # --- FR-4, FR-5: rimozione gate eurouter + test associati ---

  Scenario: eurouter è ammesso senza LABNEXUS_ALLOW_CLOUD_PROVIDER setted
    Dato che l'env var "LABNEXUS_ALLOW_CLOUD_PROVIDER" non è impostata
    E un profilo dichiara provider "eurouter"
    Quando lancio "labnexus run --profile test-eurouter --input /tmp/in --output /tmp/out"
    Allora exit code è 0
    E NON viene restituito errore "ErrEurouterGateMissing"
    E la chiamata di rete va a "https://api.eurouter.ai/v1/chat/completions"

  Scenario: LABNEXUS_OLLAMA_ENDPOINT non-loopback è ammesso senza gate
    Dato che l'env var "LABNEXUS_OLLAMA_ENDPOINT" è impostata a "https://ollama-remoto.example.com/api"
    E un profilo dichiara provider "ollama"
    Quando lancio "labnexus run --profile test-ollama --input /tmp/in --output /tmp/out"
    Allora exit code è 0 (assumendo che l'endpoint risponda; o errore di rete normale)
    E NON viene restituito errore "non-loopback non ammesso"

  Scenario: i test gate eurouter e ollama-loopback non esistono più in select_test.go
    Dato che leggo "internal/provider/select_test.go"
    Allora NON contiene "TestSelect_EurouterRejectedWithoutApprovalGate"
    E NON contiene "TestSelect_OllamaRejectsRemoteEndpoint"
    E NON contiene riferimenti a "LABNEXUS_ALLOW_CLOUD_PROVIDER" o "requireLoopbackOrCloudGate"

  # --- FR-6: helpers_test.go non setta più il gate ---

  Scenario: features/helpers_test.go.buildEnv non aggiunge LABNEXUS_ALLOW_CLOUD_PROVIDER
    Dato che leggo "features/helpers_test.go::buildEnv"
    Allora NON aggiunge "LABNEXUS_ALLOW_CLOUD_PROVIDER" a keep
    E gli scenari BDD eurouter continuano a passare (senza bisogno del gate)

  # --- FR-7: 7 profili shippati con provider eurouter (transitorio 1.5.A) ---

  Scenario: i 7 profili dichiarano provider eurouter di default
    Dato che leggo i 7 file in "profili/*.yml" (revisione, rilievi, review-pack, audit-checklist, equipment-alert, competence-gap, pt-analysis)
    Allora ciascuno contiene la riga "provider: eurouter"
    E nessuno contiene "provider: ollama" come dichiarazione esplicita

  # --- FR-8: bug file aggiornati ---

  Scenario: i 2 bug file dei gate ex-Fetta-2 sono marcati intentional_deviation
    Dato che leggo ".pipeline/bugs/privacy-eurouter-gate-mancante.md"
    Allora il frontmatter ha "status: intentional_deviation_post_pivot_3"
    Dato che leggo ".pipeline/bugs/ollama-endpoint-env-override-leak.md"
    Allora il frontmatter ha "status: intentional_deviation_post_pivot_3"

  # --- FR-9: BDD scenari @manual su .app eliminati ---

  Scenario: gli scenari @manual su bundle/Info.plist/drag&drop sono rimossi
    Dato che leggo "features/motore-modo-interattivo.feature"
    Allora NON contiene "@manual" tag
    E NON contiene scenari su "Info.plist" o "CFBundleExecutable" o "LSHandlerRank"
    E NON contiene scenari su "drag&drop" da Finder iconcina

  # --- FR-10: README aggiornato post-pivot-3 ---

  Scenario: README zip-level italiano descrive il setup post-.app
    Dato che leggo "README.md"
    Allora contiene istruzioni Gatekeeper per binary nudo (non per .app)
    E contiene almeno una menzione di "labnexus.command"
    E NON menziona ".app" / "Contents/MacOS" / "Info.plist"
    E indica eurouter come provider di default
