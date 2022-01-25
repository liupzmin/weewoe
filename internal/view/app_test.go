package view_test

import (
	"testing"

	"github.com/liupzmin/weewoe/internal/config"
	"github.com/liupzmin/weewoe/internal/view"
	"github.com/stretchr/testify/assert"
)

func TestAppNew(t *testing.T) {
	a := view.NewApp(config.NewConfig(ks{}))
	a.Init("blee", 10)

	assert.Equal(t, 11, len(a.GetActions()))
}
