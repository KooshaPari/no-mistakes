package main

import (
	"slices"
	"strings"
	"testing"
)

func TestCIWorkflowRunsTestsOnAllSupportedDesktopPlatforms(t *testing.T) {
	wf := loadCIWorkflowDoc(t)
	testJob, ok := wf.Jobs["test"]
	if !ok {
		t.Fatal("CI workflow has no test job")
	}
	if got, want := workflowExpression(wfScalar(testJob.RunsOn)), "matrix.os"; got != want {
		t.Fatalf("CI test job runner expression = %q, want %q", got, want)
	}

	platforms := make(map[string]bool)
	for _, leg := range testJob.Strategy.Matrix.Include {
		platforms[leg["os"]] = true
	}
	for _, osName := range []string{"ubuntu-latest", "macos-latest", "windows-latest"} {
		if !platforms[osName] {
			t.Fatalf("CI workflow must test %q", osName)
		}
	}
}

func workflowExpression(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "${{")
	value = strings.TrimSuffix(value, "}}")
	return strings.TrimSpace(value)
}

func TestCIWorkflowUsesRaceTestsOnUnixRunners(t *testing.T) {
	wf := loadCIWorkflowDoc(t)
	testJob, ok := wf.Jobs["test"]
	if !ok {
		t.Fatal("CI workflow has no test job")
	}

	for _, step := range testJob.Steps {
		if unixOnly(step.If) && slices.Equal(strings.Fields(step.Run), []string{"go", "test", "-race", "./..."}) {
			return
		}
	}
	t.Fatal("CI test job must run go test -race ./... on Unix runners")
}

func unixOnly(condition string) bool {
	condition = strings.TrimSpace(condition)
	condition = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(condition, "${{"), "}}"))
	fields := strings.Fields(condition)
	return len(fields) == 3 && fields[0] == "runner.os" && fields[1] == "!=" &&
		(fields[2] == "'Windows'" || fields[2] == `"Windows"`)
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
