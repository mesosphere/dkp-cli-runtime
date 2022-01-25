// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"strings"

	symlinkexechandler "github.com/mesosphere/dkp-cli-runtime/symlink-reexec"
)

func main() {
	if symlinkexechandler.HandleExec("simple") {
		return
	}

	fmt.Print(strings.Trim(fmt.Sprintf("%q", os.Args), "[]"))
}
