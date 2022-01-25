package render_test

import (
	"testing"

	"github.com/liupzmin/weewoe/internal/render"
	"github.com/stretchr/testify/assert"
)

func TestReplicaSetRender(t *testing.T) {
	c := render.ReplicaSet{}
	r := render.NewRow(4)
	c.Render(load(t, "rs"), "", &r)

	assert.Equal(t, "icx/icx-db-7d4b578979", r.ID)
	assert.Equal(t, render.Fields{"icx", "icx-db-7d4b578979", "1", "1", "1"}, r.Fields[:5])
}
