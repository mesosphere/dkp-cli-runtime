package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/mesosphere/dkp-cli-runtime/core/cmd/root"
	"github.com/mesosphere/dkp-cli-runtime/core/output"
	"github.com/mesosphere/dkp-cli-runtime/extensions/cmd/get"
	"github.com/mesosphere/dkp-cli-runtime/extensions/options"
)

func NewCommand(in io.Reader, out, errOut io.Writer) (*cobra.Command, output.Output) {
	rootCmd, rootOpts := root.NewCommand(out, errOut)

	clientOpts := options.NewClientOptions(true)
	clientOpts.AddFlags(rootCmd.PersistentFlags())

	ioStreams := genericclioptions.IOStreams{In: in, Out: out, ErrOut: errOut}
	rootCmd.AddCommand(get.NewCommand(ioStreams, clientOpts, "pods"))

	return rootCmd, rootOpts.Output
}

func Execute() {
	rootCmd, out := NewCommand(os.Stdin, os.Stdout, os.Stderr)
	rootCmd.SilenceErrors = true

	if err := rootCmd.Execute(); err != nil {
		out.Error(err, "")
		os.Exit(1)
	}
}
