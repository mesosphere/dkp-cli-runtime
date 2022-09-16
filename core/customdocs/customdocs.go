// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package customdocs

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// Template interface is implemented by core packages text/template and html/template.
type Template interface {
	Execute(wr io.Writer, data interface{}) error
}

type templateValues struct {
	Name        string
	Short       string
	Long        string
	UseLine     string
	Example     string
	Flags       string
	ParentFlags string
	Links       []linkValue
}

type linkValue struct {
	Name  string
	Short string
}

// GenWithTemplate outputs CLI docs based on the provided template (text/template or html/template).
func GenWithTemplate(cmd *cobra.Command, w io.Writer, template Template) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	values := templateValues{
		Name:    cmd.CommandPath(),
		Short:   cmd.Short,
		Long:    cmd.Long,
		Example: cmd.Example,
	}

	if cmd.Runnable() {
		values.UseLine = cmd.UseLine()
	}

	flags := cmd.NonInheritedFlags()
	if flags.HasAvailableFlags() {
		tempBuf := new(bytes.Buffer)
		flags.SetOutput(tempBuf)
		flags.PrintDefaults()
		values.Flags = tempBuf.String()
	}

	parentFlags := cmd.InheritedFlags()
	if parentFlags.HasAvailableFlags() {
		tempBuf := new(bytes.Buffer)
		parentFlags.SetOutput(tempBuf)
		parentFlags.PrintDefaults()
		values.ParentFlags = tempBuf.String()
	}

	values.Links = make([]linkValue, 0)
	if hasSeeAlso(cmd) {
		if cmd.HasParent() {
			parent := cmd.Parent()
			values.Links = append(values.Links, linkValue{
				Name:  parent.CommandPath(),
				Short: parent.Short,
			})
		}

		children := cmd.Commands()
		sort.Sort(byName(children))
		for _, child := range children {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			values.Links = append(values.Links, linkValue{
				Name:  child.CommandPath(),
				Short: child.Short,
			})
		}
	}

	return template.Execute(w, values)
}

// GenTreeWithTemplate outputs CLI docs for the command and all sub-commands into a directory,
// using the provided template.
func GenTreeWithTemplate(cmd *cobra.Command, dir string, template Template) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := GenTreeWithTemplate(c, dir, template); err != nil {
			return err
		}
	}

	basename := strings.ReplaceAll(cmd.CommandPath(), " ", "_")
	filename := filepath.Join(dir, basename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := GenWithTemplate(cmd, f, template); err != nil {
		return err
	}
	return nil
}

func hasSeeAlso(cmd *cobra.Command) bool {
	if cmd.HasParent() {
		return true
	}
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		return true
	}
	return false
}

type byName []*cobra.Command

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }
