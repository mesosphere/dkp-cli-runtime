// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"fmt"

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClientOptions holds options for Kubernetes clients.
type ClientOptions struct {
	cliRuntimeFlags *genericclioptions.ConfigFlags

	fieldOwner string
	qps        float32
	burst      int
}

// RESTClientGetter is an interface that the ClientOptions describe to provide an easier way to mock for commands
// and eliminate the direct coupling to a struct type. Users may wish to duplicate this type in their own packages
// as per the golang type overlapping.
type RESTClientGetter interface {
	// ToRESTConfig returns restconfig
	ToRESTConfig() (*rest.Config, error)
	// ToDiscoveryClient returns discovery client
	ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error)
	// ToRESTMapper returns a restmapper
	ToRESTMapper() (meta.RESTMapper, error)
	// ToRawKubeConfigLoader return kubeconfig loader as-is
	ToRawKubeConfigLoader() clientcmd.ClientConfig
}

var _ RESTClientGetter = &ClientOptions{}

// NewClientOptions creates ClientOptions with default settings.
func NewClientOptions(usePersistentConfig bool) *ClientOptions {
	configFlags := genericclioptions.NewConfigFlags(usePersistentConfig)

	configFlags.APIServer = nil
	configFlags.AuthInfoName = nil
	configFlags.BearerToken = nil
	configFlags.CAFile = nil
	configFlags.CacheDir = nil
	configFlags.CertFile = nil
	configFlags.ClusterName = nil
	configFlags.Impersonate = nil
	configFlags.ImpersonateGroup = nil
	configFlags.Insecure = nil
	configFlags.KeyFile = nil
	configFlags.TLSServerName = nil

	return &ClientOptions{
		cliRuntimeFlags: configFlags,
		qps:             20.0,
		burst:           50,
	}
}

// WithImpersonation enables support for impersonation.
func (o *ClientOptions) WithImpersonation() *ClientOptions {
	o.cliRuntimeFlags.Impersonate = pointer.String("")
	impersonateGroup := []string{}
	o.cliRuntimeFlags.ImpersonateGroup = &impersonateGroup
	return o
}

// WithDiscoveryBurst sets the RESTClient burst for discovery.
func (o *ClientOptions) WithDiscoveryBurst(discoveryBurst int) *ClientOptions {
	o.cliRuntimeFlags = o.cliRuntimeFlags.WithDiscoveryBurst(discoveryBurst)
	return o
}

// WithKubeAPIQPS sets the the maximum QPS to the master from this client.
func (o *ClientOptions) WithKubeAPIQPS(qps float32) *ClientOptions {
	o.qps = qps
	return o
}

// WithKubeAPIBurst sets the maximum burst for the rate limit.
func (o *ClientOptions) WithKubeAPIBurst(burst int) *ClientOptions {
	o.burst = burst
	return o
}

// AddFlags adds flags for configuring the Kubernetes clients to the provided FlagSet.
func (o *ClientOptions) AddFlags(flags *pflag.FlagSet) {
	o.cliRuntimeFlags.AddFlags(flags)

	flags.Float32Var(&o.qps, "kube-api-qps", 20.0,
		"The maximum queries-per-second of requests sent to the Kubernetes API.")
	flags.IntVar(&o.burst, "kube-api-burst", 50,
		"The maximum burst queries-per-second of requests sent to the Kubernetes API.")

	flags.StringVar(
		&o.fieldOwner, "field-owner", o.fieldOwner,
		"Name of the field owner to use with server-side apply (SSA).",
	)
}

// WithFieldOwner sets the field owner.
func (o *ClientOptions) WithFieldOwner(fieldOwner string) *ClientOptions {
	o.fieldOwner = fieldOwner
	return o
}

// FieldOwner returns the field owner.
func (o *ClientOptions) FieldOwner() string {
	return o.fieldOwner
}

// ToRESTConfig implements RESTClientGetter.
// Returns a REST client configuration based on a provided path
// to a .kubeconfig file, loading rules, and config flag overrides.
// Expects the AddFlags method to have been called. If WrapConfigFn
// is non-nil this function can transform config before return.
func (o *ClientOptions) ToRESTConfig() (*rest.Config, error) {
	o.cliRuntimeFlags.WrapConfigFn = func(c *rest.Config) *rest.Config {
		c.QPS = o.qps
		c.Burst = o.burst
		return c
	}
	return o.cliRuntimeFlags.ToRESTConfig()
}

// ToRawKubeConfigLoader binds config flag values to config overrides
// Returns an interactive clientConfig if the password flag is enabled,
// or a non-interactive clientConfig otherwise.
func (o *ClientOptions) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return o.cliRuntimeFlags.ToRawKubeConfigLoader()
}

// ToDiscoveryClient implements RESTClientGetter.
// Expects the AddFlags method to have been called.
// Returns a CachedDiscoveryInterface using a computed RESTConfig.
func (o *ClientOptions) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	o.cliRuntimeFlags.WrapConfigFn = func(c *rest.Config) *rest.Config {
		c.QPS = o.qps
		c.Burst = o.burst
		return c
	}
	return o.cliRuntimeFlags.ToDiscoveryClient()
}

// ToRESTMapper returns a mapper.
func (o *ClientOptions) ToRESTMapper() (meta.RESTMapper, error) {
	return o.cliRuntimeFlags.ToRESTMapper()
}

// KubernetesClient returns a (typed) configured Kubernetes client.
func (o *ClientOptions) KubernetesClient() (kubernetes.Interface, error) {
	restConfig, err := o.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes REST config: %v", err)
	}

	kc, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	return kc, nil
}

// ControllerRuntimeClient returns a (dynamic) controller runtime Kubernets client with default settings.
func (o *ClientOptions) ControllerRuntimeClient() (client.Client, error) {
	return o.ControllerRuntimeClientWithOptions(client.Options{})
}

// ControllerRuntimeClient returns a (dynamic) controller runtime Kubernets client with the provided settings.
func (o *ClientOptions) ControllerRuntimeClientWithOptions(opts client.Options) (client.Client, error) {
	restConfig, err := o.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes REST config: %v", err)
	}

	kc, err := client.New(restConfig, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	return kc, nil
}
