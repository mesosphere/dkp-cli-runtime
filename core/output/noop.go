// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package output

import "io"

// DiscardingOutput discards all output, which can be useful for testing, among other purposes.
type DiscardingOutput struct{}

// Convention used to verify, at compile time, that DiscardingOutput implements the Output interface.
var _ Output = DiscardingOutput{}

func (o DiscardingOutput) Info(msg string)                                      {}
func (o DiscardingOutput) Infof(format string, args ...interface{})             {}
func (o DiscardingOutput) InfoWriter() io.Writer                                { return io.Discard }
func (o DiscardingOutput) Warn(msg string)                                      {}
func (o DiscardingOutput) Warnf(format string, args ...interface{})             {}
func (o DiscardingOutput) WarnWriter() io.Writer                                { return io.Discard }
func (o DiscardingOutput) Error(err error, msg string)                          {}
func (o DiscardingOutput) Errorf(err error, format string, args ...interface{}) {}
func (o DiscardingOutput) ErrorWriter() io.Writer                               { return io.Discard }
func (o DiscardingOutput) StartOperation(status string)                         {}
func (o DiscardingOutput) EndOperation(success bool)                            {}
func (o DiscardingOutput) Result(result string)                                 {}
func (o DiscardingOutput) ResultWriter() io.Writer                              { return io.Discard }
func (o DiscardingOutput) WithValues(keysAndValues ...interface{}) Output       { return o }
func (o DiscardingOutput) V(level int) Output                                   { return o }
