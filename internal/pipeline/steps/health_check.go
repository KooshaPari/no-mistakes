package steps

import (
	"context"
	"strconv"
	"strings"

	"github.com/kunchenguid/no-mistakes/internal/pipeline"
	"github.com/kunchenguid/no-mistakes/internal/types"
)

// HealthCheckStep flags a worktree that is not on main, has dirty tree, or
// is ahead/behind origin/main.
type HealthCheckStep struct{}

func (s *HealthCheckStep) Name() types.StepName { return types.StepHealthCheck }

func (s *HealthCheckStep) Execute(sctx *pipeline.StepContext) (*pipeline.StepOutcome, error) {
	if sctx == nil || sctx.Ctx == nil {
		return &pipeline.StepOutcome{}, nil
	}
	var problems []string
	if branch, err := gitTrim(sctx.Ctx, sctx.WorkDir, "rev-parse", "--abbrev-ref", "HEAD"); err == nil {
		branch = strings.TrimSpace(branch)
		if branch != "" && branch != "main" && branch != "master" {
			problems = append(problems, "not on main (HEAD="+branch+")")
		}
	}
	if status, err := gitTrim(sctx.Ctx, sctx.WorkDir, "status", "--porcelain"); err == nil {
		if status != "" {
			problems = append(problems, "dirty ("+strconv.Itoa(strings.Count(status, "\n")+1)+" files)")
		}
	}
	if ahead, err := gitTrim(sctx.Ctx, sctx.WorkDir, "rev-list", "--count", "origin/main..HEAD"); err == nil {
		if n, e := strconv.Atoi(strings.TrimSpace(ahead)); e == nil && n > 0 {
			problems = append(problems, "ahead "+strconv.Itoa(n))
		}
	}
	if behind, err := gitTrim(sctx.Ctx, sctx.WorkDir, "rev-list", "--count", "HEAD..origin/main"); err == nil {
		if n, e := strconv.Atoi(strings.TrimSpace(behind)); e == nil && n > 0 {
			problems = append(problems, "behind "+strconv.Itoa(n))
		}
	}
	sev := "info"
	msg := "healthy"
	if len(problems) > 0 {
		sev = "warn"
		msg = strings.Join(problems, "; ")
	}
	return &pipeline.StepOutcome{
		Findings: `[{"severity":"` + sev + `","message":"` + escapeJSON(msg) + `"}]`,
	}, nil
}

var _ = context.Background
