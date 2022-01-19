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
