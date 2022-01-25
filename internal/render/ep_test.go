package render_test

import (
	"testing"

	"github.com/liupzmin/weewoe/internal/render"
	"github.com/stretchr/testify/assert"
)

func TestEndpointsRender(t *testing.T) {
	c := render.Endpoints{}
	r := render.NewRow(4)
	c.Render(load(t, "ep"), "", &r)

	assert.Equal(t, "default/dictionary1", r.ID)
	assert.Equal(t, render.Fields{"default", "dictionary1", "<none>"}, r.Fields[:3])
}
