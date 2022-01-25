// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/utils/pointer"
)

// GetOptions contains options for retrieving Kubernetes resources.
type GetOptions struct {
	cliRuntimeFlags *genericclioptions.ResourceBuilderFlags
}

// NewGetOptions creates GetOptions with default settings.
func NewGetOptions() *GetOptions {
	cliRuntimeFlags := &genericclioptions.ResourceBuilderFlags{
		LabelSelector: pointer.String(""),
		FieldSelector: pointer.String(""),
		AllNamespaces: pointer.Bool(false),
	}
	cliRuntimeFlags = cliRuntimeFlags.WithFile(true)

	return &GetOptions{
		cliRuntimeFlags: cliRuntimeFlags,
	}
}

// AddFlags adds flags for configuring get options to the provided FlagSet.
func (o *GetOptions) AddFlags(flags *pflag.FlagSet) {
	o.cliRuntimeFlags.AddFlags(flags)
}

// AllNamespaces returns true if getting resources from all namespaces is enabled.
func (o *GetOptions) AllNamespaces() bool {
	return o.cliRuntimeFlags != nil && *o.cliRuntimeFlags.AllNamespaces
}

// ToResourceBuilder creates a resource builder that can be used to retrieve generic resources
// from the API server. The `getArgs` slice must be
// of the form `(<type1>[,<type2>,...]|<type> <name1>[,<name2>,...])`. When one argument is
// received, the types provided will be retrieved from the server (and be comma delimited).
// When two or more arguments are received, they must be a single type and resource name(s).
// ToBuilder gives you back a resource finder to visit resources that are located.
func (o *GetOptions) ToResourceBuilder(restClientGetter RESTClientGetter, resources ...string) *resource.Builder {
	namespace, enforceNamespace, namespaceErr := restClientGetter.ToRawKubeConfigLoader().Namespace()

	builder := resource.NewBuilder(restClientGetter).
		NamespaceParam(namespace).DefaultNamespace()

	if o.cliRuntimeFlags.AllNamespaces != nil {
		builder.AllNamespaces(*o.cliRuntimeFlags.AllNamespaces)
	}

	if o.cliRuntimeFlags.Scheme != nil {
		builder.WithScheme(o.cliRuntimeFlags.Scheme, o.cliRuntimeFlags.Scheme.PrioritizedVersionsAllGroups()...)
	} else {
		builder.Unstructured()
	}

	if o.cliRuntimeFlags.FileNameFlags != nil {
		opts := o.cliRuntimeFlags.FileNameFlags.ToOptions()
		builder.FilenameParam(enforceNamespace, &opts)
	}

	if o.cliRuntimeFlags.Local == nil || !*o.cliRuntimeFlags.Local {
		builder.ResourceTypeOrNameArgs(true, resources...)
		// label selectors only work non-local (for now)
		if o.cliRuntimeFlags.LabelSelector != nil {
			builder.LabelSelectorParam(*o.cliRuntimeFlags.LabelSelector)
		}
		// field selectors only work non-local (forever)
		if o.cliRuntimeFlags.FieldSelector != nil {
			builder.FieldSelectorParam(*o.cliRuntimeFlags.FieldSelector)
		}
		// latest only works non-local (forever)
		if o.cliRuntimeFlags.Latest {
			builder.Latest()
		}
	} else {
		builder.Local()

		if len(resources) > 0 {
			builder.AddError(resource.LocalResourceError)
		}
	}

	if !o.cliRuntimeFlags.StopOnFirstError {
		builder.ContinueOnError()
	}

	return builder.
		Flatten().
		AddError(namespaceErr)
}
