// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package ui_test

import (
	"testing"

	"github.com/liupzmin/weewoe/internal/model"
	"github.com/liupzmin/weewoe/internal/ui"
	"github.com/stretchr/testify/assert"
)

func TestKeyActionsHints(t *testing.T) {
	kk := ui.KeyActions{
		ui.KeyF: ui.NewKeyAction("fred", nil, true),
		ui.KeyB: ui.NewKeyAction("blee", nil, true),
		ui.KeyZ: ui.NewKeyAction("zorg", nil, false),
	}

	hh := kk.Hints()

	assert.Equal(t, 3, len(hh))
	assert.Equal(t, model.MenuHint{Mnemonic: "b", Description: "blee", Visible: true}, hh[0])
}
