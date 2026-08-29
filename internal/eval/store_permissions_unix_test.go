//go:build unix

package eval

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteJSONKeepsEvaluationArtifactsPrivate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "decision.json")

	for _, value := range []any{
		map[string]string{"decision": "first"},
		map[string]string{"decision": "replacement"},
	} {
		if err := writeJSON(path, value); err != nil {
			t.Fatalf("writeJSON: %v", err)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat evaluation artifact: %v", err)
		}
		if got, want := info.Mode().Perm(), os.FileMode(0o600); got != want {
			t.Fatalf("evaluation artifact permissions = %04o, want %04o", got, want)
		}
	}
}
