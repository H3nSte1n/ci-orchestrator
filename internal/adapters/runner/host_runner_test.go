package runner

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func collectEvents(events <-chan domain.LogEvent) []domain.LogEvent {
	var out []domain.LogEvent
	for ev := range events {
		out = append(out, ev)
	}
	return out
}

func TestHostRunner_Start_Success_StreamsStdout(t *testing.T) {
	runner := NewHostRunner()
	ctx := context.Background()
	workdir := t.TempDir()

	events, waitFn, err := runner.Start(ctx, workdir, `echo "hello"`, nil)

	require.NoError(t, err)
	require.NotNil(t, waitFn)

	exitCode, runErr := waitFn()
	evs := collectEvents(events)

	assert.NoError(t, runErr)
	assert.Equal(t, 0, exitCode)

	found := false
	for _, ev := range evs {
		if ev.Stream == domain.LogStdout && strings.Contains(ev.Line, "hello") {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestHostRunner_Start_Failure_StreamsStderrAndExitCode(t *testing.T) {
	runner := NewHostRunner()
	ctx := context.Background()
	workdir := t.TempDir()

	cmd := `echo "oops" 1>&2; exit 7`
	events, waitFn, err := runner.Start(ctx, workdir, cmd, nil)

	require.NoError(t, err)
	require.NotNil(t, waitFn)

	exitCode, runErr := waitFn()
	evs := collectEvents(events)

	assert.Error(t, runErr)
	assert.Equal(t, 7, exitCode)

	found := false
	for _, ev := range evs {
		if ev.Stream == domain.LogStderr && strings.Contains(ev.Line, "oops") {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestHostRunner_Start_Cancel_ContextStopsCommand(t *testing.T) {
	runner := NewHostRunner()
	ctx, cancel := context.WithCancel(context.Background())
	workdir := t.TempDir()

	events, waitFn, err := runner.Start(ctx, workdir, `sleep 5; echo "done"`, nil)
	require.NoError(t, err)
	require.NotNil(t, waitFn)

	time.Sleep(100 * time.Millisecond)
	cancel()

	exitCode, runErr := waitFn()
	_ = collectEvents(events)

	assert.Error(t, runErr)
	_ = exitCode
}
