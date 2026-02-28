package domain

import "time"

type LogStream string

const (
	LogStdout LogStream = "stdout"
	LogStderr LogStream = "stderr"
)

type LogEvent struct {
	Stream LogStream
	Line   string
	Time   time.Time
}
