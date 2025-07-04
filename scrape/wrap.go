// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of weewoe

package scrape

import (
	"crypto/md5"
	"fmt"
	"strconv"

	"github.com/liupzmin/weewoe/internal/render"
	"github.com/liupzmin/weewoe/log"
	"github.com/liupzmin/weewoe/serialize"
)

type WrapperFunc func(serialize.Serializable) TableListener

type Wrapper struct {
	serialize.Serializable

	ch chan []byte
}

func NewWrapper(s serialize.Serializable) TableListener {
	return &Wrapper{
		Serializable: s,
		ch:           make(chan []byte),
	}
}

func (w Wrapper) TableDataChanged(data render.Rows) {
	b, err := w.Encode(data)
	if err != nil {
		log.Errorf("Encode error: %s", err)
	}
	w.ch <- b
}

func (w Wrapper) TableLoadFailed(err error) {}

func (w Wrapper) Chan() chan []byte {
	return w.ch
}

type NSWrapper struct {
	serialize.Serializable
	ch chan []byte
}

func NewNSWrapper(s serialize.Serializable) TableListener {
	return &NSWrapper{
		Serializable: s,
		ch:           make(chan []byte),
	}
}

func (w NSWrapper) TableDataChanged(data render.Rows) {
	b, err := w.Encode(w.peep(data))
	if err != nil {
		log.Errorf("Encode error: %s", err)
	}
	w.ch <- b
}

func (w NSWrapper) TableLoadFailed(err error) {}

func (w NSWrapper) Chan() chan []byte {
	return w.ch
}

func (w *NSWrapper) peep(rows render.Rows) render.Rows {
	type ns struct {
		Name    string
		Process int
		Host    int
	}

	nmap := make(map[string]*ns)
	hmap := make(map[string]map[string]struct{})
	for _, v := range rows {
		if _, ok := nmap[v.Fields[0]]; !ok {
			nmap[v.Fields[0]] = &ns{Name: v.Fields[0]}
			hmap[v.Fields[0]] = make(map[string]struct{})
		}
		hmap[v.Fields[0]][v.Fields[3]] = struct{}{}
		nmap[v.Fields[0]].Process++
		nmap[v.Fields[0]].Host = len(hmap[v.Fields[0]])
	}

	var nrows render.Rows
	for _, v := range nmap {
		row := render.Row{
			ID:     fmt.Sprintf("%x", md5.Sum([]byte(v.Name))),
			Fields: []string{v.Name, strconv.Itoa(v.Process), strconv.Itoa(v.Host)},
		}
		nrows = append(nrows, row)
	}
	nrows.Sort(0, true, false, false)
	return nrows
}
