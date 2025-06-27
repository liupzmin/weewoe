// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// DefaultRefreshRate represents the refresh interval.
	DefaultRefreshRate = 2 // secs

	// DefaultLogLevel represents the default log level.
	DefaultLogLevel = "info"

	// DefaultCommand represents the default command to run.
	DefaultCommand = ""

	DefaultHost = "127.0.0.1"

	DefaultPort = "9527"
)

// DefaultLogFile represents the default W2 log file.
var DefaultLogFile = filepath.Join(os.TempDir(), fmt.Sprintf("w2-%s.log", MustW2User()))

// Flags represents W2 configuration flags.
type Flags struct {
	RefreshRate   *int
	Host          *string
	Port          *string
	LogLevel      *string
	LogFile       *string
	Headless      *bool
	Logoless      *bool
	Command       *string
	AllNamespaces *bool
	ReadOnly      *bool
	Write         *bool
	Crumbsless    *bool
	ScreenDumpDir *string
}

// NewFlags returns new configuration flags.
func NewFlags() *Flags {
	return &Flags{
		RefreshRate:   intPtr(DefaultRefreshRate),
		Host:          strPtr(DefaultHost),
		Port:          strPtr(DefaultPort),
		LogLevel:      strPtr(DefaultLogLevel),
		LogFile:       strPtr(DefaultLogFile),
		Headless:      boolPtr(false),
		Logoless:      boolPtr(false),
		Command:       strPtr(DefaultCommand),
		AllNamespaces: boolPtr(false),
		ReadOnly:      boolPtr(false),
		Write:         boolPtr(false),
		Crumbsless:    boolPtr(false),
		ScreenDumpDir: strPtr(W2DefaultScreenDumpDir),
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
