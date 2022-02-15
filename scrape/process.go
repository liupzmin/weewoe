package scrape

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/liupzmin/weewoe/internal/render"

	"github.com/spf13/viper"

	"github.com/liupzmin/weewoe/log"
)

var (
	SendMail = make(chan struct{})
)

type ProcessDetail struct {
	*ProcessCache
	running   bool
	listeners []TableListener
	re, done  chan struct{}
	sync.RWMutex
}

func NewProcessDetail() *ProcessDetail {
	p := &ProcessDetail{
		ProcessCache: &ProcessCache{},
		listeners:    make([]TableListener, 0),
		re:           make(chan struct{}),
		done:         make(chan struct{}),
	}
	return p
}

func (p *ProcessDetail) Start() error {
	if p.running {
		return nil
	}
	interval := viper.GetInt("scrape.interval")
	if interval == 0 {
		return fmt.Errorf("the scrape interval can't be zero")
	}
	du := time.Duration(interval)

	go p.collect(du)
	p.running = true
	return nil
}

func (p *ProcessDetail) Refresh() {
	p.re <- struct{}{}
}

func (p *ProcessDetail) Running() bool {
	return p.running
}

func (p *ProcessDetail) AddListener(l TableListener) {
	p.Lock()
	p.listeners = append(p.listeners, l)
	p.Unlock()
	log.Debug("Process Collector register!")

	go l.TableDataChanged(p.Peek())
}

func (p *ProcessDetail) RemoveListener(l TableListener) {
	p.Lock()
	defer p.Unlock()

	log.Debug("begin to leave!")

	victim := -1
	for i, lis := range p.listeners {
		if lis == l {
			victim = i
			break
		}
	}

	if victim >= 0 {
		log.Debug("Process Collector left!")
		p.listeners = append(p.listeners[:victim], p.listeners[victim+1:]...)
	}
}

func (p *ProcessDetail) Peek() render.Rows {
	return p.Render()
}

func (p *ProcessDetail) Peep(rows render.Rows) render.Rows {
	return rows
}

func (p *ProcessDetail) SendCommand(i int64) {
	switch i {
	case 1:
		p.re <- struct{}{}
	case 2:
		SendMail <- struct{}{}
	default:

	}
}

func (p *ProcessDetail) collect(du time.Duration) {
	p.scrape()
	t := time.NewTicker(du * time.Minute)
	for {
		select {
		case <-t.C:
			log.Debugf("tick! tick! it's time to collect data")
			p.scrape()
		case <-p.re:
			t.Reset(du * time.Minute)
			log.Debugf("attention! attention! it's time to reload data")
			p.scrape()
		case <-p.done:
			log.Debugf("Process Collector: work down! I'm quitting.")
			for _, v := range instances {
				if v.Conn.IsValid() {
					v.Close()
				}
			}
			p.done <- struct{}{}
			return
		}
	}
}

func (p *ProcessDetail) scrape() {
	p.collectProcess()
	p.collectPort()
	p.fireDataChanged()
}

func (p *ProcessDetail) fireDataChanged() {
	p.RLock()
	defer p.RUnlock()

	t := p.Peek()
	for _, l := range p.listeners {
		l.TableDataChanged(t)
	}
}

func (p *ProcessDetail) collectProcess() {
	collection := make([]*ProcessState, 0)
	var (
		wg  sync.WaitGroup
		mux sync.Mutex
	)
	for _, v := range processInfo.Processes {
		wg.Add(1)
		go func(pc ProcessConfig) {
			defer wg.Done()
			imux.RLock()
			t := instances[pc.Host]
			imux.RUnlock()

			mux.Lock()
			defer mux.Unlock()
			for _, p := range pc.Process {
				p := p
				p.Host = pc.Host
				cmd := NewCommand(t, p)
				ps, err := cmd.GetProcessStat()
				if err != nil {
					log.Errorf("GetProcessStat failed: %s, process: %v", err.Error(), p)
					collection = append(collection, &ProcessState{
						Process:   p,
						State:     Bad,
						Timestamp: time.Now().Unix(),
					})
					continue
				}
				collection = append(collection, ps)
			}
		}(v)
	}

	wg.Wait()

	log.Debug("the new process collection ", log.Any("collection", collection))

	p.SyncPro(collection)
}

func (p *ProcessDetail) collectPort() {
	collection := make([]*PortState, 0)
	var (
		wg  sync.WaitGroup
		mux sync.Mutex
	)
	for _, v := range processInfo.Processes {
		for _, p := range v.Process {
			wg.Add(1)
			go func(p Process, host string) {
				defer wg.Done()

				p.Host = host
				states := make([]*Port, 0)
				for _, port := range p.Ports {
					isOpen := RawConnect(p.Host, port)
					states = append(states, &Port{
						Number: strconv.Itoa(int(port)),
						State:  transBool(isOpen),
					})
				}

				ps := &PortState{
					Process:   p,
					States:    states,
					Timestamp: time.Now().Unix(),
				}

				mux.Lock()
				defer mux.Unlock()
				collection = append(collection, ps)
			}(p, v.Host)
		}
	}

	wg.Wait()

	log.Debug("the new port collection ", log.Any("collection", collection))

	p.SyncPort(collection)
}

func (p *ProcessDetail) Stop() {
	p.running = false
	p.done <- struct{}{}
	<-p.done
}

func transBool(b bool) int64 {
	if b {
		return 1
	} else {
		return 0
	}
}