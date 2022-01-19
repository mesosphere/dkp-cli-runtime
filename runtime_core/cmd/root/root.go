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
	"k8s.io/klog/v2"

	"github.com/mesosphere/dkp-cli-runtime/runtime_core/cmd/help"
	"github.com/mesosphere/dkp-cli-runtime/runtime_core/cmd/version"
	"github.com/mesosphere/dkp-cli-runtime/runtime_core/output"
	"github.com/mesosphere/dkp-cli-runtime/runtime_core/plugin"
	"github.com/mesosphere/dkp-cli-runtime/runtime_core/term"
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
	outputVerbosity := -1
	klogVmodule := ""
	rootCmd.PersistentFlags().IntVarP(&outputVerbosity, "verbose", "v", -1, "Output verbosity")
	rootCmd.PersistentFlags().StringVar(&klogVmodule, "vmodule", "",
		"Comma-separated list of pattern=N settings for file-filtered logging")
	rootCmd.PersistentFlags().MarkHidden("vmodule") // nolint:errcheck // flag just created, guaranteed to succeed
	ensureTitleCaseForHelpFlagUsage(rootCmd)

	rootCmd.AddCommand(version.NewCommand(out))
	rootCmd.AddCommand(plugin.NewDiscoveryCommand(out, rootCmd))
	rootCmd.SetHelpCommand(help.NewHelpCommandWrapper(rootCmd))

	// make sure flags are parsed
	// nolint:errcheck // will be parsed again by cobra
	rootCmd.PersistentFlags().Parse(os.Args)

	rootOpts := &RootOptions{
		Profiling: profilingOpts,
		Output:    configureOutput(out, errOut, outputVerbosity, klogVmodule),
	}
	return rootCmd, rootOpts
}

func configureOutput(out, errOut io.Writer, verbosity int, klogVmodule string) output.Output {
	o := newOutput(out, errOut, verbosity)

	// send output of standard logger to Info, verbosity 1
	log.SetFlags(0)
	log.SetOutput(o.V(1).InfoWriter())

	// send klog logs to output if verbosity flag is set
	if verbosity >= 0 || klogVmodule != "" {
		o := newOutput(out, errOut, math.MaxInt)
		configureKlog(o, verbosity, klogVmodule)
	} else {
		klog.SetLogger(logr.Discard())
	}

	return o
}

func newOutput(out, errOut io.Writer, verbosity int) output.Output {
	if verbosity < 0 {
		verbosity = 0
	}
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
