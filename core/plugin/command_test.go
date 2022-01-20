package plugin_test

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/mesosphere/dkp-cli-runtime/core/plugin"
)

func TestCommandsToSpecAndBack(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cmd := &cobra.Command{Use: "simple"}
		testCommandSpec(t, cmd)
	})

	t.Run("with flags", func(t *testing.T) {
		cmdWithFlags := &cobra.Command{
			Use:   "with_flags",
			Short: "a command with flags",
		}
		cmdWithFlags.Flags().String("flag1", "", "usage")
		cmdWithFlags.PersistentFlags().String("flag2", "", "usage")
		testCommandSpec(t, cmdWithFlags)
	})

	t.Run("with flags runnable", func(t *testing.T) {
		cmdWithFlagsRunnable := &cobra.Command{
			Use:   "with_flags_runnable",
			Short: "a command with flags",
			Run:   func(cmd *cobra.Command, args []string) {},
		}
		cmdWithFlagsRunnable.Flags().String("flag1", "", "usage")
		cmdWithFlagsRunnable.PersistentFlags().String("flag2", "", "usage")
		testCommandSpec(t, cmdWithFlagsRunnable)
	})

	t.Run("with subcommands", func(t *testing.T) {
		cmdWithSubCommands := &cobra.Command{
			Use:   "with_subcommands",
			Short: "a command with subcommands",
		}
		cmdWithSubCommands.Flags().String("flag1", "", "usage")
		cmdWithSubCommands.PersistentFlags().String("flag2", "", "usage")

		subCmd := &cobra.Command{
			Use: "subcommand",
		}
		subCmd.Flags().String("flag3", "", "usage")
		cmdWithSubCommands.AddCommand(subCmd)

		subCmd2 := &cobra.Command{
			Use: "subcommand2",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		subCmd2.Flags().String("flag4", "", "usage")
		cmdWithSubCommands.AddCommand(subCmd2)

		testCommandSpec(t, cmdWithSubCommands)
	})
}

func testCommandSpec(t *testing.T, cmd *cobra.Command) {
	spec := plugin.SpecFromCommand(cmd)
	runE := func(cmd *cobra.Command, args []string) error { return nil }
	cmdFromSpec := spec.ToCommand(runE)
	assertCommandsProduceSameHelpOutput(t, cmd, cmdFromSpec)
	assertCommandsProduceSameCompletion(t, cmd, cmdFromSpec)
}

func assertCommandsProduceSameHelpOutput(t *testing.T, expected, actual *cobra.Command) {
	expectedOutput := bytes.Buffer{}
	expected.SetOut(&expectedOutput)
	expected.Usage()

	actualOutput := bytes.Buffer{}
	actual.SetOut(&actualOutput)
	actual.Usage()

	assert.Equal(t, expectedOutput.String(), actualOutput.String())

	for _, expectedSubCmd := range expected.Commands() {
		for _, actualSubCmd := range actual.Commands() {
			if actualSubCmd.Use == expectedSubCmd.Use {
				assertCommandsProduceSameHelpOutput(t, expectedSubCmd, actualSubCmd)
			}
		}
	}
}

func assertCommandsProduceSameCompletion(t *testing.T, expected, actual *cobra.Command) {
	for _, expectedSubCmd := range expected.Commands() {
		for _, actualSubCmd := range actual.Commands() {
			if actualSubCmd.Use == expectedSubCmd.Use {
				assertCommandsProduceSameCompletion(t, expectedSubCmd, actualSubCmd)
			}
		}
	}

	expectedOutput := bytes.Buffer{}
	actualOutput := bytes.Buffer{}

	expected.GenBashCompletion(&expectedOutput)
	actual.GenBashCompletion(&actualOutput)
	assert.Equal(t, expectedOutput.String(), actualOutput.String())

	expectedOutput.Reset()
	actualOutput.Reset()
	expected.GenFishCompletion(&expectedOutput, true)
	actual.GenFishCompletion(&actualOutput, true)
	assert.Equal(t, expectedOutput.String(), actualOutput.String())

	expectedOutput.Reset()
	actualOutput.Reset()
	expected.GenZshCompletion(&expectedOutput)
	actual.GenZshCompletion(&actualOutput)
	assert.Equal(t, expectedOutput.String(), actualOutput.String())

	expectedOutput.Reset()
	actualOutput.Reset()
	expected.GenPowerShellCompletion(&expectedOutput)
	actual.GenPowerShellCompletion(&actualOutput)
	assert.Equal(t, expectedOutput.String(), actualOutput.String())
}
