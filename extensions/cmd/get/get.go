// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/rest"

	"github.com/mesosphere/dkp-cli-runtime/extensions/options"
)

// NewCommand creates a command for retrieving Kubernetes resources.
func NewCommand(
	ioStreams genericclioptions.IOStreams,
	clientOpts *options.ClientOptions,
	allowedResources ...string,
) *cobra.Command {
	getOpts := options.NewGetOptions()
	printOpts := options.NewPrintOptions()

	rootGetCmd := &cobra.Command{
		Use: "get",
	}

	clientOpts.AddFlags(rootGetCmd.PersistentFlags())
	getOpts.AddFlags(rootGetCmd.PersistentFlags())

	for _, resource := range allowedResources {
		getResourceCmd := &cobra.Command{
			Use:  resource,
			RunE: getRunEFn(ioStreams, getOpts, clientOpts, printOpts, resource),
		}
		printOpts.AddFlags(getResourceCmd)

		rootGetCmd.AddCommand(getResourceCmd)
	}

	if len(allowedResources) == 0 {
		printOpts.AddFlags(rootGetCmd)
		rootGetCmd.RunE = getRunEFn(ioStreams, getOpts, clientOpts, printOpts)
	} else {
		rootGetCmd.RunE = func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("unknown command %q", args[0])
		}
		rootGetCmd.Args = cobra.NoArgs
	}

	return rootGetCmd
}

func getRunEFn(
	ioStreams genericclioptions.IOStreams,
	getOpts *options.GetOptions,
	clientOpts *options.ClientOptions,
	printOpts *options.PrintOptions,
	resources ...string,
) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var requestTransformers []resource.RequestTransform
		if printOpts.OutputFormat() == "table" {
			requestTransformers = append(requestTransformers, serverSideTableRequestTransformer)
		}

		resourceResult := getOpts.ToResourceBuilder(clientOpts, append(resources, args...)...).
			TransformRequests(requestTransformers...).Do()

		infos, err := resourceResult.Infos()
		if err != nil {
			return fmt.Errorf("failed to retrieve result infos: %w", err)
		}

		resourcePrinter, err := printOpts.ToPrinter(
			getOpts.AllNamespaces(), false, schema.GroupKind{},
		)
		if err != nil {
			return fmt.Errorf("failed to create resource printer: %w", err)
		}
		withKind := multipleGVKsRequested(infos)

		resourceResultObject, err := resourceResult.Object()
		if err != nil {
			return fmt.Errorf("failed to get resource result: %w", err)
		}

		if meta.IsListType(resourceResultObject) {
			resourceResultObject.GetObjectKind().SetGroupVersionKind(schema.FromAPIVersionAndKind("v1", "List"))
			objUnst, err := runtime.DefaultUnstructuredConverter.ToUnstructured(resourceResultObject)
			if err != nil {
				return fmt.Errorf("failed to convert result to unstructured: %w", err)
			}
			objList := &metav1.List{}
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(objUnst, objList); err != nil {
				return fmt.Errorf("failed to convert result to list object: %w", err)
			}
			if len(objList.Items) > 0 {
				objRawItem := objList.Items[0]
				if objRawItem.Object == nil {
					unstItem := &unstructured.Unstructured{}
					if err := unstItem.UnmarshalJSON(objRawItem.Raw); err != nil {
						return fmt.Errorf("failed to unmarshal raw JSON: %w", err)
					}
					objRawItem.Object = unstItem
				}
				if objRawItem.Object.GetObjectKind().GroupVersionKind() != metav1.SchemeGroupVersion.WithKind("Table") {
					return resourcePrinter.PrintObj(resourceResultObject, ioStreams.Out)
				}
			}
		}

		return resourceResult.Visit(
			func(i *resource.Info, e error) error {
				obj := i.Object

				if obj.GetObjectKind().GroupVersionKind() == metav1.SchemeGroupVersion.WithKind("Table") {
					if unst, ok := obj.(*unstructured.Unstructured); ok {
						objTable := &metav1.Table{}
						err := runtime.DefaultUnstructuredConverter.FromUnstructured(unst.UnstructuredContent(), objTable)
						if err != nil {
							return nil
						}

						for idx := range objTable.Rows {
							objRow := &objTable.Rows[idx]
							if objRow.Object.Object == nil {
								unstRow := &unstructured.Unstructured{}
								if err := unstRow.UnmarshalJSON(objRow.Object.Raw); err != nil {
									return fmt.Errorf("failed to unmarshal table row: %w", err)
								}
								unstRow.SetGroupVersionKind(i.Mapping.GroupVersionKind)
								objRow.Object.Object = unstRow
							}
						}

						obj = objTable

						resourcePrinter, err = printOpts.ToPrinter(
							getOpts.AllNamespaces(), withKind, i.Mapping.GroupVersionKind.GroupKind(),
						)
						if err != nil {
							return fmt.Errorf("failed to create resource printer: %w", err)
						}
					}
				}

				return resourcePrinter.PrintObj(obj, ioStreams.Out)
			},
		)
	}
}

func serverSideTableRequestTransformer(req *rest.Request) {
	req.SetHeader("Accept", strings.Join([]string{
		fmt.Sprintf("application/json;as=Table;v=%s;g=%s", metav1.SchemeGroupVersion.Version, metav1.GroupName),
		"application/json",
	}, ","))
}

func multipleGVKsRequested(infos []*resource.Info) bool {
	if len(infos) < 2 {
		return false
	}
	gvk := infos[0].Mapping.GroupVersionKind
	for _, info := range infos {
		if info.Mapping.GroupVersionKind != gvk {
			return true
		}
	}
	return false
}
