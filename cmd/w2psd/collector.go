package main

import (
	"container/list"
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
	client *list.List
	sync.RWMutex
}

func GetCollector() *Collector {
	return collector
}

func (c *Collector) Collect() {
	t := time.NewTicker(5 * time.Minute)
	time.Sleep(5 * time.Second)
	c.collect()
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
	log.Debugf("tick! tick! it's time to collect data")
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

	log.Debug("the new collection ", log.Any("collection", collection))

	go func() {
		c.RLock()
		defer c.RUnlock()
		for i := c.client.Front(); i != nil; i = i.Next() {
			i.Value.(chan []*ProcessState) <- collection
		}
	}()

	go sc.Sync(collection)

	return nil
}

func (c *Collector) FetchFromCache() []*ProcessState {
	return sc.Fetch()
}

func (c *Collector) ReCollect() {
	c.re <- struct{}{}
}

func (c *Collector) RegisterChan(ch chan<- []*ProcessState) {
	c.Lock()
	defer c.Unlock()
	c.client.PushBack(ch)
}

func (c *Collector) UnRegisterChan(ch chan<- []*ProcessState) {
	log.Debugf("Unregister channel")
	c.Lock()
	defer c.Unlock()
	for i := c.client.Front(); i != nil; i = i.Next() {
		if i.Value.(chan<- []*ProcessState) == ch {
			c.client.Remove(i)
		}
	}
}

func init() {
	once.Do(func() {
		collector = &Collector{
			client: list.New(),
			re:     make(chan struct{}),
		}
		go collector.Collect()
	})
}
