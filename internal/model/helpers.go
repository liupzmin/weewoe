// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package model

import (
	"context"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/liupzmin/tview"
	runewidth "github.com/mattn/go-runewidth"
)

// Truncate a string to the given l and suffix ellipsis if needed.
func Truncate(str string, width int) string {
	return runewidth.Truncate(str, width, string(tview.SemigraphicsHorizontalEllipsis))
}

// NewExpBackOff returns a new exponential backoff timer.
func NewExpBackOff(ctx context.Context, start, max time.Duration) backoff.BackOffContext {
	bf := backoff.NewExponentialBackOff()
	bf.InitialInterval, bf.MaxElapsedTime = start, max
	return backoff.WithContext(bf, ctx)
}
