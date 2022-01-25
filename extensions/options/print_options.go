// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
)

// PrintOptions holds settings for output formatting.
type PrintOptions struct {
	printFlags *genericclioptions.PrintFlags

	outputFormat string
}

// NewPrintOptions creates PrintOptions with default settings.
func NewPrintOptions() *PrintOptions {
	printFlags := genericclioptions.NewPrintFlags("get")
	printFlags.OutputFormat = nil
	return &PrintOptions{
		printFlags:   printFlags,
		outputFormat: "table",
	}
}

// AddFlags adds flags for setting print options to a command.
func (o *PrintOptions) AddFlags(cmd *cobra.Command) {
	o.printFlags.AddFlags(cmd)

	cmd.Flags().StringVarP(&o.outputFormat, "output", "o", o.outputFormat,
		fmt.Sprintf("Output format. One of: %s.", strings.Join(o.AllowedFormats(), "|")))
	o.printFlags.OutputFlagSpecified = func() bool {
		return cmd.Flag("output").Changed
	}
}

// WithOutputFormat sets the output format.
func (o *PrintOptions) WithOutputFormat(outputFormat string) *PrintOptions {
	o.outputFormat = outputFormat
	return o
}

// OutputFormat returns the currently set output format.
func (o *PrintOptions) OutputFormat() string {
	return o.outputFormat
}

// AllowedFormats returns a list of supported output formats.
func (o *PrintOptions) AllowedFormats() []string {
	allowedFormats := append([]string{"table"}, o.printFlags.AllowedFormats()...)
	sort.Strings(allowedFormats)
	return allowedFormats
}

// ToPrinter creates a ResourcePrinter with the configured format.
func (o *PrintOptions) ToPrinter(
	allNamespaces, multipleGVKs bool, gk schema.GroupKind,
) (printers.ResourcePrinter, error) {
	if o.outputFormat == "table" {
		return printers.NewTablePrinter(printers.PrintOptions{
				WithNamespace: allNamespaces,
				WithKind:      multipleGVKs,
				Kind:          gk,
			}),
			nil
	}

	o.printFlags.OutputFlagSpecified = func() bool { return true }
	o.printFlags.OutputFormat = &o.outputFormat
	printer, err := o.printFlags.ToPrinter()
	return printer, err
}
