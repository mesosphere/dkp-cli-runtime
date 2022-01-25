// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package symlinkexechandler

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// HandleExec translates the command invoked into any necessary subcommands, re-execing
// the target executable as necessary, and returns true if the executable is re-execed, false
// otherwise.
//
// Idiomatic use is as first lines in your main:
//
// 	if symlinkexechandler.HandleExec("simple") {
// 		return
// 	}
func HandleExec(prefixToStrip string) bool {
	if !strings.HasSuffix(prefixToStrip, "-") {
		prefixToStrip += "-"
	}
	reexecedName := strings.TrimSuffix(prefixToStrip, "-")

	// Store the name used to invoke this process.
	executableInvoked := filepath.Base(os.Args[0])

	// If the executable invoked doesn't start with the specified prefix, then this is just a
	// standard symlink - return here rather than re-execing.
	if !strings.HasPrefix(executableInvoked, prefixToStrip) {
		return false
	}
	executableInvoked = strings.TrimPrefix(executableInvoked, prefixToStrip)

	// Construct the args to pass to the re-execed executable splitting on "-" and replacing
	// "_" with "-" for each subcommand.
	executableInvokedSplit := strings.Split(executableInvoked, "-")
	args := make([]string, 0, len(executableInvokedSplit)+len(os.Args[1:]))
	for _, s := range executableInvokedSplit {
		args = append(args, strings.ReplaceAll(s, "_", "-"))
	}
	// Append all the other args specified too.
	args = append(args, os.Args[1:]...)

	// Get the actual executable used to start this process.
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// Executable could be a symlink, so resolve to the actual file if necessary.
	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		panic(err)
	}

	// Re-exec using the actual executable, passing subcommands derived from
	// the name of the executable ivoked and flags/arguments passed to the original invocation.
	if err := syscall.Exec(executable, append([]string{reexecedName}, args...), os.Environ()); err != nil {
		panic(err)
	}

	return true
}
