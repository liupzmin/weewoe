package render_test

import (
	"testing"

	"github.com/liupzmin/weewoe/internal/render"
	"github.com/stretchr/testify/assert"
)

func TestRoleRender(t *testing.T) {
	c := render.Role{}
	r := render.NewRow(3)
	c.Render(load(t, "ro"), "", &r)

	assert.Equal(t, "default/blee", r.ID)
	assert.Equal(t, render.Fields{"default", "blee"}, r.Fields[:2])
}
