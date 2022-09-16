// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package customdocs_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mesosphere/dkp-cli-runtime/core/customdocs"
)

func TestGenWithTemplate(t *testing.T) {
	_, cmd := setupTestCommands()

	tpl, err := template.ParseFiles(filepath.Join("testfiles", "template"))
	require.NoError(t, err)

	output := new(bytes.Buffer)
	err = customdocs.GenWithTemplate(cmd, output, tpl)
	require.NoError(t, err)

	expected, err := os.ReadFile(filepath.Join("testfiles", "example_command"))
	require.NoError(t, err)

	assert.Equal(t, string(expected), output.String())
}

func TestGenTreeWithTemplate(t *testing.T) {
	cmd, _ := setupTestCommands()

	tpl, err := template.ParseFiles(filepath.Join("testfiles", "template"))
	require.NoError(t, err)

	tempDir := t.TempDir()
	err = customdocs.GenTreeWithTemplate(cmd, tempDir, tpl)
	require.NoError(t, err)

	files, err := os.ReadDir(tempDir)
	require.NoError(t, err)
	for _, file := range files {
		expected, err := os.ReadFile(filepath.Join("testfiles", file.Name()))
		require.NoError(t, err)
		actual, err := os.ReadFile(filepath.Join(tempDir, file.Name()))
		require.NoError(t, err)

		assert.Equal(t, string(expected), string(actual))
	}
}

func setupTestCommands() (*cobra.Command, *cobra.Command) {
	rootCmd := &cobra.Command{
		Use:   "example",
		Short: "Example command",
	}
	rootCmd.PersistentFlags().Bool("global", false, "a global flag")

	cmd := &cobra.Command{
		Use:     "command",
		Short:   "Example sub-command",
		Long:    "This is just a sample for testing doc generation",
		Example: "Usage example",
		Run:     func(cmd *cobra.Command, args []string) {},
	}
	cmd.Flags().String("local", "default", "a local flag")
	rootCmd.AddCommand(cmd)

	subCmd := &cobra.Command{
		Use:   "sub",
		Short: "Example sub-sub-command",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	cmd.AddCommand(subCmd)

	return rootCmd, cmd
}
