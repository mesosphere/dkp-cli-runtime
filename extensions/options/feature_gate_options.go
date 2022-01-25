// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"fmt"

	"github.com/spf13/pflag"
	"k8s.io/component-base/featuregate"
)

// FeatureGateOptions manages a CLIs feature gates.
type FeatureGateOptions struct {
	featureGate featuregate.FeatureGate
}

// NewFeatureGateOptions initializes a FeatureGateOptions object with the provided features.
func NewFeatureGateOptions(features map[featuregate.Feature]featuregate.FeatureSpec) (*FeatureGateOptions, error) {
	featureGate := featuregate.NewFeatureGate()

	if err := featureGate.Add(features); err != nil {
		return nil, fmt.Errorf("failed to create feature gate options: %w", err)
	}

	return &FeatureGateOptions{
		featureGate: featureGate,
	}, nil
}

// AddFlag adds flags for enabling feature gates to the provided FlagSet.
func (o *FeatureGateOptions) AddFlag(fs *pflag.FlagSet) {
	o.featureGate.(featuregate.MutableFeatureGate).AddFlag(fs)
	if len(o.featureGate.KnownFeatures()) <= 2 {
		_ = fs.MarkHidden("feature-gates")
	}
}

// Enabled answers wheather a feature is enabled.
func (o *FeatureGateOptions) Enabled(key string) bool {
	return o.featureGate.Enabled(featuregate.Feature(key))
}

// KnownFeatures returns a list of text descriptions of all known features.
func (o *FeatureGateOptions) KnownFeatures() []string {
	return o.featureGate.KnownFeatures()
}
