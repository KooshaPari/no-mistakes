package airlock

import (
	"context"
	"os/exec"
	"testing"
	"time"
)

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if opts.WorkspaceRoot == "" {
		t.Fatal("WorkspaceRoot must be set")
	}
	if opts.Timeout <= 0 {
		t.Fatal("Timeout must be positive")
	}
}

func TestRun_RequiresWorkspace(t *testing.T) {
	_, err := Run(context.Background(), AutocommitOptions{})
	if err == nil {
		t.Fatal("expected error when WorkspaceRoot empty")
	}
}

func TestRun_RejectsBadPath(t *testing.T) {
	_, err := Run(context.Background(), AutocommitOptions{
		WorkspaceRoot: "/nonexistent/path/here",
		Timeout:       100 * time.Millisecond,
	})
	if err == nil {
		t.Fatal("expected error for bad path")
	}
}

// Skipped live test — only run if airlock-v2 binary exists
func TestRun_LiveSkipped(t *testing.T) {
	if _, err := exec.LookPath("airlock-v2"); err != nil {
		t.Skip("airlock-v2 not in PATH; live test skipped")
	}
	opts := DefaultOptions()
	opts.WorkspaceRoot = "/tmp"
	opts.Timeout = 500 * time.Millisecond
	opts.DryRun = true
	_, err := Run(context.Background(), opts)
	if err != nil {
		t.Logf("airlock-v2 not callable: %v", err)
	}
}
