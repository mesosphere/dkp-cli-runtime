// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package output

import "io"

// noopOutput discards the output at a specific message level, but preserves the desired verbosity, by delegating V() to
// the embedded Output.
type noopOutput struct {
	Output
}

func (o noopOutput) Info(msg string)                                      {}
func (o noopOutput) Infof(format string, args ...interface{})             {}
func (o noopOutput) InfoWriter() io.Writer                                { return io.Discard }
func (o noopOutput) Warn(msg string)                                      {}
func (o noopOutput) Warnf(format string, args ...interface{})             {}
func (o noopOutput) WarnWriter() io.Writer                                { return io.Discard }
func (o noopOutput) Error(err error, msg string)                          {}
func (o noopOutput) Errorf(err error, format string, args ...interface{}) {}
func (o noopOutput) ErrorWriter() io.Writer                               { return io.Discard }
func (o noopOutput) StartOperation(status string)                         {}
func (o noopOutput) EndOperation(success bool)                            {}
func (o noopOutput) Result(result string)                                 {}
func (o noopOutput) ResultWriter() io.Writer                              { return io.Discard }
func (o noopOutput) WithValues(keysAndValues ...interface{}) Output       { return o }
