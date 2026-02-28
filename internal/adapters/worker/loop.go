package worker

import (
	"context"
	"fmt"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/ports"
	"os"
	"time"
)

type worker struct {
	workerId        string
	buildService    ports.BuildService
	buildLogService ports.BuildLogService
	interval        time.Duration
	runner          ports.Runner
}

func NewWorker(workerId string, buildService ports.BuildService, buildLogService ports.BuildLogService, interval time.Duration, runner ports.Runner) *worker {
	return &worker{
		workerId:        workerId,
		buildService:    buildService,
		buildLogService: buildLogService,
		interval:        interval,
		runner:          runner,
	}
}

func (w *worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := w.claimAndProcess(ctx)
			if err != nil {
				fmt.Println("Error claiming build:", err)
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (w *worker) claimAndProcess(ctx context.Context) error {
	build, err := w.buildService.ClaimNext(ctx, w.workerId)
	if err != nil {
		return err
	}

	if build == nil {
		return nil
	}

	workdir := fmt.Sprintf("/tmp/ci-orchestrator/%s", build.ID)
	if err := os.MkdirAll(workdir, 0o755); err != nil {
		finishedAt := time.Now()
		runErr := fmt.Errorf("create workdir: %w", err)
		return w.buildService.CompleteBuild(ctx, build.ID, -1, &finishedAt, runErr)
	}
	defer os.RemoveAll(workdir)

	events, waitFn, err := w.runner.Start(ctx, workdir, build.Command, nil)

	if err != nil {
		finishedAt := time.Now()
		runErr := fmt.Errorf("start runner: %w", err)
		return w.buildService.CompleteBuild(ctx, build.ID, -1, &finishedAt, runErr)
	}

	logErrCh := make(chan error, 1)
	go w.persistLogs(ctx, events, build.ID, logErrCh)

	exitCode, runErr := waitFn()
	finishedAt := time.Now()
	logErr := <-logErrCh
	if logErr != nil && runErr == nil {
		runErr = fmt.Errorf("persist logs: %w", logErr)
	}

	return w.buildService.CompleteBuild(ctx, build.ID, exitCode, &finishedAt, runErr)
}

func (w *worker) persistLogs(ctx context.Context, events <-chan domain.LogEvent, buildId string, logErrCh chan<- error) {
	var firstErr error

	for ev := range events {
		if firstErr != nil {
			continue
		}

		if err := w.buildLogService.AppendLog(ctx, buildId, ev); err != nil {
			firstErr = err
		}
	}

	logErrCh <- firstErr
}
