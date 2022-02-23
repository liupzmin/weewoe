package model_test

import (
	"testing"

	"github.com/liupzmin/weewoe/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestTruncate(t *testing.T) {
	uu := map[string]struct {
		data string
		size int
		e    string
	}{
		"same": {
			data: "fred",
			size: 4,
			e:    "fred",
		},
		"small": {
			data: "fred",
			size: 10,
			e:    "fred",
		},
		"larger": {
			data: "fred",
			size: 3,
			e:    "fr…",
		},
	}

	for k := range uu {
		u := uu[k]
		t.Run(k, func(t *testing.T) {
			assert.Equal(t, u.e, model.Truncate(u.data, u.size))
		})
	}
}
