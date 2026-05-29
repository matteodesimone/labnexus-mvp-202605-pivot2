---
date: 2026-05-29
status: planned
source: compound-learning bugfix KB v2.1 rimappatura profili
priority: medium
tags: profili, kb, config, DRY, master-inheritance, sprint-2
---

# `core_kb_files` ereditato dal master config

## Idea

Far ereditare ai profili un blocco **core di `kb_files`** dal master `labnexus.config.toml`, come già accade per `provider`/`modello`/`temperature` (vedi skill `master-profile-config-inheritance`). Il profilo dichiarerebbe solo i `kb_files` *specifici* della capability; il core (i 9 file `00_*`/`01_SYSTEM` + `02_LOGICA_ISPETTIVA/01_COME_PENSA`) verrebbe iniettato dal merge.

## Razionale

Dopo la rimappatura sul KB v2.1 (bugfix 2026-05-29, mappatura "ampio"), il blocco core di ~10 file è **duplicato verbatim nei 7 profili**. Conseguenze:
- ogni futuro cambiamento del core (es. nuovo file `01_SYSTEM/07_*`) richiede 7 edit manuali → drift garantito;
- i profili sono lunghi e il segnale specifico-per-capability si perde nel rumore del core ripetuto.

## Schizzo

- Master: `core_kb_files = [...]` (lista del blocco sempre-in-contesto).
- Profilo: `kb_files = [...]` (solo i requisiti/appendici specifici).
- `config.Merge`: se il profilo non disabilita esplicitamente, prepende `core_kb_files` a `kb_files` (dedup). Flag opt-out per profili che vogliono controllo totale.
- `profile.Validate` resta invariato (valida la lista mergiata).

## Note

- Non urgente: la duplicazione attuale funziona ed è coperta dal test shipped-vs-shipped.
- Sinergico con un eventuale refinement "mirato" dei profili pesanti (revisione/audit-checklist).
- Vedi `.pipeline/solutions/2026-05-29-bugfix-kb-v2.1-rimappatura-profili.md` (Tech Debt).
