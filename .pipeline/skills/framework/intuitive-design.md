# Skill: Intuitive Design

Based on Everett McKay's "Intuitive Design: Eight Steps to an Intuitive UI", adapted for all interface types we build: web GUI, CLI, game UI, API responses, and error messages.

Core definition: **if the target user must reason, memorize, experiment, seek help, or train to use the interface, it is not intuitive.**

---

## The 8 attributes across interface types

### 1. Discoverability

The user can find features and starting points without searching.

| Interface | Intuitive | Not intuitive |
|-----------|-----------|---------------|
| Web GUI | Navigation shows all primary actions. Empty states guide the user. | Features hidden in nested menus or behind unlabeled icons. |
| CLI | `--help` lists all commands with descriptions. Tab completion works. | Commands exist but aren't in `--help`. Hidden flags. |
| Game UI | Tutorial highlights interactive elements. UI hints on first encounter. | Player must discover mechanics by accident. |
| API/Error | Error response includes what went wrong and what to do. | Error code 500 with no body. `{"error": "invalid"}` with no context. |

**Review question**: can a new user find every primary action within 30 seconds?

### 2. Affordance

Elements communicate what they do through appearance or naming.

| Interface | Intuitive | Not intuitive |
|-----------|-----------|---------------|
| Web GUI | Buttons look clickable. Links are underlined or colored. Form fields have labels. | Flat text that's actually a button. Icons without labels. |
| CLI | Command names are verbs: `create`, `list`, `delete`. Flags describe their effect: `--verbose`, `--dry-run`. | Cryptic commands: `proc`, `xfr`. Single-letter flags without long form. |
| Game UI | Interactable objects glow or pulse. Cursor changes on hover. | No visual distinction between scenery and interactable objects. |
| API/Error | Field names match the domain: `user_email`, not `field_3`. Error type is descriptive: `validation_error`. | Generic field names. Error codes with no human meaning. |

**Review question**: for every interactive element, is its purpose obvious from its appearance/name alone?

### 3. Comprehensibility

The user understands what will happen before acting.

| Interface | Intuitive | Not intuitive |
|-----------|-----------|---------------|
| Web GUI | "Delete project and all its data" vs "Delete". Confirmation dialogs explain consequences. | "Submit" (submit what?). "Process" (what happens?). |
| CLI | `vibbly git sync` help says "Merges development into main. Pre-checks: tests pass, clean tree." | Command does multiple things with no description of what. |
| Game UI | Tooltip explains item effect before use. Skill tree shows prerequisites. | Item descriptions are vague. "Increases power." (How much? What power?) |
| API/Error | `"message": "Email is already registered. Use /auth/login or /auth/reset-password."` | `"message": "duplicate key"` |

**Review question**: before the user clicks/types/acts, do they know what will happen?

### 4. Responsive Feedback

The interface confirms what happened after every action.

| Interface | Intuitive | Not intuitive |
|-----------|-----------|---------------|
| Web GUI | "Project created successfully" toast. Loading spinner during wait. Progress bar for long operations. | Button click → nothing visible happens for 3 seconds. Silent success. |
| CLI | `✓ Project initialized in ./my-project`. Progress dots for long operations. Clear error with fix hint. | Command exits silently on success. Error goes to stderr with no prefix. |
| Game UI | Hit effect, sound, number popup on damage. Screen shake on impact. | Damage happens but there's no visual/audio confirmation. |
| API/Error | `201 Created` with the created resource in body. `422` with field-level errors. | `200 OK` for everything, including errors in the body. |

**Review question**: after every action, does the user know it worked (or why it didn't)?

### 5. Predictability

Same action → same result. No surprises. Consistent behavior.

| Interface | Intuitive | Not intuitive |
|-----------|-----------|---------------|
| Web GUI | Back button always goes back. Save always saves. Same layout on every page. | Back button sometimes resets the form. Different pages have different save behaviors. |
| CLI | Same flag means the same thing across all commands. Output format is consistent. | `--force` means "skip confirm" in one command and "overwrite" in another. |
| Game UI | Jump button always jumps. Attack always attacks the same way (until upgraded). | Controls change between levels without warning. |
| API/Error | Same endpoint, same input → same output. Error format is consistent across all endpoints. | Some endpoints return `{error: "..."}`, others return `{message: "..."}`. |

**Review question**: will the user be surprised by anything?

### 6. Efficiency

Minimum steps to accomplish the task. No unnecessary friction.

| Interface | Intuitive | Not intuitive |
|-----------|-----------|---------------|
| Web GUI | One-click for common actions. Smart defaults. Autofill where possible. | 5 clicks to do the most common task. Required fields that could be auto-populated. |
| CLI | Common workflows are one command. Sensible defaults so flags are optional. Piping works. | Need 3 commands for a common task. Every flag is required with no defaults. |
| Game UI | Inventory has shortcuts. Frequently used items are accessible. Quick save exists. | 4 menu levels to equip an item. No keyboard shortcuts. |
| API/Error | Batch endpoints for common multi-operations. Pagination with sensible defaults. | Need 10 API calls for what should be one. No batch support. |

**Review question**: can the most common task be done in fewer steps?

### 7. Forgiveness

The user can recover from mistakes. Destructive actions are guarded.

| Interface | Intuitive | Not intuitive |
|-----------|-----------|---------------|
| Web GUI | Undo for deletions. "Are you sure?" for destructive actions (with consequence stated). Autosave for drafts. | Delete is permanent with no confirmation. No undo. Lost work on navigation. |
| CLI | `--dry-run` to preview. Confirmation for destructive commands (default: No). Git as undo. | `rm -rf` with no confirmation. No dry-run option. No way to undo. |
| Game UI | Undo last action. Multiple save slots. Checkpoint system. | Permadeath without warning. One save slot that auto-overwrites. |
| API/Error | Soft delete (recoverable). Idempotent operations. Clear error on invalid input before processing. | Hard delete with no recovery. Non-idempotent PUT that corrupts on retry. |

**Review question**: if the user makes a mistake, can they recover without losing significant work?

### 8. Explorability

The user can try things safely. The interface is safe to explore.

| Interface | Intuitive | Not intuitive |
|-----------|-----------|---------------|
| Web GUI | Sandbox/preview mode. Tooltips explain features. Nothing breaks from clicking around. | Clicking a menu item triggers an irreversible action. No preview. |
| CLI | `--help` on every subcommand. `--dry-run` on destructive commands. Safe to run without args (shows help, not error). | Running without args does something unexpected. No help text. |
| Game UI | Tutorial area. Safe zones to practice. Controls displayed during gameplay. | Thrown into combat with no explanation. No safe area to learn. |
| API/Error | Sandbox/test environment. Detailed docs with examples. Clear 400 errors that explain what's wrong. | Only production. Docs are wrong or missing. Cryptic error codes. |

**Review question**: can the user explore the interface without fear of breaking something?

---

## Using this skill in the pipeline

**During `/spec`**: for every feature with a user interface (including CLI), write the **main instruction** — one sentence per screen/command that tells the user what to do. If it can't be one sentence, the screen does too much.

**During `/plan`**: design decisions about UI/CLI should reference these attributes. "The delete command uses `--dry-run` by default (forgiveness, explorability)" is a legitimate architecture decision.

**During `/review`**: evaluate every interface element against the 8 review questions. Any "no" is a finding. Severity: CRITICAL if discoverability or comprehensibility fails on a primary task, HIGH if forgiveness fails on destructive actions, MEDIUM for efficiency/predictability gaps, LOW for minor affordance issues.

**During `/compound`**: if users struggle with an interface (from feedback, testing, or observation), update this skill with the specific pattern that failed and how it was fixed.
