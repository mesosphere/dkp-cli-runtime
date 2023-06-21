// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package output_test

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mesosphere/dkp-cli-runtime/core/output"
)

func TestOutputLogr(t *testing.T) {
	assert := assert.New(t)

	o := &outputMock{}
	logger := output.NewOutputLogr(o)

	logger.Info("info message")
	assert.Equal("info message", o.msg)
	assert.Nil(o.err)
	assert.Equal(0, o.verbosity)
	o.Reset()

	logger.Info("info message", "key", "value")
	assert.Equal("info message", o.msg)
	assert.Equal([]interface{}{"key", "value"}, o.keysAndValues)
	assert.Equal(0, o.verbosity)
	o.Reset()

	err := errors.New("an error")
	logger.Error(err, "error message")
	assert.Equal("error message", o.msg)
	assert.Equal(err, o.err)
	assert.Equal(0, o.verbosity)
	o.Reset()

	logger.Error(err, "error message", "key", "value")
	assert.Equal("error message", o.msg)
	assert.Equal(err, o.err)
	assert.Equal([]interface{}{"key", "value"}, o.keysAndValues)
	assert.Equal(0, o.verbosity)
	o.Reset()

	logger.WithValues("key", "value").WithValues("key2", "value2").Info("info message", "key3", "value3")
	assert.Equal("info message", o.msg)
	assert.Equal([]interface{}{
		"key", "value",
		"key2", "value2",
		"key3", "value3",
	}, o.keysAndValues)
	assert.Equal(0, o.verbosity)
	o.Reset()

	logger.V(1).Info("info message")
	assert.Equal("info message", o.msg)
	assert.Nil(o.err)
	assert.Equal(1, o.verbosity)
	o.Reset()

	logger.V(1).V(2).Info("info message")
	assert.Equal("info message", o.msg)
	assert.Nil(o.err)
	assert.Equal(3, o.verbosity)
	o.Reset()
}

type outputMock struct {
	msg           string
	err           error
	keysAndValues []interface{}
	verbosity     int
}

func (o *outputMock) Reset() {
	o.msg = ""
	o.err = nil
	o.keysAndValues = []interface{}{}
	o.verbosity = 0
}

func (o *outputMock) Info(msg string) {
	o.msg = msg
}

func (o *outputMock) Error(err error, msg string) {
	o.err = err
	o.msg = msg
}

func (o *outputMock) V(level int) output.Output {
	o.verbosity = level
	return o
}

func (o *outputMock) WithValues(keysAndValues ...interface{}) output.Output {
	o.keysAndValues = append(o.keysAndValues, keysAndValues...)
	return o
}

func (o *outputMock) Infof(format string, args ...interface{})                {}
func (o *outputMock) InfoWriter() io.Writer                                   { return io.Discard }
func (o *outputMock) Warn(msg string)                                         {}
func (o *outputMock) Warnf(format string, args ...interface{})                {}
func (o *outputMock) WarnWriter() io.Writer                                   { return io.Discard }
func (o *outputMock) Errorf(err error, format string, args ...interface{})    {}
func (o *outputMock) ErrorWriter() io.Writer                                  { return io.Discard }
func (o *outputMock) StartOperation(status string)                            {}
func (o *outputMock) StartOperationWithProgress(gauge *output.ProgressGauge)  {}
func (o *outputMock) EndOperation(success bool)                               {}
func (o *outputMock) EndOperationWithStatus(status output.EndOperationStatus) {}
func (o *outputMock) Result(result string)                                    {}
func (o *outputMock) ResultWriter() io.Writer                                 { return io.Discard }
func (o *outputMock) Enabled(level int) bool                                  { return true }
