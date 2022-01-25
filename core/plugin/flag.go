// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package plugin

import "github.com/spf13/pflag"

// FlagSpec is a serializable object that has all the information that describes a pflag.Flag.
type FlagSpec struct {
	Type                string              `json:"type,omitempty"`
	Name                string              `json:"name,omitempty"`
	Shorthand           string              `json:"shorthand,omitempty"`
	Usage               string              `json:"usage,omitempty"`
	DefaultValue        string              `json:"default_value,omitempty"`
	NoOptDefaultValue   string              `json:"no_opt_default_value,omitempty"`
	Deprecated          string              `json:"deprecated,omitempty"`
	Hidden              bool                `json:"hidden,omitempty"`
	ShorthandDeprecated string              `json:"shorthand_deprecated,omitempty"`
	Annotations         map[string][]string `json:"annotations,omitempty"`
}

// SpecFromFlag creates a FlagSpec describing a flag.
func SpecFromFlag(flag *pflag.Flag) FlagSpec {
	return FlagSpec{
		Type:                flag.Value.Type(),
		Name:                flag.Name,
		Shorthand:           flag.Shorthand,
		Usage:               flag.Usage,
		DefaultValue:        flag.DefValue,
		NoOptDefaultValue:   flag.NoOptDefVal,
		Deprecated:          flag.Deprecated,
		Hidden:              flag.Hidden,
		ShorthandDeprecated: flag.ShorthandDeprecated,
		Annotations:         flag.Annotations,
	}
}

// SpecsFromFlagset creates a slice of FlagSpecs from a pflag.FlagSet.
func SpecsFromFlagset(flags *pflag.FlagSet) []FlagSpec {
	flagSpecs := []FlagSpec{}
	flags.VisitAll(func(flag *pflag.Flag) {
		flagSpecs = append(flagSpecs, SpecFromFlag(flag))
	})
	return flagSpecs
}

// ToFlag creates a pflag.Flag based on this FlagSpec.
func (spec FlagSpec) ToFlag() *pflag.Flag {
	return &pflag.Flag{
		Name:                spec.Name,
		Shorthand:           spec.Shorthand,
		Usage:               spec.Usage,
		Value:               typedFlagValue(spec.Type),
		DefValue:            spec.DefaultValue,
		NoOptDefVal:         spec.NoOptDefaultValue,
		Deprecated:          spec.Deprecated,
		Hidden:              spec.Hidden,
		ShorthandDeprecated: spec.ShorthandDeprecated,
		Annotations:         spec.Annotations,
	}
}

// typedFlagValue is a dummy pflag.Value that knows its type. That's all we need here.
// All standard pflag.Value implementations are private.
type typedFlagValue string

func (s typedFlagValue) Set(val string) error {
	return nil
}

func (s typedFlagValue) Type() string {
	return string(s)
}

func (s typedFlagValue) String() string { return "" }
