// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/utils/pointer"
)

// DeleteOptions contains settings for deleting Kubernetes resources.
type DeleteOptions struct {
	cliRuntimeFlags *genericclioptions.ResourceBuilderFlags
}

// NewDeleteOptions creates DeleteOptions with default settings.
func NewDeleteOptions() *DeleteOptions {
	cliRuntimeFlags := &genericclioptions.ResourceBuilderFlags{
		LabelSelector: pointer.String(""),
		FieldSelector: pointer.String(""),
		AllNamespaces: pointer.Bool(false),
		All:           pointer.Bool(false),
	}
	cliRuntimeFlags = cliRuntimeFlags.WithFile(true)

	return &DeleteOptions{
		cliRuntimeFlags: cliRuntimeFlags,
	}
}

// AddFlags adds flags for configuring resource deletion to the provided FlagSet.
func (o *DeleteOptions) AddFlags(flags *pflag.FlagSet) {
	o.cliRuntimeFlags.AddFlags(flags)
}

// ToResourceBuilder creates a resource builder that can be used to retrieve generic resources
// from the API server. The `getArgs` slice must be
// of the form `(<type1>[,<type2>,...]|<type> <name1>[,<name2>,...])`. When one argument is
// received, the types provided will be retrieved from the server (and be comma delimited).
// When two or more arguments are received, they must be a single type and resource name(s).
func (o *DeleteOptions) ToResourceBuilder(
	restClientGetter RESTClientGetter, getArgs ...string,
) genericclioptions.ResourceFinder {
	return o.cliRuntimeFlags.ToBuilder(restClientGetter, getArgs)
}
