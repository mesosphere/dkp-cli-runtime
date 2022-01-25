// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package plugin

import (
	"github.com/spf13/cobra"
)

// CommandSpec is a serializable object that has all the information that describes a cobra.Command.
type CommandSpec struct {
	Use                        string            `json:"use"`
	Aliases                    []string          `json:"aliases,omitempty"`
	SuggestFor                 []string          `json:"suggest_for,omitempty"`
	Short                      string            `json:"short,omitempty"`
	Long                       string            `json:"long,omitempty"`
	Example                    string            `json:"example,omitempty"`
	Deprecated                 string            `json:"deprecated,omitempty"`
	Annotations                map[string]string `json:"annotations,omitempty"`
	Version                    string            `json:"version,omitempty"`
	TraverseChildren           bool              `json:"traverse_children,omitempty"`
	Hidden                     bool              `json:"hidden,omitempty"`
	DisableAutoGenTag          bool              `json:"disable_auto_gen_tag,omitempty"`
	DisableFlagsInUseLine      bool              `json:"disable_flags_in_use_line,omitempty"`
	DisableSuggestions         bool              `json:"disable_suggestions,omitempty"`
	SuggestionsMinimumDistance int               `json:"suggestions_minimum_distance,omitempty"`

	LocalFlags      []FlagSpec    `json:"local_flags,omitempty"`
	PersistentFlags []FlagSpec    `json:"persistent_flags,omitempty"`
	Runnable        bool          `json:"runnable,omitempty"`
	SubCommands     []CommandSpec `json:"sub_commands,omitempty"`
}

// SpecFromCommand creates a CommandSpec describing a cobra.Command.
func SpecFromCommand(cmd *cobra.Command) CommandSpec {
	result := CommandSpec{
		Use:                        cmd.Use,
		Aliases:                    cmd.Aliases,
		SuggestFor:                 cmd.SuggestFor,
		Short:                      cmd.Short,
		Long:                       cmd.Long,
		Example:                    cmd.Example,
		Deprecated:                 cmd.Deprecated,
		Annotations:                cmd.Annotations,
		Version:                    cmd.Version,
		TraverseChildren:           cmd.TraverseChildren,
		Hidden:                     cmd.Hidden,
		DisableAutoGenTag:          cmd.DisableAutoGenTag,
		DisableFlagsInUseLine:      cmd.DisableFlagsInUseLine,
		DisableSuggestions:         cmd.DisableSuggestions,
		SuggestionsMinimumDistance: cmd.SuggestionsMinimumDistance,

		LocalFlags:      SpecsFromFlagset(cmd.LocalNonPersistentFlags()),
		PersistentFlags: SpecsFromFlagset(cmd.PersistentFlags()),
		Runnable:        cmd.Runnable(),
	}

	for _, subCmd := range cmd.Commands() {
		switch subCmd.Name() {
		case DiscoveryCommandName, "completion":
			continue
		}
		result.SubCommands = append(result.SubCommands, SpecFromCommand(subCmd))
	}
	return result
}

// ToCommand creates a cobra.Command based on this CommandSpec.
func (c CommandSpec) ToCommand(runOverride func(*cobra.Command, []string) error) *cobra.Command {
	command := &cobra.Command{
		Use:                        c.Use,
		Aliases:                    c.Aliases,
		SuggestFor:                 c.SuggestFor,
		Short:                      c.Short,
		Long:                       c.Long,
		Example:                    c.Example,
		Deprecated:                 c.Deprecated,
		Annotations:                c.Annotations,
		Version:                    c.Version,
		TraverseChildren:           c.TraverseChildren,
		Hidden:                     c.Hidden,
		DisableAutoGenTag:          c.DisableAutoGenTag,
		DisableFlagsInUseLine:      c.DisableFlagsInUseLine,
		DisableSuggestions:         c.DisableSuggestions,
		SuggestionsMinimumDistance: c.SuggestionsMinimumDistance,
	}

	for _, flagSpec := range c.PersistentFlags {
		command.PersistentFlags().AddFlag(flagSpec.ToFlag())
	}
	for _, flagSpec := range c.LocalFlags {
		command.Flags().AddFlag(flagSpec.ToFlag())
	}

	for _, subSpec := range c.SubCommands {
		command.AddCommand(subSpec.ToCommand(runOverride))
	}

	if c.Runnable {
		command.RunE = runOverride
	}
	return command
}
