# Fork Status — KooshaPari/no-mistakes

**Last updated:** 2026-09-03 (G3 forensic pass)
**Authoritative:** see upstream [kunchenguid/no-mistakes](https://github.com/kunchenguid/no-mistakes) for the base project; this repo is the actively-diverged living gate.

## TL;DR

This repo is a **true, deliberately-diverged fork** of `kunchenguid/no-mistakes` (a Go git-review gate / no-mistakes pipeline). It is the **living no-mistakes implementation** — it carries a substantial local evaluation and release toolkit that does not exist upstream. It is not a passive mirror and not a delete candidate.

## Repo identity

| Field | Value |
|---|---|
| Kind | Fork (true — parent `kunchenguid/no-mistakes` alive) |
| Languages | Go (≈7.1 MB of Go source + `.no-mistakes.yaml` governance) |
| Default branch | `main` |
| Remotes | `origin` (our GH fork) · `upstream` (kunchenguid/no-mistakes) · `no-mistakes` (gate's local mirror under `~/.no-mistakes/repos/`) |

## Divergence from upstream (powered by local gitAccess)

| Metric | Count |
|---|---|
| **Ours ahead of upstream** | **76 commits** |
| **Upstream ahead of us** (not yet merged) | **137 commits** |
| Local branch state | `main` clean, tracking `origin/main` |

### What the 76 local commits add (the valuable divergence)

- **Local evaluation/review toolkit** — `feat(eval): add local review evaluation toolkit`, score replay findings against human gold, auto-collect a local corpus, pin gold-only holdout, ingest confirmed post-PR misses as false-negative gold.
- **Release automation** — release-please driven, **1.45 → 1.54** (release tooling, semantic-main reconciliation, generated-state preservation).
- **CI / security hardening** — restore security + protection gates, split Windows test legs, bind CI platform matrix to runner selection, require complete runner expressions, workflow-for-release-please (avoid protected-check deadlock).
- **SonarCloud hygiene** — extract duplicated string-literal constants (internal/config, internal/git) to clear SonarCloud duplicate-literal findings; tests pass, lint clean.
- **Daemon/eviting fixes** — reap lingering run-worktree processes, publish evidence to an orphan branch, build-identity on run records, move test evidence out of system temp.
- **Docs** — Operations + Demo quadrants, on-device demo, vision statement update, AI-slop badges.

## Syncing posture

- `upstream/main` is tracked at `refs/remotes/upstream/main`.
- 137 upstream commits are not yet merged into our `main`. Merging them is **fork-sync hygiene** to consider separately (rebases/merges may conflict with the local eval toolkit).
- Do **not** blindly force-push `main`; keep the local eval/release divergence until a deliberate reconciliation.

## Provenance

Fork created from `kunchenguid/no-mistakes`. The divergence is intentional and governed by the no-mistakes gate itself (this repo runs `no-mistakes` review on its own PRs). See `README.md` and the gate config under `.no-mistakes.yaml`.

## Contact

For questions, open an issue or contact via the KooshaPari org.