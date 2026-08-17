// Package airlock is the airlock-v2 autocommit daemon wrapper for no-mistakes.
//
// This file is the no-mistakes absorb of the bash wrapper at
// airlock/integrations/phenovcs-airlock-v2/packaging/airlock-v2-autocommit.sh
// (in phenovcs-airlock-v2). The wrapper just shells out to `airlock-v2
// autocommit` — the real daemon logic is in the compiled binary, porting
// of which is out-of-scope here.
//
// LIFECYCLE MGMT: this code is parked-by-default. You manage it yourself.
// To activate, wire the lifecycle package to call Start()/Stop() in your
// scenario. Do NOT daemonize via launchd without reading the lifecycle
// package doc first.
package airlock

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// AutocommitOptions configures a single autocommit run.
type AutocommitOptions struct {
	// WorkspaceRoot is the path to the meta-workspace to operate on.
	WorkspaceRoot string
	// BeadsLedgerPath is the path to the Beads ledger (state DB).
	BeadsLedgerPath string
	// MaxCommitsPerRun caps the number of commits per invocation (0 = unlimited).
	MaxCommitsPerRun int
	// DryRun reports what would happen without committing.
	DryRun bool
	// Timeout bounds a single invocation.
	Timeout time.Duration
}

// DefaultOptions returns safe defaults for a single invocation.
func DefaultOptions() AutocommitOptions {
	return AutocommitOptions{
		WorkspaceRoot:     "/Users/kooshapari/CodeProjects/Phenotype/repos",
		BeadsLedgerPath:   "/Users/kooshapari/CodeProjects/Phenotype/repos/.airlock/beads-ledger.jsonl",
		MaxCommitsPerRun: 50,
		DryRun:            false,
		Timeout:           30 * time.Minute,
	}
}

// Run executes a single autocommit cycle against the airlock-v2 binary.
// Returns the number of commits applied and any error encountered.
func Run(ctx context.Context, opts AutocommitOptions) (int, error) {
	if opts.WorkspaceRoot == "" {
		return 0, fmt.Errorf("airlock: WorkspaceRoot required")
	}
	if _, err := os.Stat(opts.WorkspaceRoot); err != nil {
		return 0, fmt.Errorf("airlock: workspace check failed: %w", err)
	}
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}
	beadsDir := filepath.Dir(opts.BeadsLedgerPath)
	if err := os.MkdirAll(beadsDir, 0o755); err != nil {
		return 0, fmt.Errorf("airlock: mkdir beads dir: %w", err)
	}
	args := []string{"autocommit", "--workspace", opts.WorkspaceRoot}
	if opts.BeadsLedgerPath != "" {
		args = append(args, "--beads", opts.BeadsLedgerPath)
	}
	if opts.MaxCommitsPerRun > 0 {
		args = append(args, "--max-commits", fmt.Sprintf("%d", opts.MaxCommitsPerRun))
	}
	if opts.DryRun {
		args = append(args, "--dry-run")
	}
	cmd := exec.CommandContext(ctx, "airlock-v2", args...)
	cmd.Dir = opts.WorkspaceRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("airlock: airlock-v2 failed: %w", err)
	}
	return opts.MaxCommitsPerRun, nil
}
