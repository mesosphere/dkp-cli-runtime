// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package kubernetes_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mesosphere/dkp-cli-runtime/extensions/kubernetes"
)

func TestServerSideApply(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	ctx := context.Background()
	t.Run("with object", func(t *testing.T) {
		client := &MockedClient{}
		pod := &corev1.Pod{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Pod",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "namespace",
				Name:      "name",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Image: "nginx",
						Name:  "nginx",
					},
				},
			},
		}
		err := kubernetes.ServerSideApplyResource(ctx, pod, client)
		require.NoError(err)
		require.Len(client.calledWithObject, 1)
		assert.Equal(getPodObjectFull("namespace", "name"), client.calledWithObject[0])
	})

	t.Run("with YAML string", func(t *testing.T) {
		client := &MockedClient{}
		err := kubernetes.ServerSideApplyResourceFromString(ctx, podYAML, client)
		require.NoError(err)
		require.Len(client.calledWithObject, 1)
		assert.Equal(getPodObject("namespace", "name"), client.calledWithObject[0])
	})

	t.Run("with JSON string", func(t *testing.T) {
		client := &MockedClient{}
		err := kubernetes.ServerSideApplyResourceFromString(ctx, podJSON, client)
		require.NoError(err)
		require.Len(client.calledWithObject, 1)
		assert.Equal(getPodObject("namespace", "name"), client.calledWithObject[0])
	})

	t.Run("with invalid JSON string", func(t *testing.T) {
		client := &MockedClient{}
		err := kubernetes.ServerSideApplyResourceFromString(ctx, "{definitely not JSON", client)
		require.Error(err)
		require.Len(client.calledWithObject, 0)
	})

	t.Run("with YAML string multiple", func(t *testing.T) {
		client := &MockedClient{}
		err := kubernetes.ServerSideApplyResourceFromString(ctx, multiPodYAML, client)
		require.NoError(err)
		require.Len(client.calledWithObject, 2)
		assert.Equal(getPodObject("namespace", "name"), client.calledWithObject[0])
		assert.Equal(getPodObject("namespace2", "name2"), client.calledWithObject[1])
	})

	t.Run("with YAML bytes", func(t *testing.T) {
		client := &MockedClient{}
		err := kubernetes.ServerSideApplyResourceFromBytes(ctx, []byte(podYAML), client)
		require.NoError(err)
		require.Len(client.calledWithObject, 1)
		assert.Equal(getPodObject("namespace", "name"), client.calledWithObject[0])
	})

	t.Run("with YAML reader", func(t *testing.T) {
		client := &MockedClient{}
		reader := bytes.NewBuffer([]byte(podYAML))
		err := kubernetes.ServerSideApplyResourceFromReader(ctx, reader, client)
		require.NoError(err)
		require.Len(client.calledWithObject, 1)
		assert.Equal(getPodObject("namespace", "name"), client.calledWithObject[0])
	})

	t.Run("server side error", func(t *testing.T) {
		client := &MockedClient{
			returnErr: errors.New("server side error"),
		}
		err := kubernetes.ServerSideApplyResourceFromString(ctx, podJSON, client)
		require.Error(err)
	})
}

const podYAML = `
apiVersion: v1
kind: Pod
metadata:
    namespace: namespace
    name: name
spec:
    containers:
    - name: nginx
      image: nginx
`

const podJSON = `
{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "namespace": "namespace",
        "name": "name"
    },
    "spec": {
        "containers": [
            {
                "name": "nginx",
                "image": "nginx"
            }
        ]
    }
}
`

const multiPodYAML = `
apiVersion: v1
kind: Pod
metadata:
    namespace: namespace
    name: name
spec:
    containers:
    - name: nginx
      image: nginx
---
apiVersion: v1
kind: Pod
metadata:
    namespace: namespace2
    name: name2
spec:
    containers:
    - name: nginx
      image: nginx
`

func getPodObject(namespace, name string) client.Object {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{
						"name":  "nginx",
						"image": "nginx",
					},
				},
			},
		},
	}
}

func getPodObjectFull(namespace, name string) client.Object {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]interface{}{
				"name":              name,
				"namespace":         namespace,
				"creationTimestamp": nil,
			},
			"spec": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{
						"name":      "nginx",
						"image":     "nginx",
						"resources": map[string]interface{}{},
					},
				},
			},
			"status": map[string]interface{}{},
		},
	}
}

type MockedClient struct {
	calledWithObject []client.Object
	returnErr        error
}

func (c *MockedClient) Patch(
	ctx context.Context,
	obj client.Object,
	patch client.Patch,
	opts ...client.PatchOption,
) error {
	c.calledWithObject = append(c.calledWithObject, obj)
	return c.returnErr
}

func (c *MockedClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object) error {
	return nil
}

func (c *MockedClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	return nil
}

func (c *MockedClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	return nil
}

func (c *MockedClient) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	return nil
}

func (c *MockedClient) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	return nil
}

func (c *MockedClient) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	return nil
}

func (c *MockedClient) Status() client.StatusWriter {
	return nil
}

func (c *MockedClient) Scheme() *runtime.Scheme {
	return nil
}

func (c *MockedClient) RESTMapper() meta.RESTMapper {
	return nil
}
