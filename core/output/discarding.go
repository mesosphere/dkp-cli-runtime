// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package output

// NewDiscardingOutput returns an Output that discards all output. Useful for testing, among other purposes.
func NewDiscardingOutput() Output {
	return &discardingOutput{}
}

type discardingOutput struct{ noopOutput }

func (o *discardingOutput) WithValues(keysAndValues ...interface{}) Output { return o }
func (o *discardingOutput) V(level int) Output                             { return o }

// Convention used to verify, at compile time, that DiscardingOutput implements the Output interface.
var _ Output = &discardingOutput{}
