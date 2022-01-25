package scrape

import "github.com/liupzmin/weewoe/log"

var (
	p            = NewProcessDetail()
	CollectorMap = map[string]Collector{
		"process":   p,
		"namespace": p,
	}
	WrapperFuncMap = map[string]WrapperFunc{
		"process":   NewWrapper,
		"namespace": NewNSWrapper,
	}
)

func takeOff() {
	for k, v := range CollectorMap {
		if err := v.Start(); err != nil {
			log.Panicf("Collector %s start failed: %s", k, v)
		}
	}
}
