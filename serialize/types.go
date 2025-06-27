// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of weewoe

package serialize

import "github.com/liupzmin/weewoe/internal/render"

type Serializable interface {
	Encode(render.Rows) ([]byte, error)
	Decode([]byte) (render.Rows, error)
}
