// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of weewoe

//go:build windows
// +build windows

package xcolor

import (
	"fmt"
	"math/rand"
	"strconv"
)

var _ = RandomColor()

// RandomColor generates a random color.
func RandomColor() string {
	return fmt.Sprintf("#%s", strconv.FormatInt(int64(rand.Intn(16777216)), 16))
}

// Yellow ...
func Yellow(msg string) string {
	return fmt.Sprintf("%s", msg)
}

// Red ...
func Red(msg string) string {
	return fmt.Sprintf("%s", msg)
}

// Redf ...
func Redf(msg string, arg interface{}) string {
	return fmt.Sprintf("%s %+v\n", msg, arg)
}

// Blue ...
func Blue(msg string) string {
	return fmt.Sprintf("%s", msg)
}

// Green ...
func Green(msg string) string {
	return fmt.Sprintf("%s", msg)
}

// Greenf ...
func Greenf(msg string, arg interface{}) string {
	return fmt.Sprintf("%s %+v\n", msg, arg)
}
