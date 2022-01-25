package render_test

import (
	"testing"

	"github.com/liupzmin/weewoe/internal/render"
	"github.com/stretchr/testify/assert"
)

func TestRoleBindingRender(t *testing.T) {
	c := render.RoleBinding{}
	r := render.NewRow(6)
	c.Render(load(t, "rb"), "", &r)

	assert.Equal(t, "default/blee", r.ID)
	assert.Equal(t, render.Fields{"default", "blee", "blee", "SvcAcct", "fernand"}, r.Fields[:5])
}
