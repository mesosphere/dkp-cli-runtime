// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"runtime"
	"sort"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

var (
	major        = "0"
	minor        = "0"
	gitVersion   = "v0.0.0-dev"
	gitCommit    = ""
	gitTreeState = ""
	commitDate   = "1970-01-01T00:00:00Z"
)

// Version is a struct for version information.
type Version struct {
	Major        string `json:"major"`
	Minor        string `json:"minor"`
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	CommitDate   string `json:"commitDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// GetVersion returns this binary's version.
func GetVersion() Version {
	return Version{
		Major:        major,
		Minor:        minor,
		GitVersion:   gitVersion,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		CommitDate:   commitDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// Versions is a map of multiple components' version information.
type Versions = map[string]Version

// options is a struct to support version command.
type options struct {
	Long   bool
	Output string

	outWriter   io.Writer
	getVersions func() (Versions, error)
}

// NewCommand returns a cobra command for fetching versions.
func NewCommand(output io.Writer) *cobra.Command {
	return NewCommandWithVersionGetter(output, func() (Versions, error) {
		return Versions{
			"": GetVersion(),
		}, nil
	})
}

// NewCommandWithVersionGetter returns a custom cobra command for fetching versions.
func NewCommandWithVersionGetter(output io.Writer, getVersions func() (Versions, error)) *cobra.Command {
	options := &options{
		outWriter:   output,
		getVersions: getVersions,
	}
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := options.Validate()
			if err != nil {
				return err
			}
			return options.Run()
		},
	}
	cmd.Flags().BoolVar(&options.Long, "long", options.Long, "If true, print additional version information.")
	cmd.Flags().StringVarP(&options.Output, "output", "o", options.Output, "One of 'yaml' or 'json'.")
	return cmd
}

// Validate validates the provided options.
func (o *options) Validate() error {
	if o.Output != "" && o.Output != "yaml" && o.Output != "json" {
		return errors.New(`--output must be 'yaml' or 'json'`)
	}

	return nil
}

// Run executes version command.
func (o *options) Run() error {
	versions, err := o.getVersions()
	if err != nil {
		return err
	}
	var v interface{}
	if len(versions) > 1 {
		v = versions
	} else {
		v = versions[""]
	}

	switch o.Output {
	case "":
		switch v := v.(type) {
		case Version:
			if !o.Long {
				fmt.Fprintln(o.outWriter, v.GitVersion)
			} else {
				fmt.Fprintf(o.outWriter, "%#v\n", v)
			}
		case Versions:
			keys := orderedKeys(v)
			for _, name := range keys {
				if !o.Long {
					fmt.Fprintln(o.outWriter, name+":", v[name].GitVersion)
				} else {
					fmt.Fprintf(o.outWriter, "%s: %#v\n", name, v[name])
				}
			}
		}
	case "yaml":
		marshalled, err := yaml.Marshal(&v)
		if err != nil {
			return err
		}
		fmt.Fprintln(o.outWriter, string(marshalled))
	case "json":
		marshalled, err := json.MarshalIndent(&v, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(o.outWriter, string(marshalled))
	default:
		// There is a bug in the program if we hit this case.
		// However, we follow a policy of never panicking.
		return fmt.Errorf("VersionOptions were not validated: --output=%q should have been rejected", o.Output)
	}

	return nil
}

func orderedKeys(mapp Versions) []string {
	keys := make([]string, 0, len(mapp))
	for key := range mapp {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
