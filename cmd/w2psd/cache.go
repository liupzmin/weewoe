package main

import (
	"sync/atomic"

	"github.com/liupzmin/weewoe/log"
)

var sc = &StateCache{}

type StateCache struct {
	m atomic.Value
}

func (s *StateCache) Sync(collection []*ProcessState) {
	sc := make([]*ProcessState, len(collection))
	copy(sc, collection)
	s.m.Store(sc)
	log.Debugf("cache sync")
}

func (s *StateCache) Fetch() []*ProcessState {
	return s.m.Load().([]*ProcessState)
}
