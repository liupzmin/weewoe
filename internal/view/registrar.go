// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package view

import "github.com/liupzmin/weewoe/internal/client"

func loadCustomViewers() MetaViewers {
	m := make(MetaViewers, 30)
	coreViewers(m)

	return m
}

func coreViewers(vv MetaViewers) {
	// todo: 参数问题
	vv["process"] = MetaViewer{
		viewerFn: NewProcess,
	}
	client.AddCat("process")

	vv["namespace"] = MetaViewer{
		viewerFn: NewNamespace,
	}
	client.AddCat("namespace")
}
