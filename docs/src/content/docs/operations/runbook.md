---
title: Runbook
---

# Runbook

## Pipeline stuck

1. `no-mistakes pipeline status`
2. `no-mistakes pipeline cancel <id>`
3. If that fails: `no-mistakes worktree prune --all`

## Daemon unresponsive

1. `no-mistakes daemon ping`
2. `no-mistakes daemon restart`
3. If restart fails: check `/var/log/no-mistakes/daemon.log`
