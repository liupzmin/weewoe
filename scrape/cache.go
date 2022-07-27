package scrape

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/liupzmin/weewoe/internal/render"
)

const (
	RUNNING   = "RUNNING"
	EXCEPTION = "EXCEPTION"
	OPEN      = "OPEN"
	CLOSED    = "CLOSED"
)

type Group struct {
	Name      string          `json:"name"`
	Processes []*CacheProcess `json:"processes"`
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

type NameSpace struct {
	g []Group
	sync.RWMutex
}

func NewNameSpace() *NameSpace {
	return &NameSpace{}
}

func (n *NameSpace) Flat(ns string) [][]string {
	n.RLock()
	defer n.RUnlock()
	data := n.g
	var ok bool
	// todo: 优雅处理
	if ns != "" {
		if data, ok = n.Exist(ns); !ok {
			return nil
		}
	}
	top := make([][]string, 0)
	for _, v := range data {
		for _, p := range v.Processes {
			pss := cps(p.Ports).Status()
			if p.State != Good {
				pss = ""
			}
			row := []string{
				v.Name,
				p.Name,
				GetProcessStateDesc(p.State),
				p.Host,
				cps(p.Ports).String(),
				pss,
				p.Path,
				p.StartTime,
				p.TimeStamp,
				fmt.Sprintf("%t", p.Suspend),
				p.Flag,
			}
			top = append(top, row)
		}
	}
	return top
}

func (n *NameSpace) Erect(rows render.Rows) {
	gmaps := make(map[string]*Group)
	for _, v := range rows {
		if _, ok := gmaps[v.Fields[0]]; !ok {
			gmaps[v.Fields[0]] = &Group{
				Name:      v.Fields[0],
				Processes: nil,
			}
		}
		g := gmaps[v.Fields[0]]

		f := func(c rune) bool {
			return c == '|'
		}
		ports := strings.FieldsFunc(v.Fields[4], f)
		ps := strings.FieldsFunc(v.Fields[5], f)

		pss := make([]CachePort, 0)
		for i, num := range ports {
			if FromProcessStateDesc(v.Fields[2]) != Good {
				break
			}
			cp := CachePort{
				Number: num,
				State:  FromPortStateDesc(ps[i]),
			}
			pss = append(pss, cp)
		}

		pc := &CacheProcess{
			Name:      v.Fields[1],
			Host:      v.Fields[3],
			Path:      v.Fields[6],
			Flag:      v.Fields[7],
			Ports:     pss,
			State:     FromProcessStateDesc(v.Fields[2]),
			StartTime: v.Fields[8],
			TimeStamp: v.Fields[9],
			Suspend:   Suspend(v.Fields[10]),
		}
		g.Processes = append(g.Processes, pc)
	}

	n.Lock()
	defer n.Unlock()
	g := make([]Group, 0)
	for _, v := range gmaps {
		g = append(g, *v)
	}
	n.g = g
}

func (n *NameSpace) Exist(ns string) ([]Group, bool) {
	n.RLock()
	defer n.RUnlock()
	for _, v := range n.g {
		if v.Name == ns {
			return []Group{v}, true
		}
	}
	return nil, false
}

func (n *NameSpace) Suck(ns string, rows *render.Rows) {
	for _, v := range n.Flat(ns) {
		var row render.Row
		row.ID = fmt.Sprintf("%x", md5.Sum([]byte(v[0]+v[1]+v[7])))
		row.Fields = v
		*rows = append(*rows, row)
	}
}

func (n *NameSpace) SetGroups(g []Group) {
	n.Lock()
	defer n.Unlock()
	n.g = g
}

func (n *NameSpace) Groups() []Group {
	n.RLock()
	defer n.RUnlock()

	g := make([]Group, len(n.g))
	copy(g, n.g)

	return g
}

type CacheProcess struct {
	Name      string      `json:"name"`
	Host      string      `json:"host"`
	Path      string      `json:"path"`
	Flag      string      `json:"flag"`
	Ports     []CachePort `json:"ports"`
	State     int64       `json:"state"`
	StartTime string      `json:"start_time"`
	TimeStamp string      `json:"time_stamp"`
	Suspend   bool        `json:"suspend"`
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
	Number string `json:"number"`
	State  int64  `json:"state"`
}

type cps []CachePort

func (c cps) String() string {
	var s string
	for _, v := range c {
		s = fmt.Sprintf("%s|%s", s, v.Number)
	}
	return strings.Trim(s, "|")
}

func (c cps) Status() (s string) {
	for _, v := range c {
		s = fmt.Sprintf("%s|%s", s, GetPortStateDesc(v.State))
	}
	return strings.Trim(s, "|")
}

type ProcessCache struct {
	m, n atomic.Value
}

func (c *ProcessCache) Render() render.Rows {
	rows := make(render.Rows, 0)

	g := c.List()
	n := NewNameSpace()
	n.SetGroups(g)
	n.Suck("", &rows)

	return rows
}

func (c *ProcessCache) Empty() bool {
	if len(c.FetchPro()) == 0 || len(c.FetchPort()) == 0 {
		return true
	}
	return false
}

func (c *ProcessCache) List() []Group {
	return c.MergeSort(c.FetchPro(), c.FetchPort())
}

func (c *ProcessCache) SyncPro(collection []*ProcessState) {
	sc := make([]*ProcessState, len(collection))
	copy(sc, collection)
	c.m.Store(sc)
}

func (c *ProcessCache) FetchPro() []*ProcessState {
	p, ok := c.m.Load().([]*ProcessState)
	if ok {
		return p
	}
	return nil
}

func (c *ProcessCache) SyncPort(collection []*PortState) {
	sc := make([]*PortState, len(collection))
	copy(sc, collection)
	c.n.Store(sc)
}

func (c *ProcessCache) FetchPort() []*PortState {
	p, ok := c.n.Load().([]*PortState)
	if ok {
		return p
	}
	return nil
}

func (c *ProcessCache) MergeSort(pros []*ProcessState, ports []*PortState) []Group {
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
		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			loc = time.FixedZone("CST-8", 8*3600)
		}
		if v.StartTime > 0 {
			st = time.Unix(v.StartTime, 0).In(loc).Format(TimeLayout)
		}
		g[v.Name] = &CacheProcess{
			Name:      v.Name,
			Host:      v.Host,
			Path:      v.Path,
			Flag:      v.Flag,
			State:     v.State,
			StartTime: st,
			Suspend:   v.Suspend,
			TimeStamp: time.Unix(v.Timestamp, 0).In(loc).Format(TimeLayout),
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

//  Helpers

func GetProcessStateDesc(state int64) string {
	switch state {
	case Bad:
		return EXCEPTION
	case Good:
		return RUNNING
	default:
		return ""
	}
}

func FromProcessStateDesc(desc string) int64 {
	switch desc {
	case EXCEPTION:
		return Bad
	case RUNNING:
		return Good
	default:
		return Bad
	}
}

func GetPortStateDesc(state int64) string {
	switch state {
	case Bad:
		return CLOSED
	case Good:
		return OPEN
	default:
		return ""
	}
}

func FromPortStateDesc(desc string) int64 {
	switch desc {
	case CLOSED:
		return Bad
	case OPEN:
		return Good
	default:
		return Bad
	}
}

func Suspend(s string) bool {
	if s == "true" {
		return true
	}
	return false
}
