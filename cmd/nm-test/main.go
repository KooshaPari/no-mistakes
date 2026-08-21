package main

import (
	"context"
	"fmt"
	"os"

	steps "github.com/kunchenguid/no-mistakes/internal/pipeline/steps"
	"github.com/kunchenguid/no-mistakes/internal/pipeline"
)

func main() {
	dir := os.Args[1]
	ctx := context.Background()
	sctx := &pipeline.StepContext{Ctx: ctx, WorkDir: dir}

	fmt.Println("=== worktree-gc ===")
	r, err := (&steps.WorktreeGcStep{}).Execute(sctx)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Printf("findings: %s\n", r.Findings)

	fmt.Println("=== branch-lint ===")
	r, _ = (&steps.BranchLintStep{}).Execute(sctx)
	fmt.Printf("findings: %s\n", r.Findings)

	fmt.Println("=== health-check ===")
	r, _ = (&steps.HealthCheckStep{}).Execute(sctx)
	fmt.Printf("findings: %s\n", r.Findings)
}
