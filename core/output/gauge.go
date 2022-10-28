// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package output

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const maxProgressBarWidth = 40

// ProgressGauge is not really a Gauge in the truest sense and is used only to display a progress bar.
// It is a gauge in the sense that the value can be incremented or decremented.
// The correlation with word gauge ends here.
//
// This is used to display a progress bar that looks like:
//
//	static-status [====>                                    1/10] (time elapsed 00s)
//	static-status [============>                            3/10] (time elapsed 00s)
//	static-status [=======================================>10/10] (time elapsed 00s)
type ProgressGauge struct {
	status    string
	current   int
	capacity  int
	startTime time.Time
}

func (g *ProgressGauge) IsReady() bool {
	if g == nil {
		return false
	}
	if g.current < 0 {
		return false
	}
	if g.capacity <= 0 {
		return false
	}
	if g.current > g.capacity {
		return false
	}
	if g.startTime.IsZero() {
		return false
	}
	return true
}

func (g *ProgressGauge) SetCapacity(capacity int) {
	if g == nil {
		return
	}
	if capacity < 0 {
		return
	}
	g.capacity = capacity
}

func (g *ProgressGauge) SetStatus(status string) {
	if g == nil {
		return
	}
	g.status = status
}

func (g *ProgressGauge) Set(current int) {
	if g == nil {
		return
	}
	g.current = current
}

func (g *ProgressGauge) Inc() {
	if g == nil {
		return
	}
	g.current += 1
}

func (g *ProgressGauge) Dec() {
	if g == nil {
		return
	}
	g.current -= 1
}

func (g *ProgressGauge) InitStartTime() {
	if g == nil {
		return
	}
	if g.startTime.IsZero() {
		g.startTime = time.Now()
	}
}

// String generates a string representation of the progress based on the current and capacity values of the gauge.
// It ensures that the progress bar generated is of fixed length format
// It also appends the elapsed time to string representation (if timer is not set, this will initialize it).
func (g *ProgressGauge) String() string {
	if g == nil {
		return ""
	}
	if g.startTime.IsZero() {
		g.startTime = time.Now()
	}
	if !g.IsReady() {
		return fmt.Sprintf(" %s", g.status)
	}
	duration := HumanReadableDuration(time.Since(g.startTime))
	ratio := fmt.Sprintf("%d/%d", g.current, g.capacity)
	availableSpace := maxProgressBarWidth - len(ratio)
	progress := (availableSpace * g.current) / g.capacity
	if progress < 0 {
		progress = 0
	}
	progressStr := ""
	if progress > 0 {
		backtrack := 0
		if progress == availableSpace {
			backtrack--
		}
		progressStr = fmt.Sprintf("%s>", strings.Repeat("=", progress+backtrack))
	}
	spaces := availableSpace - len(progressStr)
	if spaces < 0 {
		spaces = 0
	}
	return fmt.Sprintf(" %s [%s%s%s] (time elapsed %s) ",
		g.status,
		progressStr,
		strings.Repeat(" ", spaces),
		ratio,
		duration)
}

// HumanReadableDuration converts duration to a human-readable format like:
//
//	00s
//	30s
//	1m00s
//	1m30s
//	61m30s (we never use hours)
func HumanReadableDuration(duration time.Duration) string {
	minutes := int(duration.Minutes())
	seconds := int(math.Mod(duration.Seconds(), 60))
	durationStr := ""
	if minutes > 0 {
		durationStr = fmt.Sprintf("%dm", minutes)
	}
	durationStr = fmt.Sprintf("%s%02ds", durationStr, seconds)
	return durationStr
}
