// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package output

import "io"

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

// DiscardingOutput discards all output, which can be useful for testing, among other purposes.
type DiscardingOutput struct{ noopOutput }

func (o DiscardingOutput) WithValues(keysAndValues ...interface{}) Output { return o }
func (o DiscardingOutput) V(level int) Output                             { return o }

// Convention used to verify, at compile time, that DiscardingOutput implements the Output interface.
var _ Output = DiscardingOutput{}
