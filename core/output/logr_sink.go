package output

import (
	"strings"

	"github.com/go-logr/logr"
)

func NewOutputLogr(output Output) logr.Logger {
	return logr.New(&logrSink{
		output: output,
	})
}

type outputEnabled interface {
	Enabled(level int) bool
}

type logrSink struct {
	output Output
	level  int
}

func (l *logrSink) Init(info logr.RuntimeInfo) {}

func (l *logrSink) Enabled(level int) bool {
	if output, ok := l.output.(outputEnabled); ok {
		return output.Enabled(level)
	}
	return true
}

func (l *logrSink) V(level int) logr.LogSink {
	level += l.level
	return &logrSink{
		output: l.output.V(level),
		level:  level,
	}
}

func (l *logrSink) Info(level int, msg string, keysAndValues ...interface{}) {
	l.output.V(level).WithValues(keysAndValues...).Info(strings.TrimRight(msg, "\n"))
}

func (l *logrSink) Error(err error, msg string, keysAndValues ...interface{}) {
	l.output.WithValues(keysAndValues...).Error(err, strings.TrimRight(msg, "\n"))
}

func (l *logrSink) WithValues(keysAndValues ...interface{}) logr.LogSink {
	return &logrSink{
		output: l.output.WithValues(keysAndValues...),
	}
}

func (l *logrSink) WithName(name string) logr.LogSink {
	// not using the logger name
	return l
}
