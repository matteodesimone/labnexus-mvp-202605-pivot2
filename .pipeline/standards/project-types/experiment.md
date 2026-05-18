# Experiment — Project Type Definition

A lightweight, time-boxed exploration. No client, no formal process — just a question to answer.

When `PROJECT_TYPE` is `experiment`, the pipeline runs in a reduced mode optimized for speed over rigor.

---

## What It Is

A spike, proof of concept, or technical exploration. The goal is to answer a question ("can we do X?", "does Y work with Z?") within a fixed time box. The output is knowledge, not a product.

## What Changes in the Pipeline

| Step | Modification |
|---|---|
| `/v-spec` | Lightweight: just the question to answer, success criteria, and time box. No full FR list needed. |
| `/v-plan` | Minimal: approach and what to try first. No fette — experiments are single-pass. |
| `/v-review` | Optional. Experiments are disposable — review only if the code will be kept. |
| `/v-compound` | Important: the whole point is the learning. What did we discover? |

## What Doesn't Apply

- Gettoni and packages
- Discovery Goal
- Checkpoints
- Fette (single-pass by design)
- Client reports
- Full review rigor (unless the code survives)

## When to Use

- Technical spikes ("can Rails 8 do X?")
- Library evaluation ("is libY better than libZ for our use case?")
- Architecture exploration ("what happens if we use SQLite for this?")
- Learning exercises

## Time Box

Experiments should have a hard time limit. If the question isn't answered by the deadline, the answer is "not easily" — and that's a valid outcome. Document it in `/v-compound`.
