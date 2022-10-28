// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package output

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

func NewNonInteractiveShell(out, errOut io.Writer, verbosity int) Output {
	return &nonInteractiveShellOutput{
		out:       out,
		errOut:    errOut,
		verbosity: verbosity,
	}
}

// formatExtended formats a log message in the following format:
// <timestamp> <level> <message> [key=value]
// e.g. 2022-01-10 18:07:49 INF some message key=value
// .
func formatExtended(level string, msg string, keysAndValues []interface{}) string {
	result := fmt.Sprintf("%s %s %s", time.Now().Format("2006-01-02 15:04:05"), level, msg)
	result += formatKeysAndValues(keysAndValues)
	return result
}

func formatKeysAndValues(keysAndValues []interface{}) string {
	if len(keysAndValues) == 0 {
		return ""
	}
	result := "   "
	for i := 1; i < len(keysAndValues); i += 2 {
		key := keysAndValues[i-1]
		value := fmt.Sprintf("%v", keysAndValues[i])
		if strings.Contains(value, " ") || strings.Contains(value, "=") {
			value = `"` + value + `"`
		}
		result += fmt.Sprintf(" %s=%v", key, value)
	}
	return result
}

type nonInteractiveShellOutput struct {
	out           io.Writer
	errOut        io.Writer
	verbosity     int
	status        string
	keysAndValues []interface{}
	lock          sync.Mutex
}

func (o *nonInteractiveShellOutput) Info(msg string) {
	fmt.Fprintln(o.errOut, formatExtended("INF", msg, o.keysAndValues))
}

func (o *nonInteractiveShellOutput) Infof(format string, args ...interface{}) {
	o.Info(fmt.Sprintf(format, args...))
}

func (o *nonInteractiveShellOutput) InfoWriter() io.Writer {
	return msgWriter(o.Info)
}

func (o *nonInteractiveShellOutput) Warn(msg string) {
	fmt.Fprintln(o.errOut, formatExtended("WRN", msg, o.keysAndValues))
}

func (o *nonInteractiveShellOutput) Warnf(format string, args ...interface{}) {
	o.Warn(fmt.Sprintf(format, args...))
}

func (o *nonInteractiveShellOutput) WarnWriter() io.Writer {
	return msgWriter(o.Warn)
}

func (o *nonInteractiveShellOutput) Error(err error, msg string) {
	fmt.Fprintln(o.errOut, formatExtended("ERR", msg, append([]interface{}{"err", err}, o.keysAndValues...)))
}

func (o *nonInteractiveShellOutput) Errorf(err error, format string, args ...interface{}) {
	o.Error(err, fmt.Sprintf(format, args...))
}

func (o *nonInteractiveShellOutput) ErrorWriter() io.Writer {
	return msgWriter(func(msg string) {
		o.Error(nil, msg)
	})
}

func (o *nonInteractiveShellOutput) StartOperation(status string) {
	o.EndOperation(true)

	o.lock.Lock()
	defer o.lock.Unlock()

	o.status = status
	o.Infof(" • %s...", o.status)
}

func (o *nonInteractiveShellOutput) StartOperationWithProgress(gauge *ProgressGauge) {
	o.EndOperation(true)

	o.lock.Lock()
	defer o.lock.Unlock()

	o.status = strings.TrimPrefix(gauge.String(), " ")
	o.Infof(" • %s...", o.status)
}

func (o *nonInteractiveShellOutput) EndOperation(success bool) {
	o.lock.Lock()
	defer o.lock.Unlock()

	if o.status == "" {
		return
	}
	if success {
		o.Infof(" ✓ %s", o.status)
	} else {
		o.Infof(" ✗ %s", o.status)
	}
	o.status = ""
}

func (o *nonInteractiveShellOutput) Result(result string) {
	fmt.Fprintln(o.out, result)
}

func (o *nonInteractiveShellOutput) ResultWriter() io.Writer {
	return o.out
}

func (o *nonInteractiveShellOutput) Enabled(level int) bool {
	return level <= o.verbosity
}

func (o *nonInteractiveShellOutput) V(level int) Output {
	if !o.Enabled(level) {
		return &noopOutput{Output: o}
	}
	return &nonInteractiveShellOutput{
		out:           o.out,
		errOut:        o.errOut,
		verbosity:     o.verbosity,
		keysAndValues: o.keysAndValues,
	}
}

func (o *nonInteractiveShellOutput) WithValues(keysAndValues ...interface{}) Output {
	return &nonInteractiveShellOutput{
		out:           o.out,
		errOut:        o.errOut,
		verbosity:     o.verbosity,
		keysAndValues: append(o.keysAndValues, keysAndValues...),
	}
}
