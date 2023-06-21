// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package output

import (
	"fmt"
	"io"

	"github.com/jwalton/gchalk"
)

type EndOperationStatus interface {
	Fprintln(w io.Writer, format string, a ...any) (n int, err error)
}

type status struct {
	statusCharacter string
	color           *gchalk.Builder
}

func (s status) Fprintln(w io.Writer, format string, a ...any) (n int, err error) {
	format = fmt.Sprintf(" %s %s\n", s.color.Sprintf(s.statusCharacter), format)
	return fmt.Fprintf(w, format, a...)
}

func NewStatus(statusCharacter string, color *gchalk.Builder) EndOperationStatus {
	return status{
		statusCharacter: statusCharacter,
		color:           color,
	}
}

func Success() EndOperationStatus {
	return NewStatus("✓", gchalk.Stderr.WithGreen())
}

func Failure() EndOperationStatus {
	return NewStatus("✗", gchalk.Stderr.WithRed())
}

func Skipped() EndOperationStatus {
	return NewStatus("∅", gchalk.Stderr.WithYellow())
}
