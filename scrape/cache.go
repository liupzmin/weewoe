package scrape

import (
	"fmt"
	"sort"
	"sync/atomic"
	"time"

	"github.com/liupzmin/weewoe/log"
)

type Group struct {
	Name      string
	Processes []*CacheProcess
}

func (g Group) CountNormal() int {
	var count int
	for _, p := range g.Processes {
		if p.State == Good {
			count++
		}
	}
	return count
}

func (g Group) CountSuspend() int {
	var count int
	for _, p := range g.Processes {
		if p.Suspend {
			count++
		}
	}
	return count
}

func (g Group) CountPortException() int {
	var count int
	for _, p := range g.Processes {
		if p.State == Good {
			for _, v := range p.Ports {
				if v.State == 0 {
					count++
				}
			}
		}
	}
	return count
}

func (g Group) CountProcessException() int {
	var count int
	for _, p := range g.Processes {
		if p.State == Bad && !p.Suspend {
			count++
		}
	}
	return count
}

type CacheProcess struct {
	Name      string
	Host      string
	Ports     []CachePort
	State     int64
	StartTime string
	Suspend   bool
}

func (cp *CacheProcess) String() string {
	return fmt.Sprintf("CacheProcess{Name: %s,Host:%s,Ports:%v,State:%d,StartTime:%s,Suspend:%t}",
		cp.Name,
		cp.Host,
		cp.Ports,
		cp.State,
		cp.StartTime,
		cp.Suspend)
}

type CachePort struct {
	Number string
	State  int64
}

var SC = &StateCache{
	mch: make(chan []*ProcessState),
	nch: make(chan []*PortState),
}

type StateCache struct {
	m, n atomic.Value
	mch  chan []*ProcessState
	nch  chan []*PortState
}

func (s *StateCache) SyncPro(collection []*ProcessState) {
	go func() {
		s.mch <- collection
	}()
	sc := make([]*ProcessState, len(collection))
	copy(sc, collection)
	s.m.Store(sc)
	log.Debugf("process cache sync")
}

func (s *StateCache) FetchPro() []*ProcessState {
	p, ok := s.m.Load().([]*ProcessState)
	if ok {
		return p
	}
	return nil
}

func (s *StateCache) SyncPort(collection []*PortState) {
	go func() {
		s.nch <- collection
	}()
	sc := make([]*PortState, len(collection))
	copy(sc, collection)
	s.n.Store(sc)
	log.Debugf("port cache sync")
}

func (s *StateCache) FetchPort() []*PortState {
	p, ok := s.n.Load().([]*PortState)
	if ok {
		return p
	}
	return nil
}

func (s *StateCache) MergeSort(pros []*ProcessState, ports []*PortState) []Group {
	if len(pros) != len(ports) {
		return nil
	}
	gmap := make(map[string]map[string]*CacheProcess)
	for _, v := range pros {
		g, ok := gmap[v.Group]
		if !ok {
			g = make(map[string]*CacheProcess, 0)
		}
		var st string
		if v.StartTime > 0 {
			st = time.Unix(v.StartTime, 0).Format(TimeLayout)
		}
		g[v.Name] = &CacheProcess{
			Name:      v.Name,
			Host:      v.Host,
			State:     v.State,
			StartTime: st,
			Suspend:   v.Suspend,
		}
		gmap[v.Group] = g
	}

	for _, v := range ports {
		pts := make([]CachePort, 0)
		for _, pt := range v.States {
			port := CachePort{
				Number: pt.Number,
				State:  pt.State,
			}
			pts = append(pts, port)
		}
		p := gmap[v.Group][v.Name]
		p.Ports = pts
	}

	gs := make([]Group, 0)
	for k, v := range gmap {
		ps := make([]*CacheProcess, 0)
		for _, v2 := range v {
			ps = append(ps, v2)
		}
		sort.Slice(ps, func(i, j int) bool {
			return ps[i].Name < ps[j].Name
		})
		g := Group{
			Name:      k,
			Processes: ps,
		}
		gs = append(gs, g)
	}

	sort.Slice(gs, func(i, j int) bool {
		return gs[i].Name < gs[j].Name
	})

	return gs
}
