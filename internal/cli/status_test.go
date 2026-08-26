package cli

import (
	"testing"

	"github.com/kunchenguid/no-mistakes/internal/branchsync"
	"github.com/kunchenguid/no-mistakes/internal/db"
)

func TestCachedBranchSummary(t *testing.T) {
	tests := []struct {
		name  string
		state branchsync.State
		want  string
	}{
		{
			name: "clean branch",
			state: branchsync.State{
				State: branchsync.StateSynchronized,
				Local: branchsync.LocalState{Branch: "feature/state", Head: "0123456789abcdef", Clean: true},
			},
			want: "cached: feature/state 01234567 (clean; already synchronized with the pipeline-pushed head)",
		},
		{
			name: "dirty branch",
			state: branchsync.State{
				State: branchsync.StateDirty,
				Local: branchsync.LocalState{Branch: "feature/state", Head: "fedcba9876543210", Reason: "uncommitted changes"},
			},
			want: "cached: feature/state fedcba98 (dirty: uncommitted changes; dirty)",
		},
		{
			name:  "unavailable",
			state: branchsync.State{State: branchsync.StateAmbiguousContext},
			want:  "cached: unavailable (ambiguous context)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cachedBranchSummary(tt.state); got != tt.want {
				t.Fatalf("cachedBranchSummary() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStatusFingerprintIncludesCachedSummary(t *testing.T) {
	run := &db.Run{ID: "run-1", Branch: "feature/test", Status: "running", HeadSHA: "head-one"}
	before := statusFingerprint("repo", "running", run, "cached: main 01234567 (clean; synchronized)")
	after := statusFingerprint("repo", "running", run, "cached: main 89abcdef (dirty; dirty)")
	if before == after {
		t.Fatal("changing displayed cached evidence must change the status fingerprint")
	}
}
