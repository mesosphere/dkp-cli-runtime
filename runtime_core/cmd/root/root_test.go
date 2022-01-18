package root_test

import (
	"io"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"

	"github.com/mesosphere/dkp-cli-runtime/runtime_core/cmd/root"
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
