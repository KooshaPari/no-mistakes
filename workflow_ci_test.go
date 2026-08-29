package main

import (
	"os"
	"slices"
	"strings"
	"testing"
)

func TestCIWorkflowRunsTestsOnAllSupportedDesktopPlatforms(t *testing.T) {
	data, err := os.ReadFile(".github/workflows/ci.yml")
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}

	content := string(data)
	for _, osName := range []string{"ubuntu-latest", "macos-latest", "windows-latest"} {
		if !strings.Contains(content, osName) {
			t.Fatalf("CI workflow must test %q", osName)
		}
	}
}

func TestCIWorkflowUsesRaceTestsOnUnixRunners(t *testing.T) {
	data, err := os.ReadFile(".github/workflows/ci.yml")
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "if: runner.os != 'Windows'") {
		t.Fatalf("CI workflow must keep the Unix test branch so macOS runs the Unix suite")
	}
	if !strings.Contains(content, "run: go test -race ./...") {
		t.Fatalf("CI workflow must run the race-enabled suite on Unix runners")
	}
}

func TestCIWorkflowEmitsStableProtectedCheckNames(t *testing.T) {
	wf := loadCIWorkflowDoc(t)

	lint, ok := wf.Jobs["check"]
	if !ok {
		t.Fatal("CI workflow has no check job")
	}
	if got, want := lint.DisplayName, "ci / lint"; got != want {
		t.Fatalf("check job name = %q, want protected context %q", got, want)
	}

	testGate, ok := wf.Jobs["test-gate"]
	if !ok {
		t.Fatal("CI workflow has no stable aggregate test gate")
	}
	if got, want := testGate.DisplayName, "ci / test"; got != want {
		t.Fatalf("aggregate test job name = %q, want protected context %q", got, want)
	}
	if got, want := testGate.needs(), []string{"test", "e2e"}; !slices.Equal(got, want) {
		t.Fatalf("aggregate test job needs = %v, want %v", got, want)
	}
	if !strings.Contains(testGate.If, "always()") {
		t.Fatalf("aggregate test job condition = %q, want an always-run verdict", testGate.If)
	}
	for _, result := range []string{"needs.test.result", "needs.e2e.result"} {
		if !strings.Contains(testGate.allRun(), result) {
			t.Fatalf("aggregate test gate must include %s in its verdict", result)
		}
	}
}
