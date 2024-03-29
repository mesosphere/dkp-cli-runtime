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

	"github.com/stretchr/testify/assert"

	"github.com/mesosphere/dkp-cli-runtime/core/output"
)

func TestNonInteractiveShellOutput(t *testing.T) {
	assert := assert.New(t)

	assertEqualExceptTimestamp := func(expected, actual string) {
		t.Helper()
		firstSpace := strings.Index(actual, " ")
		if firstSpace != -1 {
			actual = actual[firstSpace+1:]
			secondSpace := strings.Index(actual, " ")
			if secondSpace != -1 {
				actual = "<timestamp>" + actual[secondSpace:]
			}
		}
		assert.Equal(expected, actual)
	}

	t.Run("default", func(t *testing.T) {
		out := bytes.Buffer{}
		errOut := bytes.Buffer{}
		output := output.NewNonInteractiveShell(&out, &errOut, 0)

		output.Info("info message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> INF info message\n", errOut.String())
		errOut.Reset()

		output.Infof("info %s", "message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> INF info message\n", errOut.String())
		errOut.Reset()

		n, err := io.WriteString(output.InfoWriter(), "info message")
		assert.Equal(len("info message"), n)
		assert.NoError(err)
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> INF info message\n", errOut.String())
		errOut.Reset()

		output.Warn("warning message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> WRN warning message\n", errOut.String())
		errOut.Reset()

		output.Warnf("warning %s", "message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> WRN warning message\n", errOut.String())
		errOut.Reset()

		n, err = io.WriteString(output.WarnWriter(), "warning message")
		assert.Equal(len("warning message"), n)
		assert.NoError(err)
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> WRN warning message\n", errOut.String())
		errOut.Reset()

		output.Error(fmt.Errorf("error message"), "an error happened")
		assert.Empty(out.String())
		assertEqualExceptTimestamp(`<timestamp> ERR an error happened    err="error message"`+"\n", errOut.String())
		errOut.Reset()
		output.Error(nil, "an error happened")
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> ERR an error happened    err=<nil>\n", errOut.String())
		errOut.Reset()
		errOut.Reset()
		output.Error(fmt.Errorf("error message"), "")
		assert.Empty(out.String())
		assertEqualExceptTimestamp(`<timestamp> ERR     err="error message"`+"\n", errOut.String())
		errOut.Reset()

		output.Errorf(fmt.Errorf("error message"), "an error %s", "happened")
		assert.Empty(out.String())
		assertEqualExceptTimestamp(`<timestamp> ERR an error happened    err="error message"`+"\n", errOut.String())
		errOut.Reset()

		n, err = io.WriteString(output.ErrorWriter(), "an error happened")
		assert.Equal(len("an error happened"), n)
		assert.NoError(err)
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> ERR an error happened    err=<nil>\n", errOut.String())
		errOut.Reset()

		output.Result("a result")
		assert.Equal("a result\n", out.String())
		assert.Empty(errOut.String())
		out.Reset()

		n, err = io.WriteString(output.ResultWriter(), "a result")
		assert.Equal(len("a result"), n)
		assert.NoError(err)
		assert.Equal("a result", out.String())
		assert.Empty(errOut.String())
		out.Reset()

		output.WithValues("key", "value").Info("info message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> INF info message    key=value\n", errOut.String())
		errOut.Reset()

		output.WithValues("key", "value with spaces").Info("info message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp(`<timestamp> INF info message    key="value with spaces"`+"\n", errOut.String())
		errOut.Reset()

		output.WithValues("key", "value=withequalsign").Info("info message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp(`<timestamp> INF info message    key="value=withequalsign"`+"\n", errOut.String())
		errOut.Reset()

		output.WithValues("key", "value").Warn("warning message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> WRN warning message    key=value\n", errOut.String())
		errOut.Reset()

		output.WithValues("key", "value with spaces").Warn("warning message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp(`<timestamp> WRN warning message    key="value with spaces"`+"\n", errOut.String())
		errOut.Reset()

		output.WithValues("key", "value=withequalsign").Warn("warning message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp(`<timestamp> WRN warning message    key="value=withequalsign"`+"\n", errOut.String())
		errOut.Reset()
	})

	t.Run("verbosity hidden", func(t *testing.T) {
		out := bytes.Buffer{}
		errOut := bytes.Buffer{}
		output := output.NewNonInteractiveShell(&out, &errOut, 0)

		output.V(1).Info("info message")
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		output.V(1).Infof("info %s", "message")
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		_, err := io.WriteString(output.V(1).InfoWriter(), "info message")
		assert.NoError(err)
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		output.V(1).Warn("warning message")
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		output.V(1).Warnf("warning %s", "message")
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		_, err = io.WriteString(output.V(1).WarnWriter(), "warning message")
		assert.NoError(err)
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		output.V(1).Error(fmt.Errorf("error message"), "an error happened")
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		output.V(1).Errorf(fmt.Errorf("error message"), "an error %s", "happened")
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		_, err = io.WriteString(output.V(1).ErrorWriter(), "an error happened")
		assert.NoError(err)
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		output.V(1).Result("a result")
		assert.Empty(out.String())
		assert.Equal("", out.String())
		assert.Empty(errOut.String())
		out.Reset()

		_, err = io.WriteString(output.V(1).ResultWriter(), "a result")
		assert.NoError(err)
		assert.Equal("", out.String())
		assert.Empty(errOut.String())
		out.Reset()

		output.V(1).WithValues("key", "value").Info("info message")
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		output.V(1).WithValues("key", "value").Warn("warning message")
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()

		output.V(1).StartOperation("working")
		output.V(1).EndOperation(true)
		assert.Empty(out.String())
		assert.Equal("", errOut.String())
		errOut.Reset()
	})

	t.Run("verbosity", func(t *testing.T) {
		out := bytes.Buffer{}
		errOut := bytes.Buffer{}
		output := output.NewNonInteractiveShell(&out, &errOut, 1)

		output.V(1).Info("info message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> INF info message\n", errOut.String())
		errOut.Reset()

		output.V(1).Infof("info %s", "message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> INF info message\n", errOut.String())
		errOut.Reset()

		n, err := io.WriteString(output.V(1).InfoWriter(), "info message")
		assert.Equal(len("info message"), n)
		assert.NoError(err)
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> INF info message\n", errOut.String())
		errOut.Reset()

		output.V(1).Warn("warning message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> WRN warning message\n", errOut.String())
		errOut.Reset()

		output.V(1).Warnf("warning %s", "message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> WRN warning message\n", errOut.String())
		errOut.Reset()

		n, err = io.WriteString(output.V(1).WarnWriter(), "warning message")
		assert.Equal(len("warning message"), n)
		assert.NoError(err)
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> WRN warning message\n", errOut.String())
		errOut.Reset()

		output.V(1).Error(fmt.Errorf("error message"), "an error happened")
		assert.Empty(out.String())
		assertEqualExceptTimestamp(`<timestamp> ERR an error happened    err="error message"`+"\n", errOut.String())
		errOut.Reset()
		output.V(1).Error(nil, "an error happened")
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> ERR an error happened    err=<nil>\n", errOut.String())
		errOut.Reset()
		errOut.Reset()
		output.V(1).Error(fmt.Errorf("error message"), "")
		assert.Empty(out.String())
		assertEqualExceptTimestamp(`<timestamp> ERR     err="error message"`+"\n", errOut.String())
		errOut.Reset()

		output.V(1).Errorf(fmt.Errorf("error message"), "an error %s", "happened")
		assert.Empty(out.String())
		assertEqualExceptTimestamp(`<timestamp> ERR an error happened    err="error message"`+"\n", errOut.String())
		errOut.Reset()

		n, err = io.WriteString(output.V(1).ErrorWriter(), "an error happened")
		assert.Equal(len("an error happened"), n)
		assert.NoError(err)
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> ERR an error happened    err=<nil>\n", errOut.String())
		errOut.Reset()

		output.V(1).Result("a result")
		assert.Equal("a result\n", out.String())
		assert.Empty(errOut.String())
		out.Reset()

		n, err = io.WriteString(output.V(1).ResultWriter(), "a result")
		assert.Equal(len("a result"), n)
		assert.NoError(err)
		assert.Equal("a result", out.String())
		assert.Empty(errOut.String())
		out.Reset()

		output.V(1).WithValues("key", "value").Info("info message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> INF info message    key=value\n", errOut.String())
		errOut.Reset()

		output.V(1).WithValues("key", "value with spaces").Info("info message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp(`<timestamp> INF info message    key="value with spaces"`+"\n", errOut.String())
		errOut.Reset()

		output.V(1).WithValues("key", "value=withequalsign").Info("info message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp(`<timestamp> INF info message    key="value=withequalsign"`+"\n", errOut.String())
		errOut.Reset()

		output.WithValues("key", "value=withequalsign").V(1).Info("info message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp(`<timestamp> INF info message    key="value=withequalsign"`+"\n", errOut.String())
		errOut.Reset()

		output.V(1).WithValues("key", "value").Warn("warning message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp("<timestamp> WRN warning message    key=value\n", errOut.String())
		errOut.Reset()

		output.V(1).WithValues("key", "value with spaces").Warn("warning message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp(`<timestamp> WRN warning message    key="value with spaces"`+"\n", errOut.String())
		errOut.Reset()

		output.V(1).WithValues("key", "value=withequalsign").Warn("warning message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp(`<timestamp> WRN warning message    key="value=withequalsign"`+"\n", errOut.String())
		errOut.Reset()

		output.WithValues("key", "value=withequalsign").V(1).Warn("warning message")
		assert.Empty(out.String())
		assertEqualExceptTimestamp(`<timestamp> WRN warning message    key="value=withequalsign"`+"\n", errOut.String())
		errOut.Reset()
	})

	t.Run("operations", func(t *testing.T) {
		out := bytes.Buffer{}
		errOut := bytes.Buffer{}
		output := output.NewNonInteractiveShell(&out, &errOut, 0)

		output.StartOperation("working")
		time.Sleep(200 * time.Millisecond)
		output.Info("a message")
		time.Sleep(200 * time.Millisecond)
		output.EndOperation(true)
		output.StartOperation("working")
		time.Sleep(200 * time.Millisecond)
		output.Error(nil, "an error")
		time.Sleep(200 * time.Millisecond)
		output.EndOperation(false)

		result := strings.TrimSuffix(errOut.String(), "\n")
		outputLines := strings.Split(result, "\n")
		assert.Len(outputLines, 6)

		assertEqualExceptTimestamp("<timestamp> INF  • working...", outputLines[0])
		assertEqualExceptTimestamp("<timestamp> INF a message", outputLines[1])
		assertEqualExceptTimestamp("<timestamp> INF  ✓ working", outputLines[2])
		assertEqualExceptTimestamp("<timestamp> INF  • working...", outputLines[3])
		assertEqualExceptTimestamp("<timestamp> ERR an error    err=<nil>", outputLines[4])
		assertEqualExceptTimestamp("<timestamp> INF  ✗ working", outputLines[5])
	})

	t.Run("concurrent", func(t *testing.T) {
		output := output.NewNonInteractiveShell(io.Discard, io.Discard, 0)

		wg := sync.WaitGroup{}
		doStuff := func() {
			output.StartOperation("working")
			output.Info("a message")
			output.EndOperation(true)
			output.StartOperation("working")
			output.Error(nil, "an error")
			output.EndOperation(false)
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

		o := output.NewNonInteractiveShell(&out, &errOut, maxAllowedVerbosity)

		// Decreasing the level should not decrease the max allowed verbosity.
		o.V(maxAllowedVerbosity - 1).V(maxAllowedVerbosity).Info("test")
		assertEqualExceptTimestamp("<timestamp> INF test\n", errOut.String())
		errOut.Reset()

		// Increasing the level should not increase the max allowed verbosity.
		o.V(maxAllowedVerbosity + 1).V(maxAllowedVerbosity + 1).Info("should not be output")
		assert.Equal("", errOut.String())
		errOut.Reset()
	})
}
