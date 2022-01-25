// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package root

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"

	"github.com/spf13/pflag"
)

// ProfilingOptions contains settings for profiling.
type ProfilingOptions struct {
	profileName   string
	profileOutput string
}

// NewProfilingOptions initializes ProfilingOptions with defaults.
func NewProfilingOptions() *ProfilingOptions {
	return &ProfilingOptions{
		profileName:   "none",
		profileOutput: "profile.pprof",
	}
}

// AddFlags adds flags for setting profiling options to the provided FlagSet.
func (o *ProfilingOptions) AddFlags(flagSet *pflag.FlagSet) {
	flagSet.StringVar(
		&o.profileName, "profile", o.profileName,
		"Name of profile to capture. One of (none|cpu|goroutine|threadcreate|heap|allocs|block|mutex)",
	)
	_ = flagSet.MarkHidden("profile")
	flagSet.StringVar(&o.profileOutput, "profile-output", o.profileOutput, "Name of the file to write the profile to")
	_ = flagSet.MarkHidden("profile-output")
}

// InitProfiling starts profiling.
func (o *ProfilingOptions) InitProfiling() error {
	switch o.profileName {
	case "none":
		return nil
	case "cpu":
		f, err := os.Create(o.profileOutput)
		if err != nil {
			return err
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			return err
		}
	// Block and mutex profiles need a call to Set{Block,Mutex}ProfileRate to
	// output anything. We choose to sample all events.
	case "block":
		runtime.SetBlockProfileRate(1)
	case "mutex":
		runtime.SetMutexProfileFraction(1)
	default:
		// Check the profile name is valid.
		if profile := pprof.Lookup(o.profileName); profile == nil {
			return fmt.Errorf("unknown profile '%s'", o.profileName)
		}
	}

	// If the command is interrupted before the end (ctrl-c), flush the
	// profiling files
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		_ = o.FlushProfiling()
		os.Exit(0)
	}()

	return nil
}

// FlushProfiling stops profiling and writes remaining unwritten data.
func (o *ProfilingOptions) FlushProfiling() error {
	switch o.profileName {
	case "none":
		return nil
	case "cpu":
		pprof.StopCPUProfile()
	case "heap":
		runtime.GC()
		fallthrough
	default:
		profile := pprof.Lookup(o.profileName)
		if profile == nil {
			return nil
		}
		f, err := os.Create(o.profileOutput)
		if err != nil {
			return err
		}
		_ = profile.WriteTo(f, 0)
	}

	return nil
}
