package tchart_test

import (
	"strconv"
	"testing"

	"github.com/liupzmin/weewoe/internal/tchart"
	"github.com/stretchr/testify/assert"
)

func TestDial3x3(t *testing.T) {
	d := tchart.NewDotMatrix()
	for n := 0; n <= 2; n++ {
		i := n
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			assert.Equal(t, tchart.To3x3Char(i), d.Print(i))
		})
	}
}
