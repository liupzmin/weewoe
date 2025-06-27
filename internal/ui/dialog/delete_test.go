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

func TestDeleteDialog(t *testing.T) {
	p := ui.NewPages()

	okFunc := func(c, f bool) {
		assert.True(t, c)
		assert.True(t, f)
	}
	caFunc := func() {
		assert.True(t, true)
	}
	ShowDelete(config.Dialog{}, p, "Yo", okFunc, caFunc)

	d := p.GetPrimitive(deleteKey).(*tview.ModalForm)
	assert.NotNil(t, d)

	dismissDelete(p)
	assert.Nil(t, p.GetPrimitive(deleteKey))
}
