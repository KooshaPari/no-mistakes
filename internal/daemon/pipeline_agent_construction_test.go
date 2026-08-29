package daemon

import (
	"os"
	"strings"
	"testing"
)

// TestPipelineAgentConstructionHasOneOwner keeps new runs and recovered runs
// on the same agent-construction path. Duplicating this wiring has previously
// allowed profile, steering, fallback cleanup, and gate-neutralization behavior
// to drift between the two run lifecycles.
func TestPipelineAgentConstructionHasOneOwner(t *testing.T) {
	source, err := os.ReadFile("manager.go")
	if err != nil {
		t.Fatal(err)
	}

	if got := strings.Count(string(source), "agent.NewWithOptions("); got != 1 {
		t.Fatalf("manager.go has %d agent.NewWithOptions call sites, want one shared construction owner", got)
	}
}
