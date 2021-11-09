package main

import (
	"sync/atomic"

	"github.com/liupzmin/weewoe/log"
)

var sc = &StateCache{}

type StateCache struct {
	m, n atomic.Value
}

func (s *StateCache) SyncPro(collection []*ProcessState) {
	sc := make([]*ProcessState, len(collection))
	copy(sc, collection)
	s.m.Store(sc)
	log.Debugf("process cache sync")
}

func (s *StateCache) FetchPro() []*ProcessState {
	return s.m.Load().([]*ProcessState)
}

func (s *StateCache) SyncPort(collection []*PortState) {
	sc := make([]*PortState, len(collection))
	copy(sc, collection)
	s.n.Store(sc)
	log.Debugf("port cache sync")
}

func (s *StateCache) FetchPort() []*PortState {
	return s.n.Load().([]*PortState)
}
