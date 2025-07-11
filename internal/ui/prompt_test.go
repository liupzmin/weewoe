// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package ui_test

import (
	"testing"

	"github.com/liupzmin/weewoe/internal/config"
	"github.com/liupzmin/weewoe/internal/model"
	"github.com/liupzmin/weewoe/internal/ui"
	"github.com/stretchr/testify/assert"
)

func TestCmdNew(t *testing.T) {
	v := ui.NewPrompt(nil, true, config.NewStyles())
	model := model.NewFishBuff(':', model.CommandBuffer)
	v.SetModel(model)
	model.AddListener(v)
	for _, r := range "blee" {
		model.Add(r)
	}

	assert.Equal(t, "\x00> [::b]blee\n", v.GetText(false))
}

func TestCmdUpdate(t *testing.T) {
	model := model.NewFishBuff(':', model.CommandBuffer)
	v := ui.NewPrompt(nil, true, config.NewStyles())
	v.SetModel(model)

	model.AddListener(v)
	model.SetText("blee", "")
	model.Add('!')

	assert.Equal(t, "\x00> [::b]blee!\n", v.GetText(false))
	assert.False(t, v.InCmdMode())
}

func TestCmdMode(t *testing.T) {
	model := model.NewFishBuff(':', model.CommandBuffer)
	v := ui.NewPrompt(&ui.App{}, true, config.NewStyles())
	v.SetModel(model)
	model.AddListener(v)

	for _, f := range []bool{false, true} {
		model.SetActive(f)
		assert.Equal(t, f, v.InCmdMode())
	}
}
