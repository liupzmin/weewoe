package render_test

import (
	"testing"

	"github.com/liupzmin/weewoe/internal/render"
	"github.com/stretchr/testify/assert"
)

func TestStorageClassRender(t *testing.T) {
	c := render.StorageClass{}
	r := render.NewRow(4)
	c.Render(load(t, "sc"), "", &r)

	assert.Equal(t, "-/standard", r.ID)
	assert.Equal(t, render.Fields{"standard", "kubernetes.io/gce-pd"}, r.Fields[:2])
}
