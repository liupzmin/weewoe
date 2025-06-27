// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package config

// FeatureGates represents W2 opt-in features.
type FeatureGates struct {
	NodeShell bool `yaml:"nodeShell"`
}

// NewFeatureGates returns a new feature gate.
func NewFeatureGates() *FeatureGates {
	return &FeatureGates{}
}
