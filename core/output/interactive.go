// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package output

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

const (
	termRed    = "\x1b[31m"
	termGreen  = "\x1b[32m"
	termYellow = "\x1b[33m"
	termReset  = "\x1b[0m"
)

func NewInteractiveShell(out, errOut io.Writer, verbosity int) Output {
	return &interactiveShellOutput{
		out:       out,
		errOut:    newSpinner(errOut),
		verbosity: verbosity,
		level:     0,
	}
}

type interactiveShellOutput struct {
	out    io.Writer
	errOut *spinner
	// verbosity is the maximum V level that is printed to output
	verbosity int
	// level is the V level of this instance
	level         int
	status        string
	gauge         *ProgressGauge
	keysAndValues []interface{}
	lock          sync.Mutex
}

func (o *interactiveShellOutput) Info(msg string) {
	if o.level > 0 {
		msg += formatKeysAndValues(o.keysAndValues)
	}
	fmt.Fprintln(o.errOut, msg)
}

func (o *interactiveShellOutput) Infof(format string, args ...interface{}) {
	o.Info(fmt.Sprintf(format, args...))
}

func (o *interactiveShellOutput) InfoWriter() io.Writer {
	return msgWriter(o.Info)
}

func (o *interactiveShellOutput) Warn(msg string) {
	if o.level > 0 {
		msg += formatKeysAndValues(o.keysAndValues)
	}
	fmt.Fprintln(o.errOut, termYellow+msg+termReset)
}

func (o *interactiveShellOutput) Warnf(format string, args ...interface{}) {
	o.Warn(fmt.Sprintf(format, args...))
}

func (o *interactiveShellOutput) WarnWriter() io.Writer {
	return msgWriter(o.Warn)
}

func (o *interactiveShellOutput) Error(err error, msg string) {
	output := ""
	switch {
	case err == nil:
		output = msg
	case msg == "":
		output = err.Error()
	default:
		output = fmt.Sprintf("%s: %s", msg, err.Error())
	}
	if o.level > 0 {
		output += formatKeysAndValues(o.keysAndValues)
	}
	fmt.Fprintln(o.errOut, termRed+output+termReset)
}

func (o *interactiveShellOutput) Errorf(err error, format string, args ...interface{}) {
	o.Error(err, fmt.Sprintf(format, args...))
}

func (o *interactiveShellOutput) ErrorWriter() io.Writer {
	return msgWriter(func(msg string) {
		o.Error(nil, msg)
	})
}

func (o *interactiveShellOutput) StartOperation(status string) {
	o.EndOperation(true)

	o.lock.Lock()
	defer o.lock.Unlock()

	o.status = status
	o.errOut.SetSuffix(fmt.Sprintf(" %s ", o.status))
	o.errOut.Start()
}

func (o *interactiveShellOutput) StartOperationWithProgress(gauge *ProgressGauge) {
	o.EndOperation(true)

	o.lock.Lock()
	defer o.lock.Unlock()

	o.gauge = gauge
	o.errOut.SetProgressGauge(gauge)
	o.errOut.Start()
}

func (o *interactiveShellOutput) EndOperation(success bool) {
	o.lock.Lock()
	defer o.lock.Unlock()

	status := o.status
	if o.gauge != nil {
		status = strings.TrimPrefix(o.gauge.String(), " ")
	}

	if status == "" {
		return
	}
	o.errOut.Stop()
	fmt.Fprint(o.errOut, "\r")
	if success {
		fmt.Fprintf(o.errOut, " %s✓%s %s\n", termGreen, termReset, status)
	} else {
		fmt.Fprintf(o.errOut, " %s✗%s %s\n", termRed, termReset, status)
	}
	o.status = ""
	o.gauge = nil
}

func (o *interactiveShellOutput) Result(result string) {
	fmt.Fprintln(o.out, result)
}

func (o *interactiveShellOutput) ResultWriter() io.Writer {
	return o.out
}

func (o *interactiveShellOutput) Enabled(level int) bool {
	return level <= o.verbosity
}

func (o *interactiveShellOutput) V(level int) Output {
	if !o.Enabled(level) {
		return &noopOutput{Output: o}
	}
	return &interactiveShellOutput{
		out:           o.out,
		errOut:        o.errOut,
		verbosity:     o.verbosity,
		level:         level,
		keysAndValues: o.keysAndValues,
	}
}

func (o *interactiveShellOutput) WithValues(keysAndValues ...interface{}) Output {
	return &interactiveShellOutput{
		out:           o.out,
		errOut:        o.errOut,
		verbosity:     o.verbosity,
		level:         o.level,
		keysAndValues: append(o.keysAndValues, keysAndValues...),
	}
}

type msgWriter func(msg string)

func (w msgWriter) Write(p []byte) (n int, err error) {
	w(strings.TrimSuffix(string(p), "\n"))
	return len(p), nil
}
