// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package config

const (
	// SeverityLow tracks low severity.
	SeverityLow SeverityLevel = iota

	// SeverityMedium tracks medium severity level.
	SeverityMedium

	// SeverityHigh tracks high severity level.
	SeverityHigh
)

// SeverityLevel tracks severity levels.
type SeverityLevel int

// Severity tracks a resource severity levels.
type Severity struct {
	Critical int `yaml:"critical"`
	Warn     int `yaml:"warn"`
}

// NewSeverity returns a new instance.
func NewSeverity() *Severity {
	return &Severity{
		Critical: 90,
		Warn:     70,
	}
}

// Validate checks all thresholds and make sure we're cool. If not use defaults.
func (s *Severity) Validate() {
	norm := NewSeverity()
	if !validateRange(s.Warn) {
		s.Warn = norm.Warn
	}
	if !validateRange(s.Critical) {
		s.Critical = norm.Critical
	}
}

func validateRange(v int) bool {
	if v <= 0 || v > 100 {
		return false
	}
	return true
}
