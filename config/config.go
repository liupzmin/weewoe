// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of weewoe

package config

import "os"

const (
	// DefaultLogLevel represents the default log level.
	DefaultLogLevel = "info"
)

var DefaultLogFile = os.Stdout.Name()

// Flags represents carefree configuration flags.
type Flags struct {
	LogLevel *string
	LogFile  *string
}

// NewFlags returns new configuration flags.
func NewFlags() *Flags {
	return &Flags{
		LogLevel: strPtr(DefaultLogLevel),
		LogFile:  strPtr(DefaultLogFile),
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}
