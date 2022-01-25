package serialize

import (
	"bytes"
	"encoding/gob"

	"github.com/liupzmin/weewoe/internal/render"
)

type ProcessGob struct{}

func (g ProcessGob) Encode(data render.Rows) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	return buf.Bytes(), err
}

func (g ProcessGob) Decode(data []byte) (render.Rows, error) {
	buf := bytes.NewReader(data)
	dec := gob.NewDecoder(buf)

	var rows render.Rows
	err := dec.Decode(&rows)
	return rows, err
}
