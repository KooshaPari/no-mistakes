package steps

import (
	"context"
	"os/exec"
	"strings"
)

// gitTrim runs a git command in dir and returns trimmed stdout.
func gitTrim(ctx context.Context, dir string, args ...string) (string, error) {
	out, err := exec.CommandContext(ctx, "git", append([]string{"-C", dir}, args...)...).CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// escapeJSON quotes a string for inclusion in a JSON literal value.
func escapeJSON(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`, "\t", `\t`)
	return r.Replace(s)
}

// findingsFromError formats an error into a JSON Findings array.
func findingsFromError(op string, err error) string {
	return `[{"severity":"error","message":"` + escapeJSON(op+": "+err.Error()) + `"}]`
}
