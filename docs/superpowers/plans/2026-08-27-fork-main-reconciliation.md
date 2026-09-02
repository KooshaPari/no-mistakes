# Fork Main Reconciliation

Date: 2026-08-27

## Scope

This record explains the semantic reconciliation of `KooshaPari/main` at
`e0cf74db19a790e63ec33c4e3ab233d28b8f6c47` with the current candidate at
`dcb140615d722eef693851edf04127672e38677a`. The histories share merge base
`7abee5b138159077f2b654f7ad04df3ed1203810`.

No fork commit, ref, or artifact is deleted. The reconciliation merge retains
both ancestries. Tree selection is based on the dispositions below rather than
an unreviewed conflict resolution across the 82 paths reported by `git
merge-tree`.

## Machine classification

- Fork main has 54 non-merge commits after the common base.
- Candidate has 117 non-merge commits after the common base.
- Stable patch-ID comparison proves 34 fork commits are byte-equivalent to
  commits already present in the candidate ancestry.
- The remaining upstream-numbered commits are earlier forms of changes carried
  and subsequently evolved in the candidate ancestry. Release commits are not
  replayed.
- The genuine fork-authored surface is the policy/configuration group below,
  fork PRs #2 and #3, and the local snapshot `e0cf74d`.

## Fork-authored disposition

| Commit(s) | Surface | Disposition | Evidence |
|---|---|---|---|
| `2836dcf`, `e56ca39`, `8ea9384` | Stable CI gate names | Superseded | The candidate CI has the current platform matrix and gate structure; replaying the older whole-workflow replacement would discard later upstream CI behavior. |
| `a445c18` | Infisical workflow | Reject as unsafe | The workflow installs through an unpinned network script, exports secrets into a process environment, and uploads `.env` on failure. |
| `3ac4741` | Mergify rules | Reject as invalid/stale | It requests nonexistent `phenotype/core`, uses obsolete merge-message fields, assumes check names not provided by the current repository, and automatically closes stale work. |
| `ecb1479` | Generic Trunk configuration | Reject as mismatched | It enables Rust, Python, Node, Docker, and other linters for a Go repository without corresponding project configuration. |
| `1055270` | CircleCI matrix | Reject as mismatched | It installs and runs Rust, Python, and Node pipelines that are not part of this repository and contains permissive \|\| true gates. |
| `4b723aa` | OpenSSF Scorecard workflow | Replace separately if desired | The preserved version contains an invalid job-level `security` key and stale action pins; supply-chain scoring should be reintroduced as a current, independently validated change. |
| `ba63e17`, `59aa0b8` | Trunk/Prettier workflow | Superseded | Fork PR #2 already documents why the original Trunk action was broken. The replacement is a generic Prettier-only policy, not a Go correctness gate, and should not be required by this reconciliation. |
| `e26ace7` | Renovate template | Reject as mismatched | It contains Rust, Python, Node, and Docker package rules absent from this repository and grants broad dependency auto-merge behavior without current governance review. |
| `86cf6c1`, `e0cf74d` | Worktree GC, branch lint, and health-check steps | Reject incomplete implementation; preserve for redesign | `e0cf74d` adds step names and source but does not wire the steps into executable `steps.AllSteps`; it adds no behavioral tests; health-check assumes `origin/main` and warns on ordinary feature branches; worktree GC never performs GC; branch tracking parsing does not match Git's normal `[gone]` representation. The snapshot also commits a host-built 25.7 MB binary without reproducible provenance. |

## Candidate-owned follow-up

The candidate's cached repository-state work remains the intended product
change. Before merge it must retain:

1. local format, lint, race, build, and end-to-end validation;
2. fork-hosted CI on the exact reconciled head;
3. semantic review of any residual Sonar findings after the false 507-file
   new-code baseline is removed; and
4. human approval.

Potentially valuable rejected surfaces are not forbidden forever. Each may be
reintroduced through a narrow, current-architecture PR with executable tests or
the appropriate configuration validator and an explicit security review.
