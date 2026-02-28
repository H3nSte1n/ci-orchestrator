package runner

import (
	"bufio"
	"context"
	"errors"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/ports"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type HostRunner struct{}

func NewHostRunner() ports.Runner {
	return &HostRunner{}
}

func (r *HostRunner) Start(ctx context.Context, workdir string, command string, env []string) (<-chan domain.LogEvent, func() (int, error), error) {
	cmd := exec.CommandContext(ctx, "sh", "-lc", command)
	cmd.Dir = workdir
	cmd.Env = append(os.Environ(), env...)

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	events := make(chan domain.LogEvent, 200)
	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go r.scanLines(ctx, domain.LogStdout, stdout, &wg, events)
	go r.scanLines(ctx, domain.LogStderr, stderr, &wg, events)

	waitFn := func() (int, error) {
		return waitAndClose(cmd, &wg, events)
	}

	return events, waitFn, nil
}

func (r *HostRunner) scanLines(ctx context.Context, stream domain.LogStream, rc io.Reader, wg *sync.WaitGroup, events chan<- domain.LogEvent) {
	defer wg.Done()
	sc := bufio.NewScanner(rc)

	// Increase max token size (default is ~64K)
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 512*1024)

	for sc.Scan() {
		select {
		case events <- domain.LogEvent{
			Stream: stream,
			Line:   sc.Text(),
			Time:   time.Now(),
		}:
		case <-ctx.Done():
			return
		}
	}
}

func waitAndClose(cmd *exec.Cmd, wg *sync.WaitGroup, events chan domain.LogEvent) (int, error) {
	err := cmd.Wait()

	wg.Wait()
	close(events)

	if err == nil {
		return 0, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), err
	}

	// Could be start failure, signal kill, context cancel, etc.
	return -1, err
}
