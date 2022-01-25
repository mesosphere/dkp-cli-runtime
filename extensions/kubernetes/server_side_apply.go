// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const decoderBufferSize = 4096

// ServerSideApplyResource applies the provided resource using server side apply.
func ServerSideApplyResource(
	ctx context.Context,
	obj client.Object,
	c client.Client,
	options ...client.PatchOption) error {
	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return fmt.Errorf("failed to convert resource to unstructured: %w", err)
	}
	unstructuredObj := &unstructured.Unstructured{
		Object: unstructuredMap,
	}

	if err := c.Patch(ctx, unstructuredObj, client.Apply, options...); err != nil {
		return fmt.Errorf("failed to patch resource %s (%s): %w",
			obj.GetName(), obj.GetObjectKind().GroupVersionKind().String(), err)
	}

	return nil
}

// ServerSideApplyResourceFromString applies the provided resource using server side apply.
func ServerSideApplyResourceFromString(
	ctx context.Context,
	s string,
	c client.Client,
	options ...client.PatchOption) error {
	return ServerSideApplyResourceFromReader(ctx, strings.NewReader(s), c, options...)
}

// ServerSideApplyResourceFromBytes applies the provided resource using server side apply.
func ServerSideApplyResourceFromBytes(
	ctx context.Context,
	b []byte,
	c client.Client,
	options ...client.PatchOption) error {
	return ServerSideApplyResourceFromReader(ctx, bytes.NewReader(b), c, options...)
}

// ServerSideApplyResourceFromReader applies the provided resource using server side apply.
func ServerSideApplyResourceFromReader(
	ctx context.Context,
	r io.Reader,
	c client.Client,
	options ...client.PatchOption) error {
	decoder := yaml.NewYAMLOrJSONDecoder(r, decoderBufferSize)

	for {
		var obj unstructured.Unstructured
		err := decoder.Decode(&obj)
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("failed to decode manifest: %w", err)
		}

		if err := ServerSideApplyResource(ctx, &obj, c, options...); err != nil {
			return err
		}
	}

	return nil
}
