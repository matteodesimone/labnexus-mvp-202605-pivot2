# .pipeline/config.sh
# Configuration for the agent pipeline system
# Edit this file to customize behavior
#
# ╔══════════════════════════════════════════════════════════╗
# ║  ESSENTIAL — configure these for your project           ║
# ║  Everything below ADVANCED has sensible defaults.       ║
# ╚══════════════════════════════════════════════════════════╝

# --- Project type ---
# Options: prove-out, internal-tool, experiment
# See .pipeline/standards/project-types/<type>.md for type-specific rules.
PROJECT_TYPE="prove-out"

# --- Deploy / Ship ---
# SHIP_CMD is set per-domain in framework/domains/<domain>/config-domain.sh
# (cli-tool: goreleaser; web-app: fly deploy / equivalent; agent: package
# claude-code-plugin; game: itch.io upload / Steam build). The base ships
# no default because every domain's release pipeline differs — picking one
# here would re-introduce the web-app bias the refactor removed.

# --- External LLM reviewers ---
# Enable the ones you have installed. Claude always participates.
REVIEWER_GEMINI_ENABLED=false
REVIEWER_MISTRAL_ENABLED=true
REVIEWER_OLLAMA_ENABLED=true
REVIEWER_OLLAMA_MODEL="qwen3-coder:30b"
REVIEWER_CODEX_ENABLED=true

# --- Review loop ---
# REVIEW: mode and severity threshold in one value
#   "manual"    = single pass, human fixes and re-runs /review (default)
#   "critical"  = until-clean loop, auto-fix CRITICAL only (still escalates to human)
#   "high"      = until-clean loop, auto-fix HIGH+
#   "medium"    = until-clean loop, auto-fix MEDIUM+
#   "all"       = until-clean loop, fix everything
#
# REVIEW_MAX_LOOPS: max iterations for non-manual modes (1-10)
#   Each loop: review → write test for finding → fix → run all tests → re-review.
#   Ignored when REVIEW="manual".
#   CRITICAL always escalates to human regardless of REVIEW value.
#
REVIEW="all"
REVIEW_MAX_LOOPS=3

# --- Budget (prove-out only, ignored for other types) ---
# See .pipeline/standards/project-types/prove-out.md for full budget model.
PACKAGE_SIZE="S"              # S (20 gettoni, 3 weeks) | M (35 gettoni, 5 weeks)
BUDGET_TOTAL=20               # Derived from package: S=20, M=35
BUDGET_INFRA=5                # ~25% of total — infra, security, setup (internal, client doesn't see)
BUDGET_CLIENT=15              # Remainder — client-visible features
SPRINT_DURATION_WEEKS=3       # S=3, M=5

# ╔══════════════════════════════════════════════════════════╗
# ║  ADVANCED — defaults are fine for most projects         ║
# ╚══════════════════════════════════════════════════════════╝

# --- Reviewer commands (override if installed elsewhere) ---
REVIEWER_GEMINI_CMD="gemini"
REVIEWER_MISTRAL_CMD="vibe"
REVIEWER_OLLAMA_CMD="ollama run"
REVIEWER_CODEX_CMD="codex"
REVIEWER_TIMEOUT=120

# --- Review checks ---
MULTI_REVIEW_QUALITY=true
MULTI_REVIEW_TEST_COVERAGE=true
MULTI_REVIEW_SECURITY=true
MULTI_REVIEW_PRIVACY=true
MULTI_REVIEW_DATA_INTEGRITY=false
MULTI_REVIEW_DEPLOY_VERIFICATION=false

# --- Quality thresholds ---
MAX_FUNCTION_LINES=20
MAX_CYCLOMATIC_COMPLEXITY=5
MAX_NESTING_DEPTH=2

# --- Gate modes ---
GATE_SPEC="human"
GATE_PLAN="human"
GATE_TEST_SCAFFOLD="human"
GATE_REVIEW_FAIL="human"
GATE_DEPLOY="human"
GATE_COMPOUND="human"

# --- Per-gate reviewer override ---
GATE_SPEC_REVIEWERS=""
GATE_PLAN_REVIEWERS=""
GATE_TEST_SCAFFOLD_REVIEWERS=""
GATE_REVIEW_FAIL_REVIEWERS=""
GATE_DEPLOY_REVIEWERS=""
GATE_COMPOUND_REVIEWERS=""

# --- Per-check reviewer override ---
REVIEW_QUALITY_REVIEWERS=""
REVIEW_TEST_COVERAGE_REVIEWERS=""
REVIEW_SECURITY_REVIEWERS=""
REVIEW_PRIVACY_REVIEWERS=""
REVIEW_DATA_INTEGRITY_REVIEWERS=""
REVIEW_DEPLOY_VERIFICATION_REVIEWERS=""

# --- Reviewer execution ---
REVIEWER_PARALLEL=true

# --- Worktree ---
WORKTREE_ENABLED=false

# --- Skip permissions ---
SKIP_PERMISSIONS=false

# config-domain.sh — Software Development Base
# Universal foundation for every software-development domain.
# Domain-specific overlays (cli-tool, agent, web-app, game) extend this with
# their own SHIP_CMD and any extra checks.
# Core settings (reviewers, gates, review loop) live in config-base.sh.

# --- Domain checks (universal across software-dev) ---
# Core checks (always available)
CHECK_QUALITY=true
CHECK_TEST_COVERAGE=true
CHECK_SECURITY=true
CHECK_PRIVACY=true
# Extended checks (enable per project)
CHECK_DATA_INTEGRITY=false
CHECK_DEPLOY_VERIFICATION=false

# --- Per-check reviewer override ---
CHECK_QUALITY_REVIEWERS=""
CHECK_TEST_COVERAGE_REVIEWERS=""
CHECK_SECURITY_REVIEWERS=""
CHECK_PRIVACY_REVIEWERS=""
CHECK_DATA_INTEGRITY_REVIEWERS=""
CHECK_DEPLOY_VERIFICATION_REVIEWERS=""

# --- Quality thresholds (universal coding rules) ---
MAX_FUNCTION_LINES=20
MAX_CYCLOMATIC_COMPLEXITY=5
MAX_NESTING_DEPTH=2

# --- Project documentation ---
# Which docs to generate at init and keep aligned with code.
DOCS_USER_GUIDE=true       # docs/USER-GUIDE.md
DOCS_DEV_GUIDE=true        # docs/DEV-GUIDE.md
DOCS_LLM_CONTEXT=true      # CONTEXT.md

# Note: SHIP_CMD is defined per-domain in cli-tool/, web-app/, agent/, game/
# config-domain.sh — coding-base intentionally has no SHIP_CMD default.

# config-domain.sh — cli-tool Domain
# Sourced AFTER coding-base/config-domain.sh during vibbly init.
# Inherits all CHECK_*, MAX_*, DOCS_* settings; adds cli-tool-specific ship command.

# --- Ship command for cli-tool ---
# Sprint 1: NO goreleaser. Binary is NOT signed (Apple Developer cert deferred to Sprint 2).
# Deliverable is a manual zip with labnexus.app + profili/ + KB-ispettore/ + README.
# The exact build+zip script will be created during /v-plan; for now a stub keeps the
# config valid and ensures /v-deploy fails fast if invoked prematurely.
SHIP_CMD="./scripts/build-zip.sh"
