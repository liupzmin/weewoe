// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/liupzmin/weewoe/internal/model"

	"github.com/gdamore/tcell/v2"
	"github.com/liupzmin/tview"
	"github.com/liupzmin/weewoe/internal/config"
	"github.com/liupzmin/weewoe/internal/render"
)

// StatusIndicator represents a status indicator when main header is collapsed.
type StatusIndicator struct {
	*tview.TextView

	app       *App
	styles    *config.Styles
	permanent string
	cancel    context.CancelFunc
}

// NewStatusIndicator returns a new status indicator.
func NewStatusIndicator(app *App, styles *config.Styles) *StatusIndicator {
	s := StatusIndicator{
		TextView: tview.NewTextView(),
		app:      app,
		styles:   styles,
	}
	s.SetTextAlign(tview.AlignCenter)
	s.SetTextColor(tcell.ColorWhite)
	s.SetBackgroundColor(styles.BgColor())
	s.SetDynamicColors(true)
	styles.AddListener(&s)

	return &s
}

// StylesChanged notifies the skins changed.
func (s *StatusIndicator) StylesChanged(styles *config.Styles) {
	s.styles = styles
	s.SetBackgroundColor(styles.BgColor())
	s.SetTextColor(styles.FgColor())
}

const statusIndicatorFmt = "[orange::b]weewoe [aqua::]%d/[lawngreen::]%d"

// ClusterInfoChanged notifies the cluster meta was changed.
func (s *StatusIndicator) ClusterInfoChanged(c model.Cluster) {
	if !s.app.IsRunning() {
		return
	}
	s.app.QueueUpdateDraw(func() {
		s.SetPermanent(fmt.Sprintf(
			statusIndicatorFmt,
			c.Total,
			c.Healthy,
		))
	})
}

// SetPermanent sets permanent title to be reset to after updates.
func (s *StatusIndicator) SetPermanent(info string) {
	s.permanent = info
	s.SetText(info)
}

// Reset clears out the logo view and resets colors.
func (s *StatusIndicator) Reset() {
	s.Clear()
	s.SetPermanent(s.permanent)
}

// Err displays a log error state.
func (s *StatusIndicator) Err(msg string) {
	s.update(msg, "orangered")
}

// Warn displays a log warning state.
func (s *StatusIndicator) Warn(msg string) {
	s.update(msg, "mediumvioletred")
}

// Info displays a log info state.
func (s *StatusIndicator) Info(msg string) {
	s.update(msg, "lawngreen")
}

func (s *StatusIndicator) update(msg, c string) {
	s.setText(fmt.Sprintf("[%s::b] <%s> ", c, msg))
}

func (s *StatusIndicator) setText(msg string) {
	if s.cancel != nil {
		s.cancel()
	}
	s.SetText(msg)

	var ctx context.Context
	ctx, s.cancel = context.WithCancel(context.Background())
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			s.app.QueueUpdateDraw(func() {
				s.Reset()
			})
		}
	}(ctx)
}

// Helpers...

// AsPercDelta represents a percentage with a delta indicator.
func AsPercDelta(ov, nv int) string {
	prev, cur := render.IntToStr(ov), render.IntToStr(nv)
	return cur + "%" + Deltas(prev, cur)
}
