package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const rpcTimeout = 250 * time.Millisecond
const maxEventBytes = 1 << 20
const gitWorkBudget = 350 * time.Millisecond
const maxGitOutput = 256 << 10
const maxGitPathOutput = 8 << 10

type cappedBuffer struct {
	bytes.Buffer
	limit int
}

func (b *cappedBuffer) Write(p []byte) (int, error) {
	if len(p) > b.limit-b.Len() {
		return 0, io.ErrShortWrite
	}
	return b.Buffer.Write(p)
}

func gitOutput(ctx context.Context, limit int, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	var stdout cappedBuffer
	stdout.limit = limit
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return stdout.Bytes(), nil
}

type event struct {
	Agent  string `json:"agent"`
	Hook   string `json:"hook"`
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
	Dirty  int    `json:"dirty"`
	TS     string `json:"ts"`
}

func repoMetadata(cwd string) (repo, branch string, dirty int, ok bool) {
	if cwd == "" || !filepath.IsAbs(cwd) {
		return "", "", 0, false
	}
	info, err := os.Stat(cwd)
	if err != nil || !info.IsDir() {
		return "", "", 0, false
	}
	ctx, cancel := context.WithTimeout(context.Background(), gitWorkBudget)
	defer cancel()
	root, err := gitOutput(ctx, maxGitPathOutput, "-C", cwd, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", "", 0, false
	}
	repo = strings.TrimSpace(string(root))
	name, err := gitOutput(ctx, maxGitPathOutput, "-C", repo, "branch", "--show-current")
	if err != nil {
		return "", "", 0, false
	}
	branch = strings.TrimSpace(string(name))
	if branch == "" {
		return "", "", 0, false
	}
	status, err := gitOutput(ctx, maxGitOutput, "-C", repo, "status", "--porcelain", "-z")
	if err != nil {
		return "", "", 0, false
	}
	entries := 0
	rows := strings.Split(string(status), "\x00")
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		if row == "" {
			continue
		}
		entries++
		if len(row) >= 2 && (row[0] == 'R' || row[0] == 'C' || row[1] == 'R' || row[1] == 'C') {
			i++
		}
	}
	return repo, branch, entries, true
}

func socketPath() string {
	if path := os.Getenv("AIRLOCK_CODEX_HOOK_SOCKET"); path != "" {
		return path
	}
	if path := os.Getenv("AIRLOCK_SOCKET"); path != "" {
		return path
	}
	return filepath.Join(os.Getenv("HOME"), ".airlock", "codex-hook.sock")
}

func notify(path string, payload []byte) error {
	if path == "" || len(payload) == 0 {
		return errors.New("observer socket path and payload are required")
	}
	conn, err := net.DialTimeout("unix", path, rpcTimeout)
	if err != nil {
		return fmt.Errorf("connect to observer socket: %w", err)
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(rpcTimeout)); err != nil {
		return fmt.Errorf("set observer socket deadline: %w", err)
	}
	if err := writeAll(conn, append(append([]byte(nil), payload...), '\n')); err != nil {
		return fmt.Errorf("write observer event: %w", err)
	}
	return nil
}

func writeAll(w io.Writer, data []byte) error {
	for len(data) > 0 {
		n, err := w.Write(data)
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrShortWrite
		}
		data = data[n:]
	}
	return nil
}

func run(input io.Reader, path string) error {
	data, err := io.ReadAll(io.LimitReader(input, maxEventBytes+1))
	if err != nil || len(data) == 0 || len(data) > maxEventBytes {
		return nil
	}
	var inputEvent struct {
		HookEventName string `json:"hook_event_name"`
		ToolName      string `json:"tool_name"`
		CWD           string `json:"cwd"`
	}
	if json.Unmarshal(data, &inputEvent) != nil || inputEvent.HookEventName != "PostToolUse" {
		return nil
	}
	repo, branch, dirtyCount, ok := repoMetadata(inputEvent.CWD)
	if !ok {
		return nil
	}
	metadata, err := json.Marshal(event{"codex", inputEvent.ToolName, repo, branch, dirtyCount, time.Now().UTC().Format(time.RFC3339Nano)})
	if err != nil {
		return fmt.Errorf("encode observer event: %w", err)
	}
	return notify(path, metadata)
}

func main() {
	if err := run(os.Stdin, socketPath()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
