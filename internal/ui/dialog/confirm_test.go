// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package dialog

import (
	"testing"

	"github.com/liupzmin/tview"
	"github.com/liupzmin/weewoe/internal/config"
	"github.com/liupzmin/weewoe/internal/ui"
	"github.com/stretchr/testify/assert"
)

func TestConfirmDialog(t *testing.T) {
	a := tview.NewApplication()
	p := ui.NewPages()
	a.SetRoot(p, false)

	ackFunc := func() {
		assert.True(t, true)
	}
	caFunc := func() {
		assert.True(t, true)
	}
	ShowConfirm(config.Dialog{}, p, "Blee", "Yo", ackFunc, caFunc)

	d := p.GetPrimitive(confirmKey).(*tview.ModalForm)
	assert.NotNil(t, d)

	dismissConfirm(p)
	assert.Nil(t, p.GetPrimitive(confirmKey))
}
