// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package output

import (
	"strings"

	"github.com/go-logr/logr"
)

func NewOutputLogr(output Output) logr.Logger {
	return &logrLogger{
		output: output,
	}
}

type outputEnabled interface {
	Enabled(level int) bool
}

type logrLogger struct {
	output Output
	level  int
}

func (l *logrLogger) Enabled() bool {
	if output, ok := l.output.(outputEnabled); ok {
		return output.Enabled(l.level)
	}
	return true
}

func (l *logrLogger) Info(msg string, keysAndValues ...interface{}) {
	l.output.WithValues(keysAndValues...).Info(strings.TrimRight(msg, "\n"))
}

func (l *logrLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.output.WithValues(keysAndValues...).Error(err, strings.TrimRight(msg, "\n"))
}

func (l *logrLogger) V(level int) logr.Logger {
	level += l.level
	return &logrLogger{
		output: l.output.V(level),
		level:  level,
	}
}

func (l *logrLogger) WithValues(keysAndValues ...interface{}) logr.Logger {
	return &logrLogger{
		output: l.output.WithValues(keysAndValues...),
		level:  l.level,
	}
}

func (l *logrLogger) WithName(name string) logr.Logger {
	// not using the logger name
	return l
}
