// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package root

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"

	"github.com/mesosphere/dkp-cli-runtime/core/cmd/help"
	"github.com/mesosphere/dkp-cli-runtime/core/cmd/version"
	"github.com/mesosphere/dkp-cli-runtime/core/output"
	"github.com/mesosphere/dkp-cli-runtime/core/plugin"
	"github.com/mesosphere/dkp-cli-runtime/core/term"
)

// RootOptions contains options configured in the root command.
type RootOptions struct {
	Profiling *ProfilingOptions
	Output    output.Output
}

// NewCommand creates a root command with useful built-in features like:
// - profiling
// - version command with different output formats
// - help command with different output formats
// - command discovery for use as a CLI plugin.
func NewCommand(out, errOut io.Writer) (*cobra.Command, *RootOptions) {
	profilingOpts := NewProfilingOptions()

	rootCmd := &cobra.Command{
		Use:          filepath.Base(os.Args[0]),
		Args:         cobra.NoArgs,
		SilenceUsage: true,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := profilingOpts.InitProfiling(); err != nil {
				return err
			}
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return profilingOpts.FlushProfiling()
		},
	}

	profilingOpts.AddFlags(rootCmd.PersistentFlags())
	outputVerbosity := 0
	klogVmodule := ""
	rootCmd.PersistentFlags().IntVarP(&outputVerbosity, "verbose", "v", 0, "Output verbosity")
	rootCmd.PersistentFlags().StringVar(&klogVmodule, "vmodule", "",
		"Comma-separated list of pattern=N settings for file-filtered logging")
	rootCmd.PersistentFlags().MarkHidden("vmodule") //nolint:errcheck // flag just created, guaranteed to succeed
	ensureTitleCaseForHelpFlagUsage(rootCmd)

	rootCmd.AddCommand(version.NewCommand(out))
	rootCmd.AddCommand(plugin.NewDiscoveryCommand(out, rootCmd))
	rootCmd.SetHelpCommand(help.NewHelpCommandWrapper(rootCmd))

	// Make sure flags are parsed, ignoring unknown flags at this stage. This ensures that the
	// logging flags are initialized. The flags will be parsed again when the command is run, at
	// which point unknown flags will trigger an error.
	origParseErrorsWhitelist := rootCmd.PersistentFlags().ParseErrorsWhitelist
	rootCmd.PersistentFlags().ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{UnknownFlags: true}
	_ = rootCmd.PersistentFlags().Parse(os.Args)
	rootCmd.PersistentFlags().ParseErrorsWhitelist = origParseErrorsWhitelist

	verbosityFlagSet := rootCmd.PersistentFlags().Changed("verbose")

	rootOpts := &RootOptions{
		Profiling: profilingOpts,
		Output:    configureOutput(out, errOut, outputVerbosity, verbosityFlagSet, klogVmodule),
	}
	return rootCmd, rootOpts
}

func configureOutput(out, errOut io.Writer, verbosity int, verbosityFlagSet bool, klogVmodule string) output.Output {
	o := newOutput(out, errOut, verbosity)

	// send output of standard logger to Info, verbosity 1
	log.SetFlags(0)
	log.SetOutput(o.V(1).InfoWriter())

	// send klog logs to output if verbosity flag is set
	if verbosityFlagSet || klogVmodule != "" {
		o := newOutput(out, errOut, math.MaxInt)
		configureKlog(o, verbosity, klogVmodule)
	} else {
		klog.SetLogger(logr.Discard())
	}

	return o
}

func newOutput(out, errOut io.Writer, verbosity int) output.Output {
	if term.IsSmartTerminal(errOut) {
		return output.NewInteractiveShell(out, errOut, verbosity)
	} else {
		return output.NewNonInteractiveShell(out, errOut, verbosity)
	}
}

func configureKlog(o output.Output, verbosity int, vModule string) {
	klogFlags := flag.NewFlagSet("klog", flag.ContinueOnError)
	klogFlags.Usage = func() {}
	klogFlags.SetOutput(o.ErrorWriter())
	klog.InitFlags(klogFlags)
	klogFlags.Parse([]string{
		"--v", fmt.Sprint(verbosity),
		"--vmodule", vModule,
	})
	klog.SetLogger(output.NewOutputLogr(o))
}

func ensureTitleCaseForHelpFlagUsage(rootCmd *cobra.Command) {
	rootCmdHelpFunc := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, s []string) {
		if helpFlag := cmd.Flags().Lookup("help"); helpFlag != nil {
			helpFlag.Usage = "Help for " + cmd.Name()
		}
		rootCmdHelpFunc(cmd, s)
	})
}
