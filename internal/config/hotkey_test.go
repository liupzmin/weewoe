// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package config_test

import (
	"testing"

	"github.com/liupzmin/weewoe/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestHotKeyLoad(t *testing.T) {
	h := config.NewHotKeys()
	assert.Nil(t, h.LoadHotKeys("testdata/hot_key.yml"))

	assert.Equal(t, 1, len(h.HotKey))

	k, ok := h.HotKey["pods"]
	assert.True(t, ok)
	assert.Equal(t, "shift-0", k.ShortCut)
	assert.Equal(t, "Launch pod view", k.Description)
	assert.Equal(t, "pods", k.Command)
}
