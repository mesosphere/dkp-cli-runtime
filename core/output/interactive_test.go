// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package output_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jwalton/gchalk"
	"github.com/stretchr/testify/assert"

	"github.com/mesosphere/dkp-cli-runtime/core/output"
)

const (
	termRed       = "\x1b[31m"
	termGreen     = "\x1b[32m"
	termYellow    = "\x1b[33m"
	termDefaultFg = "\x1b[39m"
	termClearLine = "\x1b[2K"
)

func TestInteractiveShellOutput(t *testing.T) {
	assert := assert.New(t)

	origGchalkStderr := gchalk.Stderr
	defer func() {
		gchalk.Stderr = origGchalkStderr
	}()
	gchalk.Stderr = gchalk.New(
		gchalk.ForceLevel(gchalk.LevelAnsi256),
	)

	t.Run("default", func(t *testing.T) {
		out := bytes.Buffer{}
		errOut := bytes.Buffer{}
		tOutput := output.NewInteractiveShell(&out, &errOut, 0)

		tOutput.Info("info message")
		assert.Empty(out.String())
		assert.Equal("info message\n", errOut.String())
		errOut.Reset()

		tOutput.Infof("info %s", "message")
		assert.Empty(out.String())
		assert.Equal("info message\n", errOut.String())
		errOut.Reset()

		n, err := io.WriteString(tOutput.InfoWriter(), "info message")
		assert.Equal(len("info message"), n)
		assert.NoError(err)
		assert.Empty(out.String())
		assert.Equal("info message\n", errOut.String())
		errOut.Reset()

		tOutput.Warn("warning message")
		assert.Empty(out.String())
		assert.Equal(termYellow+"warning message"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()

		tOutput.Warnf("warning %s", "message")
		assert.Empty(out.String())
		assert.Equal(termYellow+"warning message"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()

		n, err = io.WriteString(tOutput.WarnWriter(), "warning message")
		assert.Equal(len("warning message"), n)
		assert.NoError(err)
		assert.Empty(out.String())
		assert.Equal(termYellow+"warning message"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()

		tOutput.Error(fmt.Errorf("error message"), "an error happened")
		assert.Empty(out.String())
		assert.Equal(termRed+"an error happened: error message"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()
		tOutput.Error(nil, "an error happened")
		assert.Empty(out.String())
		assert.Equal(termRed+"an error happened"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()
		tOutput.Error(fmt.Errorf("error message"), "")
		assert.Empty(out.String())
		assert.Equal(termRed+"error message"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()

		tOutput.Errorf(fmt.Errorf("error message"), "an error %s", "happened")
		assert.Empty(out.String())
		assert.Equal(termRed+"an error happened: error message"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()

		n, err = io.WriteString(tOutput.ErrorWriter(), "an error happened")
		assert.Equal(len("an error happened"), n)
		assert.NoError(err)
		assert.Empty(out.String())
		assert.Equal(termRed+"an error happened"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()

		tOutput.Result("a result")
		assert.Equal("a result\n", out.String())
		assert.Empty(errOut.String())
		out.Reset()

		n, err = io.WriteString(tOutput.ResultWriter(), "a result")
		assert.Equal(len("a result"), n)
		assert.NoError(err)
		assert.Equal("a result", out.String())
		assert.Empty(errOut.String())
		out.Reset()

		tOutput.WithValues("key", "value").Info("info message")
		assert.Empty(out.String())
		assert.Equal("info message\n", errOut.String())
		errOut.Reset()

		tOutput.WithValues("key", "value").Warn("warning message")
		assert.Empty(out.String())
		assert.Equal(termYellow+"warning message"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()
	})

	t.Run("verbosity hidden", func(t *testing.T) {
		out := bytes.Buffer{}
		errOut := bytes.Buffer{}
		tOutput := output.NewInteractiveShell(&out, &errOut, 0)

		tOutput.V(1).Info("info message")
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		tOutput.V(1).Infof("info %s", "message")
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		_, err := io.WriteString(tOutput.V(1).InfoWriter(), "info message")
		assert.NoError(err)
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		tOutput.V(1).Warn("warning message")
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		tOutput.V(1).Warnf("warning %s", "message")
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		_, err = io.WriteString(tOutput.V(1).WarnWriter(), "warning message")
		assert.NoError(err)
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		tOutput.V(1).Error(fmt.Errorf("error message"), "an error happened")
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		tOutput.V(1).Errorf(fmt.Errorf("error message"), "an error %s", "happened")
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		_, err = io.WriteString(tOutput.V(1).ErrorWriter(), "an error happened")
		assert.NoError(err)
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		tOutput.V(1).Result("a result")
		assert.Empty(out.String())
		assert.Equal("", out.String())
		assert.Empty(errOut.String())
		out.Reset()

		_, err = io.WriteString(tOutput.V(1).ResultWriter(), "a result")
		assert.NoError(err)
		assert.Equal("", out.String())
		assert.Empty(errOut.String())
		out.Reset()

		tOutput.V(1).WithValues("key", "value").Info("info message")
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		tOutput.V(1).WithValues("key", "value").Warn("warning message")
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		tOutput.V(1).StartOperation("working")
		tOutput.V(1).EndOperation(true)
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		tOutput.V(1).StartOperation("working")
		tOutput.V(1).EndOperationWithStatus(output.Success())
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()
	})

	t.Run("verbosity", func(t *testing.T) {
		out := bytes.Buffer{}
		errOut := bytes.Buffer{}
		tOutput := output.NewInteractiveShell(&out, &errOut, 1)

		tOutput.V(1).Info("info message")
		assert.Empty(out.String())
		assert.Equal("info message\n", errOut.String())
		errOut.Reset()

		tOutput.V(1).Infof("info %s", "message")
		assert.Empty(out.String())
		assert.Equal("info message\n", errOut.String())
		errOut.Reset()

		n, err := io.WriteString(tOutput.V(1).InfoWriter(), "info message")
		assert.Equal(len("info message"), n)
		assert.NoError(err)
		assert.Empty(out.String())
		assert.Equal("info message\n", errOut.String())
		errOut.Reset()

		tOutput.V(1).Warn("warning message")
		assert.Empty(out.String())
		assert.Equal(termYellow+"warning message"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()

		tOutput.V(1).Warnf("warning %s", "message")
		assert.Empty(out.String())
		assert.Equal(termYellow+"warning message"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()

		n, err = io.WriteString(tOutput.V(1).WarnWriter(), "warning message")
		assert.Equal(len("warning message"), n)
		assert.NoError(err)
		assert.Empty(out.String())
		assert.Equal(termYellow+"warning message"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()

		tOutput.V(1).Error(fmt.Errorf("error message"), "an error happened")
		assert.Empty(out.String())
		assert.Equal(termRed+"an error happened: error message"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()
		tOutput.V(1).Error(nil, "an error happened")
		assert.Empty(out.String())
		assert.Equal(termRed+"an error happened"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()
		tOutput.V(1).Error(fmt.Errorf("error message"), "")
		assert.Empty(out.String())
		assert.Equal(termRed+"error message"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()

		tOutput.V(1).Errorf(fmt.Errorf("error message"), "an error %s", "happened")
		assert.Empty(out.String())
		assert.Equal(termRed+"an error happened: error message"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()

		n, err = io.WriteString(tOutput.V(1).ErrorWriter(), "an error happened")
		assert.Equal(len("an error happened"), n)
		assert.NoError(err)
		assert.Empty(out.String())
		assert.Equal(termRed+"an error happened"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()

		tOutput.V(1).Result("a result")
		assert.Equal("a result\n", out.String())
		assert.Empty(errOut.String())
		out.Reset()

		n, err = io.WriteString(tOutput.V(1).ResultWriter(), "a result")
		assert.Equal(len("a result"), n)
		assert.NoError(err)
		assert.Equal("a result", out.String())
		assert.Empty(errOut.String())
		out.Reset()

		tOutput.WithValues("key", "value").Info("info message")
		assert.Empty(out.String())
		assert.Equal("info message\n", errOut.String())
		errOut.Reset()

		tOutput.V(1).WithValues("key", "value").Info("info message")
		assert.Empty(out.String())
		assert.Equal("info message    key=value\n", errOut.String())
		errOut.Reset()

		tOutput.WithValues("key", "value").Warn("warning message")
		assert.Empty(out.String())
		assert.Equal(termYellow+"warning message"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()

		tOutput.V(1).WithValues("key", "value").Warn("warning message")
		assert.Empty(out.String())
		assert.Equal(termYellow+"warning message    key=value"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()

		tOutput.WithValues("key", "value").Error(fmt.Errorf("error message"), "an error happened")
		assert.Empty(out.String())
		assert.Equal(termRed+"an error happened: error message"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()

		tOutput.V(1).WithValues("key", "value").Error(fmt.Errorf("error message"), "an error happened")
		assert.Empty(out.String())
		assert.Equal(termRed+"an error happened: error message    key=value"+termDefaultFg+"\n", errOut.String())
		errOut.Reset()
	})

	t.Run("operations", func(t *testing.T) {
		out := bytes.Buffer{}
		errOut := bytes.Buffer{}
		tOutput := output.NewInteractiveShell(&out, &errOut, 0)

		tOutput.StartOperation("working")
		time.Sleep(200 * time.Millisecond)
		tOutput.Info("a message")
		time.Sleep(200 * time.Millisecond)
		tOutput.EndOperation(true)
		tOutput.StartOperation("working")
		time.Sleep(200 * time.Millisecond)
		tOutput.Error(nil, "an error")
		time.Sleep(200 * time.Millisecond)
		tOutput.EndOperation(false)
		tOutput.StartOperation("working")
		time.Sleep(200 * time.Millisecond)
		tOutput.Info("another message")
		time.Sleep(200 * time.Millisecond)
		tOutput.EndOperationWithStatus(output.Success())
		tOutput.StartOperation("working")
		time.Sleep(200 * time.Millisecond)
		tOutput.Error(nil, "another error")
		time.Sleep(200 * time.Millisecond)
		tOutput.EndOperationWithStatus(output.Failure())
		tOutput.StartOperation("skipped")
		time.Sleep(200 * time.Millisecond)
		tOutput.Warn("some warning")
		time.Sleep(200 * time.Millisecond)
		tOutput.EndOperationWithStatus(output.Skipped())

		result := strings.TrimSuffix(errOut.String(), "\n")

		outputLines := strings.Split(result, "\r")
		assert.Greater(len(outputLines), 6)

		expectedFinalOutputLines := []string{
			termClearLine + "a message",
			" " + termGreen + "✓" + termDefaultFg + " working",
			termClearLine + termRed + "an error" + termDefaultFg,
			" " + termRed + "✗" + termDefaultFg + " working",
			termClearLine + "another message",
			" " + termGreen + "✓" + termDefaultFg + " working",
			termClearLine + termRed + "another error" + termDefaultFg,
			" " + termRed + "✗" + termDefaultFg + " working",
			termClearLine + termYellow + "some warning" + termDefaultFg,
			" " + termYellow + "∅" + termDefaultFg + " skipped",
		}
		actualFinalOutputLines := strings.Split(result, "\n")
		assert.Len(actualFinalOutputLines, len(expectedFinalOutputLines))
		for i, line := range actualFinalOutputLines {
			subLines := strings.Split(line, "\r")
			finalLine := subLines[len(subLines)-1]
			assert.Equal(expectedFinalOutputLines[i], finalLine)
		}
	})

	t.Run("operations with progress", func(t *testing.T) {
		out := bytes.Buffer{}
		errOut := bytes.Buffer{}
		gauge := &output.ProgressGauge{}
		gauge.SetStatus("a message")
		gauge.SetCapacity(10)
		tOutput := output.NewInteractiveShell(&out, &errOut, 0)

		tOutput.StartOperationWithProgress(gauge)
		gauge.Set(1)
		tOutput.Info(gauge.String())
		tOutput.EndOperation(true)

		tOutput.StartOperationWithProgress(gauge)
		gauge.Set(10)
		tOutput.Info(gauge.String())
		tOutput.EndOperation(true)

		tOutput.StartOperationWithProgress(gauge)
		gauge.Set(1)
		tOutput.Error(nil, "an error")
		tOutput.EndOperation(false)

		tOutput.StartOperationWithProgress(gauge)
		gauge.Set(10)
		tOutput.Info(gauge.String())
		tOutput.EndOperationWithStatus(output.Success())

		tOutput.StartOperationWithProgress(gauge)
		gauge.Set(1)
		tOutput.Error(nil, "another error")
		tOutput.EndOperationWithStatus(output.Failure())

		tOutput.StartOperationWithProgress(gauge)
		gauge.Set(10)
		tOutput.Info(gauge.String())
		tOutput.EndOperationWithStatus(output.Success())

		tOutput.StartOperation("without a gauge")
		tOutput.EndOperationWithStatus(output.Success())

		result := strings.TrimSuffix(errOut.String(), "\n")

		outputLines := strings.Split(result, "\r")
		assert.Greater(len(outputLines), 6)

		expectedFinalOutputLines := []string{
			termClearLine + " a message [===>                                1/10] (time elapsed 00s) ",
			" " + termGreen + "✓" + termDefaultFg + " a message [===>                                1/10] (time elapsed 00s) ",
			termClearLine + " a message [==================================>10/10] (time elapsed 00s) ",
			" " + termGreen + "✓" + termDefaultFg + " a message [==================================>10/10] (time elapsed 00s) ",
			termClearLine + termRed + "an error" + termDefaultFg,
			" " + termRed + "✗" + termDefaultFg + " a message [===>                                1/10] (time elapsed 00s) ",
			termClearLine + " a message [==================================>10/10] (time elapsed 00s) ",
			" " + termGreen + "✓" + termDefaultFg + " a message [==================================>10/10] (time elapsed 00s) ",
			termClearLine + termRed + "another error" + termDefaultFg,
			" " + termRed + "✗" + termDefaultFg + " a message [===>                                1/10] (time elapsed 00s) ",
			termClearLine + " a message [==================================>10/10] (time elapsed 00s) ",
			" " + termGreen + "✓" + termDefaultFg + " a message [==================================>10/10] (time elapsed 00s) ",
			" " + termGreen + "✓" + termDefaultFg + " without a gauge",
		}
		actualFinalOutputLines := strings.Split(result, "\n")
		assert.Len(actualFinalOutputLines, len(expectedFinalOutputLines))
		for i, line := range actualFinalOutputLines {
			subLines := strings.Split(line, "\r")
			finalLine := subLines[len(subLines)-1]
			assert.Equal(expectedFinalOutputLines[i], finalLine)
		}
	})

	t.Run("concurrent", func(t *testing.T) {
		tOutput := output.NewInteractiveShell(io.Discard, io.Discard, 0)

		wg := sync.WaitGroup{}
		doStuff := func() {
			tOutput.StartOperation("working")
			tOutput.Info("a message")
			tOutput.Warn("a warning")
			tOutput.EndOperation(true)
			tOutput.StartOperation("working")
			tOutput.Error(nil, "an error")
			tOutput.EndOperation(false)
			tOutput.Warn("another warning")
			tOutput.EndOperationWithStatus(output.Success())
			tOutput.StartOperation("working again")
			tOutput.Error(nil, "another error")
			tOutput.EndOperationWithStatus(output.Failure())
			wg.Done()
		}

		wg.Add(2)
		go doStuff()
		go doStuff()
		wg.Wait()
	})

	t.Run("message level should not change the maximum allowed verbosity", func(t *testing.T) {
		out := bytes.Buffer{}
		errOut := bytes.Buffer{}
		maxAllowedVerbosity := 1

		o := output.NewInteractiveShell(&out, &errOut, maxAllowedVerbosity)

		// Decreasing the level should not decrease the max allowed verbosity.
		o.V(maxAllowedVerbosity - 1).V(maxAllowedVerbosity).Info("test")
		assert.Equal("test\n", errOut.String())
		errOut.Reset()

		// Increasing the level should not increase the max allowed verbosity.
		o.V(maxAllowedVerbosity + 1).V(maxAllowedVerbosity + 1).Info("should not be output")
		assert.Equal("", errOut.String())
		errOut.Reset()
	})
}
