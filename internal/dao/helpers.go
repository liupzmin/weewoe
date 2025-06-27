// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package dao

import (
	"math"
	"regexp"

	"github.com/liupzmin/tview"
	runewidth "github.com/mattn/go-runewidth"
)

var (
	inverseRx = regexp.MustCompile(`\A\!`)
	fuzzyRx   = regexp.MustCompile(`\A\-f`)
)

// IsInverseSelector checks if inverse char has been provided.
func IsInverseSelector(s string) bool {
	if s == "" {
		return false
	}
	return inverseRx.MatchString(s)
}

// IsFuzzySelector checks if filter is fuzzy or not.
func IsFuzzySelector(s string) bool {
	if s == "" {
		return false
	}
	return fuzzyRx.MatchString(s)
}

func toPerc(v1, v2 float64) float64 {
	if v2 == 0 {
		return 0
	}
	return math.Round((v1 / v2) * 100)
}

// Truncate a string to the given l and suffix ellipsis if needed.
func Truncate(str string, width int) string {
	return runewidth.Truncate(str, width, string(tview.SemigraphicsHorizontalEllipsis))
}
