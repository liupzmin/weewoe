// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of weewoe

//go:build !linux
// +build !linux

package rotate

import (
	"os"
)

func chown(_ string, _ os.FileInfo) error {
	return nil
}
