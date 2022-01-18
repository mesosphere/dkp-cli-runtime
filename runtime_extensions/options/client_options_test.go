package options_test

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mesosphere/dkp-cli-runtime/runtime_extensions/options"
)

func TestClientOptions(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	t.Run("with_defaults", func(t *testing.T) {
		clientOptions := options.NewClientOptions(true)
		require.NotNil(clientOptions)

		assert.Equal("", clientOptions.FieldOwner())

		restConfig, err := clientOptions.ToRESTConfig()
		require.NoError(err)
		require.NotNil(restConfig)
		assert.EqualValues(20.0, restConfig.QPS)
		assert.Equal(50, restConfig.Burst)
		assert.Equal("", restConfig.Impersonate.UserName)
		assert.Len(restConfig.Impersonate.Groups, 0)

		rawConfig := clientOptions.ToRawKubeConfigLoader()
		assert.NotNil(rawConfig)

		discoveryClient, err := clientOptions.ToDiscoveryClient()
		assert.NoError(err)
		assert.NotNil(discoveryClient)

		restMapper, err := clientOptions.ToRESTMapper()
		assert.NoError(err)
		assert.NotNil(restMapper)

		kubernetesClient, err := clientOptions.KubernetesClient()
		assert.NoError(err)
		assert.NotNil(kubernetesClient)
	})

	t.Run("with_invalid kubeconfig", func(t *testing.T) {
		clientOptions := options.NewClientOptions(true)
		require.NotNil(clientOptions)

		flagSet := pflag.NewFlagSet("", pflag.ContinueOnError)
		clientOptions.AddFlags(flagSet)
		flagSet.Parse([]string{"--kubeconfig=xxx"})

		_, err := clientOptions.ToRESTConfig()
		assert.Error(err)
		_, err = clientOptions.ToDiscoveryClient()
		assert.Error(err)
		_, err = clientOptions.ToRESTMapper()
		assert.Error(err)
		_, err = clientOptions.KubernetesClient()
		assert.Error(err)
		_, err = clientOptions.ControllerRuntimeClient()
		assert.Error(err)
	})

	t.Run("settings_by_flag", func(t *testing.T) {
		clientOptions := options.NewClientOptions(true)
		require.NotNil(clientOptions)

		flagSet := pflag.NewFlagSet("", pflag.ContinueOnError)
		clientOptions.AddFlags(flagSet)
		flagSet.Parse([]string{
			"--kube-api-qps=40", "--kube-api-burst=100", "--field-owner=test-owner",
			"--as=test-user", "--as-group=test-group", "--as-group=test-group2",
		})

		assert.Equal("test-owner", clientOptions.FieldOwner())

		restConfig, err := clientOptions.ToRESTConfig()
		require.NoError(err)
		require.NotNil(restConfig)
		assert.EqualValues(40.0, restConfig.QPS)
		assert.Equal(100, restConfig.Burst)
		// no impersonation
		assert.Equal("", restConfig.Impersonate.UserName)
		assert.Len(restConfig.Impersonate.Groups, 0)
	})

	t.Run("settings_by_flag_with_impersonation", func(t *testing.T) {
		clientOptions := options.NewClientOptions(true).WithImpersonation()
		require.NotNil(clientOptions)

		flagSet := pflag.NewFlagSet("", pflag.ContinueOnError)
		clientOptions.AddFlags(flagSet)
		flagSet.Parse([]string{"--as=test-user", "--as-group=test-group", "--as-group=test-group2"})

		restConfig, err := clientOptions.ToRESTConfig()
		require.NoError(err)
		require.NotNil(restConfig)
		assert.Equal("test-user", restConfig.Impersonate.UserName)
		assert.Equal([]string{"test-group", "test-group2"}, restConfig.Impersonate.Groups)
	})

	t.Run("settings_by_method", func(t *testing.T) {
		clientOptions := options.NewClientOptions(true).
			WithFieldOwner("test-owner").
			WithKubeAPIQPS(40).
			WithKubeAPIBurst(100).
			WithDiscoveryBurst(200)
		require.NotNil(clientOptions)

		assert.Equal("test-owner", clientOptions.FieldOwner())
		restConfig, err := clientOptions.ToRESTConfig()
		require.NoError(err)
		require.NotNil(restConfig)
		assert.EqualValues(40.0, restConfig.QPS)
		assert.Equal(100, restConfig.Burst)
	})
}
