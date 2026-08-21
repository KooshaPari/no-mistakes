package steps

import (
	"context"
	"os/exec"
	"strconv"
	"strings"

	"github.com/kunchenguid/no-mistakes/internal/pipeline"
	"github.com/kunchenguid/no-mistakes/internal/types"
)

// BranchLintStep flags branches merged to main or with [gone] upstream.
type BranchLintStep struct{}

func (s *BranchLintStep) Name() types.StepName { return types.StepBranchLint }

func (s *BranchLintStep) Execute(sctx *pipeline.StepContext) (*pipeline.StepOutcome, error) {
	if sctx == nil || sctx.Ctx == nil {
		return &pipeline.StepOutcome{}, nil
	}
	var findings []string
	var gone, merged int

	upstreamOut, err := exec.CommandContext(sctx.Ctx, "git", "-C", sctx.WorkDir,
		"for-each-ref", "--format=%(refname:short) %(upstream:track)", "refs/heads/").CombinedOutput()
	if err == nil {
		for _, ln := range strings.Split(strings.TrimSpace(string(upstreamOut)), "\n") {
			ln = strings.TrimSpace(ln)
			if ln == "" {
				continue
			}
			parts := strings.SplitN(ln, " ", 2)
			if len(parts) != 2 {
				continue
			}
			branch, track := parts[0], parts[1]
			if branch == "main" || branch == "master" {
				continue
			}
			if track == "gone" {
				gone++
				findings = append(findings,
					`{"severity":"warn","message":"upstream is gone: `+escapeJSON(branch)+`"}`)
			}
		}
	}

	mergedOut, err := exec.CommandContext(sctx.Ctx, "git", "-C", sctx.WorkDir,
		"branch", "--format=%(refname:short)", "--merged", "main").CombinedOutput()
	if err == nil {
		for _, ln := range strings.Split(strings.TrimSpace(string(mergedOut)), "\n") {
			ln = strings.TrimSpace(ln)
			if ln == "" || ln == "main" || ln == "master" {
				continue
			}
			merged++
			findings = append(findings,
				`{"severity":"info","message":"merged into main: `+escapeJSON(ln)+`"}`)
		}
	}

	sev := "info"
	if gone > 0 {
		sev = "warn"
	}
	msg := strconv.Itoa(len(findings)) + " finding(s)"
	return &pipeline.StepOutcome{
		Findings: `[{"severity":"` + sev + `","message":"` + msg + `"}]`,
	}, nil
}

var _ = context.Background
