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
