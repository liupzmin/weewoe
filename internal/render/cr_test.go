package render_test

import (
	"testing"

	"github.com/liupzmin/weewoe/internal/render"
	"github.com/stretchr/testify/assert"
)

func TestClusterRoleRender(t *testing.T) {
	c := render.ClusterRole{}
	r := render.NewRow(2)
	c.Render(load(t, "cr"), "-", &r)

	assert.Equal(t, "-/blee", r.ID)
	assert.Equal(t, render.Fields{"blee"}, r.Fields[:1])
}
