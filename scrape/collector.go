package scrape

import (
	"container/list"
	"strconv"
	"sync"
	"time"

	"github.com/liupzmin/weewoe/log"
)

var (
	once      sync.Once
	collector *Collector
)

type Collector struct {
	re, done        chan struct{}
	proCli, portCli *list.List
	sync.RWMutex
}

func GetCollector() *Collector {
	return collector
}

func (c *Collector) Collect() {
	t := time.NewTicker(5 * time.Minute)
	time.Sleep(5 * time.Second)
	go c.collectProcess()
	go c.collectPort()
	for {
		select {
		case <-t.C:
			log.Debugf("tick! tick! it's time to collect data")
			go c.collectProcess()
			go c.collectPort()
		case <-c.re:
			t.Reset(5 * time.Minute)
			log.Debugf("attention! attention! it's time to reload data")
			go c.collectProcess()
			go c.collectPort()
		case <-c.done:
			log.Debugf("Collector: work down! I'll quit.")
			return
		}
	}
}

func (c *Collector) collectProcess() {

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
						State:     0,
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

	go func() {
		c.RLock()
		defer c.RUnlock()
		for i := c.proCli.Front(); i != nil; i = i.Next() {
			trySendPro(i.Value.(chan<- []*ProcessState), collection)
			// i.Value.(chan<- []*ProcessState) <- collection
		}
	}()

	go SC.SyncPro(collection)
}

func (c *Collector) collectPort() {
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

	go func() {
		c.RLock()
		defer c.RUnlock()
		for i := c.portCli.Front(); i != nil; i = i.Next() {
			trySendPort(i.Value.(chan<- []*PortState), collection)
			//i.Value.(chan<- []*PortState) <- collection
		}
	}()

	go SC.SyncPort(collection)
}

func (c *Collector) FetchProFromCache() []*ProcessState {
	return SC.FetchPro()
}

func (c *Collector) FetchPortFromCache() []*PortState {
	return SC.FetchPort()
}

func (c *Collector) ReCollect() {
	c.re <- struct{}{}
}

func (c *Collector) RegisterProChan(ch chan<- []*ProcessState) {
	log.Debugf("Register process channel")
	c.Lock()
	defer c.Unlock()
	c.proCli.PushBack(ch)
}

func (c *Collector) UnRegisterProChan(ch chan<- []*ProcessState) {
	log.Debugf("Unregister process channel")
	c.Lock()
	defer c.Unlock()
	for i := c.proCli.Front(); i != nil; i = i.Next() {
		if i.Value.(chan<- []*ProcessState) == ch {
			c.proCli.Remove(i)
			close(ch)
		}
	}
}

func (c *Collector) RegisterPortChan(ch chan<- []*PortState) {
	log.Debugf("Register port channel")
	c.Lock()
	defer c.Unlock()
	c.portCli.PushBack(ch)
}

func (c *Collector) UnRegisterPortChan(ch chan<- []*PortState) {
	log.Debugf("Unregister port channel")
	c.Lock()
	defer c.Unlock()
	for i := c.portCli.Front(); i != nil; i = i.Next() {
		if i.Value.(chan<- []*PortState) == ch {
			c.portCli.Remove(i)
			close(ch)
		}
	}
}

func (c *Collector) Stop() {
	c.RLock()
	defer c.RUnlock()
	for i := c.proCli.Front(); i != nil; i = i.Next() {
		close(i.Value.(chan<- []*ProcessState))
	}
	for i := c.portCli.Front(); i != nil; i = i.Next() {
		close(i.Value.(chan<- []*PortState))
	}
	c.done <- struct{}{}
}

func trySendPro(ch chan<- []*ProcessState, data []*ProcessState) bool {
	select {
	case ch <- data:
		return true
	default:
		return false
	}
}

func trySendPort(ch chan<- []*PortState, data []*PortState) bool {
	select {
	case ch <- data:
		return true
	default:
		return false
	}
}

func transBool(b bool) int64 {
	if b {
		return 1
	} else {
		return 0
	}
}

func init() {
	once.Do(func() {
		collector = &Collector{
			proCli:  list.New(),
			portCli: list.New(),
			re:      make(chan struct{}),
			done:    make(chan struct{}),
		}
		go collector.Collect()
	})
}
