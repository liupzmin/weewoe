package client

import "sync"

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
