package serialize

import "github.com/liupzmin/weewoe/internal/render"

type Serializable interface {
	Encode(render.Rows) ([]byte, error)
	Decode([]byte) (render.Rows, error)
}
