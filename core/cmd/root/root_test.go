// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package root_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"

	"github.com/mesosphere/dkp-cli-runtime/core/cmd/root"
)

func TestNewCommand(t *testing.T) {
	assert := assert.New(t)

	rootCmd, rootOptions := root.NewCommand(io.Discard, io.Discard)

	// all
	assert.ElementsMatch([]string{"version", "_plugin_commands"}, commandNames(rootCmd.Commands(), false))
	assert.ElementsMatch(
		[]string{"profile", "profile-output", "verbose", "v", "vmodule"},
		flagNames(rootCmd.PersistentFlags(), false),
	)

	// visible
	assert.ElementsMatch([]string{"version"}, commandNames(rootCmd.Commands(), true))
	assert.ElementsMatch([]string{"verbose", "v"}, flagNames(rootCmd.PersistentFlags(), true))

	assert.NotNil(rootOptions.Profiling)
	assert.NotNil(rootOptions.Output)
}

func commandNames(commands []*cobra.Command, onlyVisible bool) []string {
	result := make([]string, 0, len(commands))
	for _, command := range commands {
		if onlyVisible && command.Hidden {
			continue
		}
		result = append(result, command.Name())
	}
	return result
}

func flagNames(flags *pflag.FlagSet, onlyVisible bool) []string {
	names := []string{}
	flags.VisitAll(func(flag *pflag.Flag) {
		if onlyVisible && flag.Hidden {
			return
		}
		names = append(names, flag.Name)
		if flag.Shorthand != "" {
			names = append(names, flag.Shorthand)
		}
	})
	return names
}

func TestEnsureTitleCaseForHelpFlagUsage(t *testing.T) {
	assert := assert.New(t)

	output := bytes.Buffer{}
	rootCmd, _ := root.NewCommand(&output, &output)
	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)

	subCmd := &cobra.Command{
		Use: "subcommand",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	rootCmd.AddCommand(subCmd)

	output.Reset()
	rootCmd.SetArgs([]string{"--help"})
	rootCmd.Execute()
	assert.Regexp("-h, --help\\W+Help for root.test\n", output.String())

	output.Reset()
	rootCmd.SetArgs([]string{"help"})
	rootCmd.Execute()
	assert.Regexp("-h, --help\\W+Help for root.test\n", output.String())

	output.Reset()
	rootCmd.SetArgs([]string{"subcommand", "--help"})
	rootCmd.Execute()
	assert.Regexp("-h, --help\\W+Help for subcommand\n", output.String())

	output.Reset()
	rootCmd.SetArgs([]string{"help", "subcommand"})
	rootCmd.Execute()
	assert.Regexp("-h, --help\\W+Help for subcommand\n", output.String())
}

func TestVerbosity(t *testing.T) {
	for _, test := range []struct {
		name              string
		flags             []string
		regexesToMatch    []string
		regexesToNotMatch []string
	}{{
		name:              "no flags",
		regexesToMatch:    []string{"\\W+a\n"},
		regexesToNotMatch: []string{"\\W+b\n"},
	}, {
		name:           "verbose only flag",
		flags:          []string{"--verbose", "4"},
		regexesToMatch: []string{"\\W+a\n", "\\W+b\n"},
	}, {
		name:           "verbose with arbitrary flag after",
		flags:          []string{"--verbose", "4", "--wibble"},
		regexesToMatch: []string{"\\W+a\n", "\\W+b\n"},
	}, {
		name:           "verbose with arbitrary flag before",
		flags:          []string{"--wibble", "--verbose", "4"},
		regexesToMatch: []string{"\\W+a\n", "\\W+b\n"},
	}} {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			origOSArgs := os.Args
			defer func() { os.Args = origOSArgs }()

			output := bytes.Buffer{}

			os.Args = append([]string{"root"}, test.flags...)
			rootCmd, rootOpts := root.NewCommand(&output, &output)
			b := false
			rootCmd.Flags().BoolVar(&b, "wibble", false, "unused flag just for testing")
			rootCmd.SetOut(&output)
			rootCmd.SetErr(&output)
			rootCmd.Run = func(cmd *cobra.Command, args []string) {
				rootOpts.Output.Info("a")
				rootOpts.Output.V(4).Info("b")
			}

			rootCmd.Execute()

			for _, r := range test.regexesToMatch {
				assert.Regexp(r, output.String())
			}
			for _, r := range test.regexesToNotMatch {
				assert.NotRegexp(r, output.String())
			}
		})
	}
}
