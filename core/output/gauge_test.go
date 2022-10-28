// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package output_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mesosphere/dkp-cli-runtime/core/output"
)

func TestProgressGauge(t *testing.T) {
	gauge := output.ProgressGauge{}
	gauge.SetStatus("static-status")
	assert.Equal(t, " static-status", gauge.String())
	lengthShouldBeConsistent := 77
	gauge.SetCapacity(10)
	assert.Equal(t,
		" static-status [                                    0/10] (time elapsed 00s) ",
		gauge.String())
	gauge.Inc()
	assert.Equal(t,
		" static-status [===>                                1/10] (time elapsed 00s) ",
		gauge.String())
	assert.Equal(t, lengthShouldBeConsistent, len(gauge.String()))
	gauge.Inc()
	gauge.Inc()
	gauge.Dec()
	gauge.Inc()
	time.Sleep(1 * time.Second)
	assert.Equal(t,
		" static-status [==========>                         3/10] (time elapsed 01s) ",
		gauge.String())
	assert.Equal(t, lengthShouldBeConsistent, len(gauge.String()))
	gauge.Set(10)
	assert.Equal(t,
		" static-status [==================================>10/10] (time elapsed 01s) ",
		gauge.String())
	assert.Equal(t, lengthShouldBeConsistent, len(gauge.String()))
	gauge.Set(20)
	// if the gauge is incremented/decremented incorrectly, we default to static status.
	// it is the callers responsibility to make sure all the values are correct.
	assert.Equal(t, " static-status", gauge.String())
	gauge.Set(-10)
	assert.Equal(t, " static-status", gauge.String())
}

func Test_humanReadableDuration(t *testing.T) {
	assert.Equal(t, "00s", output.HumanReadableDuration(10*time.Millisecond))
	assert.Equal(t, "01s", output.HumanReadableDuration(1000*time.Millisecond))
	assert.Equal(t, "1m01s", output.HumanReadableDuration(61*1000*time.Millisecond))
	assert.Equal(t, "60m00s", output.HumanReadableDuration(60*60*1000*time.Millisecond))
	assert.Equal(t, "61m30s", output.HumanReadableDuration(61*60.5*1000*time.Millisecond))
}
