package client

import (
	"math"
	"strconv"
	"sync"
)

var defCat = &Cat{
	buf: make([]string, 0),
}

type Cat struct {
	buf []string
	sync.RWMutex
}

func (c *Cat) All() []string {
	c.RLock()
	defer c.RUnlock()
	r := make([]string, len(c.buf))
	copy(r, c.buf)
	return r
}

func (c *Cat) Add(cat string) {
	c.Lock()
	defer c.Unlock()
	c.buf = append(c.buf, cat)
}

func ListCats() []string {
	return defCat.All()
}

func AddCat(cat string) {
	defCat.Add(cat)
}

// ----------------------------------------------------------------------------
// Helpers...

// MegaByte represents a megabyte.
const MegaByte = 1024 * 1024

// ToMB converts bytes to megabytes.
func ToMB(v int64) int64 {
	return v / MegaByte
}

// ToPercentage computes percentage as string otherwise n/aa.
func ToPercentage(v1, v2 int64) int {
	if v2 == 0 {
		return 0
	}
	return int(math.Floor((float64(v1) / float64(v2)) * 100))
}

// ToPercentageStr computes percentage, but if v2 is 0, it will return NAValue instead of 0.
func ToPercentageStr(v1, v2 int64) string {
	if v2 == 0 {
		return NA
	}
	return strconv.Itoa(ToPercentage(v1, v2))
}
