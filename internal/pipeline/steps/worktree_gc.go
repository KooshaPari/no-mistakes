package steps

import (
	"context"
	"os/exec"
	"strconv"
	"strings"

	"github.com/kunchenguid/no-mistakes/internal/pipeline"
	"github.com/kunchenguid/no-mistakes/internal/types"
)

// WorktreeGcStep flags stale git worktree pointers whose dir is gone.
type WorktreeGcStep struct{}

func (s *WorktreeGcStep) Name() types.StepName { return types.StepWorktreeGc }

func (s *WorktreeGcStep) Execute(sctx *pipeline.StepContext) (*pipeline.StepOutcome, error) {
	if sctx == nil || sctx.Ctx == nil {
		return &pipeline.StepOutcome{}, nil
	}
	out, err := exec.CommandContext(sctx.Ctx, "git", "-C", sctx.WorkDir,
		"worktree", "list", "--porcelain").CombinedOutput()
	if err != nil {
		return &pipeline.StepOutcome{
			Findings: findingsFromError("worktree-gc", err),
			ExitCode: 1,
		}, nil
	}
	outStr := strings.TrimSpace(string(out))
	if outStr == "" {
		return &pipeline.StepOutcome{
			Findings: `[{"severity":"info","message":"no worktrees configured"}]`,
		}, nil
	}
	var stale, active int
	for _, block := range strings.Split(outStr, "\n\n") {
		var path string
		for _, ln := range strings.Split(block, "\n") {
			if strings.HasPrefix(ln, "worktree ") {
				path = strings.TrimPrefix(ln, "worktree ")
				break
			}
		}
		if path == "" {
			continue
		}
		check := exec.CommandContext(sctx.Ctx, "git", "-C", path, "rev-parse", "--git-dir")
		if err := check.Run(); err != nil {
			stale++
		} else {
			active++
		}
	}
	sev := "info"
	msg := strconv.Itoa(active) + " active worktree(s)"
	if stale > 0 {
		sev = "warn"
		msg = strconv.Itoa(stale) + " stale, " + strconv.Itoa(active) + " active"
	}
	return &pipeline.StepOutcome{
		Findings: `[{"severity":"` + sev + `","message":"` + escapeJSON(msg) + `"}]`,
	}, nil
}

var _ = context.Background
