// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package options_test

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/component-base/featuregate"

	"github.com/mesosphere/dkp-cli-runtime/extensions/options"
)

func TestFeatureGateOptions(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	var opts *options.FeatureGateOptions

	configureWithFlags := func(flags ...string) {
		t.Helper()
		var err error
		opts, err = options.NewFeatureGateOptions(map[featuregate.Feature]featuregate.FeatureSpec{
			featuregate.Feature("feature1"): {PreRelease: featuregate.GA, Default: true},
			featuregate.Feature("feature2"): {PreRelease: featuregate.Beta},
			featuregate.Feature("feature3"): {PreRelease: featuregate.Alpha},
			featuregate.Feature("feature4"): {PreRelease: featuregate.Deprecated},
		})
		require.NoError(err)
		flagSet := pflag.NewFlagSet("", pflag.ContinueOnError)
		opts.AddFlag(flagSet)
		assert.True(flagSet.HasAvailableFlags())
		err = flagSet.Parse(flags)
		require.NoError(err)
	}

	t.Run("default", func(t *testing.T) {
		configureWithFlags()
		assert.True(opts.Enabled("feature1"))
		assert.False(opts.Enabled("feature2"))
		assert.False(opts.Enabled("feature3"))
		assert.False(opts.Enabled("feature4"))

		knownFeatures := opts.KnownFeatures()
		// GA and Deprecated features not listed
		assert.Len(knownFeatures, 4)
		assert.Contains(knownFeatures[0], "AllAlpha")
		assert.Contains(knownFeatures[1], "AllBeta")
		assert.Contains(knownFeatures[2], "feature2")
		assert.Contains(knownFeatures[3], "feature3")
	})

	t.Run("individual", func(t *testing.T) {
		configureWithFlags("--feature-gates=feature2=true,feature3=true")
		assert.True(opts.Enabled("feature1"))
		assert.True(opts.Enabled("feature2"))
		assert.True(opts.Enabled("feature3"))
		assert.False(opts.Enabled("feature4"))
	})

	t.Run("levels", func(t *testing.T) {
		configureWithFlags("--feature-gates=AllBeta=true")
		assert.True(opts.Enabled("feature1"))
		assert.True(opts.Enabled("feature2"))
		assert.False(opts.Enabled("feature3"))
		assert.False(opts.Enabled("feature4"))
	})

	t.Run("no feature gates", func(t *testing.T) {
		var err error
		opts, err = options.NewFeatureGateOptions(nil)
		require.NoError(err)
		flagSet := pflag.NewFlagSet("", pflag.ContinueOnError)
		opts.AddFlag(flagSet)
		assert.False(flagSet.HasAvailableFlags(), "flag should be hidden")
	})
}
