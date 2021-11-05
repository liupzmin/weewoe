package main

import (
	"sync/atomic"
)

var sc = &StateCache{}

type StateCache struct {
	m atomic.Value
}

func (s *StateCache) Sync(collection []*ProcessState) {
	sc := make([]*ProcessState, len(collection))
	copy(sc, collection)
	s.m.Store(sc)
}

func (s *StateCache) Fetch() []ProcessState {
	return s.m.Load().([]ProcessState)
}
