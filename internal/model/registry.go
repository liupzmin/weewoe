// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package model

import (
	"github.com/liupzmin/weewoe/internal/render"
)

// Registry tracks resources metadata.
var Registry = map[string]ResourceMeta{
	"process": {
		Renderer: &render.Process{},
	},
	"namespace": {
		Renderer: &render.Namespace{},
	},
}
