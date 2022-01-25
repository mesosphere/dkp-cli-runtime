// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package help

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// NewHelpCommandWrapper creates an enhanced help command supporting multiple output formats.
func NewHelpCommandWrapper(rootCmd *cobra.Command) *cobra.Command {
	var (
		outputFormat  string
		showTree      bool
		treeOutputDir string
	)

	helpCmd := &cobra.Command{
		Use:   "help [-o {yaml|yml|markdown|md}] [command]",
		Short: "Help about any command",
		Long: `Help provides help for any command in the application.
Simply type ` + rootCmd.Name() + ` help [path to command] for full details.`,
		RunE: func(c *cobra.Command, args []string) error {
			cmd, _, e := c.Root().Find(args)
			if cmd == nil || e != nil {
				c.Printf("Unknown help topic %#q\n", args)
				return c.Root().Usage()
			} else {
				switch outputFormat {
				case "yaml", "yml":
					if showTree {
						return doc.GenYamlTree(cmd, treeOutputDir)
					}
					return doc.GenYaml(cmd, rootCmd.OutOrStdout())
				case "md", "markdown":
					if showTree {
						return doc.GenMarkdownTree(cmd, treeOutputDir)
					}
					return doc.GenMarkdown(cmd, rootCmd.OutOrStdout())
				case "":
					cmd.InitDefaultHelpFlag()
					return cmd.Help()
				default:
					return fmt.Errorf("unsupported help output: %s", outputFormat)
				}
			}
		},
	}

	helpCmd.ValidArgsFunction = func(
		c *cobra.Command, args []string, toComplete string,
	) ([]string, cobra.ShellCompDirective) {
		var completions []string
		cmd, _, e := c.Root().Find(args)
		if e != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		if cmd == nil {
			// Root help command.
			cmd = c.Root()
		}
		for _, subCmd := range cmd.Commands() {
			if subCmd.IsAvailableCommand() || subCmd == helpCmd {
				if strings.HasPrefix(subCmd.Name(), toComplete) {
					completions = append(completions, fmt.Sprintf("%s\t%s", subCmd.Name(), subCmd.Short))
				}
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	helpCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format for help")
	helpCmd.Flags().BoolVarP(&showTree, "tree", "t", false, "Generate help for full command tree")
	helpCmd.Flags().StringVarP(&treeOutputDir, "output-dir", "d", "", "Output for full command tree if --tree=true")

	return helpCmd
}
