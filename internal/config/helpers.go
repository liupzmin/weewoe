// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package config

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

const (
	// DefaultDirMod default unix perms for k9s directory.
	DefaultDirMod os.FileMode = 0755
	// DefaultFileMod default unix perms for k9s files.
	DefaultFileMod os.FileMode = 0600
)

// InList check if string is in a collection of strings.
func InList(ll []string, n string) bool {
	for _, l := range ll {
		if l == n {
			return true
		}
	}
	return false
}

// MustW2User establishes current user identity or fail.
func MustW2User() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal().Err(err).Msg("Die on retrieving user info")
	}
	return usr.Username
}

// EnsurePath ensures a directory exist from the given path.
func EnsurePath(path string, mod os.FileMode) {
	dir := filepath.Dir(path)
	EnsureFullPath(dir, mod)
}

// EnsureFullPath ensures a directory exist from the given path.
func EnsureFullPath(path string, mod os.FileMode) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, mod); err != nil {
			log.Fatal().Msgf("Unable to create dir %q %v", path, err)
		}
	}
}

// IsBoolSet checks if a bool prt is set.
func IsBoolSet(b *bool) bool {
	return b != nil && *b
}
