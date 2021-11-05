package main

import (
	"sync"
	"time"

	"github.com/liupzmin/weewoe/log"
)

var (
	once      sync.Once
	collector *Collector
)

type Collector struct {
	re     chan struct{}
	client []chan<- []*ProcessState
	sync.RWMutex
}

func GetCollector() *Collector {
	once.Do(func() {
		collector = &Collector{
			re: make(chan struct{}),
		}
		go collector.Collect()
	})
	return collector
}

func (c *Collector) Collect() {
	t := time.NewTicker(5 * time.Minute)

	for {
		select {
		case <-t.C:
			c.collect()
		case <-c.re:
			t.Reset(5 * time.Minute)
			c.collect()
		}
	}
}

func (c *Collector) collect() []*ProcessState {
	collection := make([]*ProcessState, 0)
	var (
		wg  sync.WaitGroup
		mux sync.Mutex
	)
	for _, v := range processInfo.Processes {
		wg.Add(len(v.Process))
		go func(pc ProcessConfig) {
			imux.RLock()
			t := instances[pc.Host]
			imux.RUnlock()

			mux.Lock()
			defer mux.Unlock()
			for _, p := range pc.Process {
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

	go func() {
		c.RLock()
		defer c.RUnlock()
		for _, ch := range c.client {
			ch <- collection
		}
	}()

	go sc.Sync(collection)

	return nil
}

func (c *Collector) FetchFromCache() []ProcessState {
	return sc.Fetch()
}

func (c *Collector) ReCollect() {
	c.re <- struct{}{}
}

func (c *Collector) RegisterChan(ch chan<- []*ProcessState) {
	c.Lock()
	defer c.Unlock()
	c.client = append(c.client, ch)
}
