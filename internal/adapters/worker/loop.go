package worker

import (
	"context"
	"fmt"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/ports"
	"time"
)

type worker struct {
	workerId     string
	buildService ports.BuildService
	interval     time.Duration
}

func NewWorker(workerId string, buildService ports.BuildService, interval time.Duration) *worker {
	return &worker{
		workerId:     workerId,
		buildService: buildService,
		interval:     interval,
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
		fmt.Println("No builds to claim, waiting...")
		return nil
	}

	// TODO: Execute build command and update status accordingly
	return nil
}
