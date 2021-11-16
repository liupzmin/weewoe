package scrape

import (
	"reflect"
	"testing"
)

var (
	ps = []*ProcessState{
		{
			Process: Process{
				Name:    "B",
				Path:    "",
				Ports:   []int64{2001},
				Group:   "BROKER",
				Host:    "1.1.1.1",
				PIDFile: "",
				Flag:    "",
				Suspend: false,
			},
			State:         1,
			StateDescribe: "S",
			StartTime:     1636700045,
			Timestamp:     1636700045,
		},
		{
			Process: Process{
				Name:    "A",
				Path:    "",
				Ports:   []int64{2001},
				Group:   "BROKER",
				Host:    "1.1.1.1",
				PIDFile: "",
				Flag:    "",
				Suspend: false,
			},
			State:         1,
			StateDescribe: "R",
			StartTime:     1636700045,
			Timestamp:     1636700045,
		},
		{
			Process: Process{
				Name:    "F",
				Path:    "",
				Ports:   []int64{2001},
				Group:   "wechat",
				Host:    "1.1.1.1",
				PIDFile: "",
				Flag:    "",
				Suspend: false,
			},
			State:         0,
			StateDescribe: "",
			StartTime:     0,
			Timestamp:     1636700045,
		},
	}
	ports = []*PortState{
		{
			Process: Process{
				Name:    "F",
				Path:    "",
				Ports:   nil,
				Group:   "wechat",
				Host:    "",
				PIDFile: "",
				Flag:    "",
				Suspend: false,
			},
			States: []*Port{
				{
					Number: "20001",
					State:  0,
				},
			},
			Timestamp: 0,
		},
		{
			Process: Process{
				Name:    "B",
				Path:    "",
				Ports:   nil,
				Group:   "BROKER",
				Host:    "",
				PIDFile: "",
				Flag:    "",
				Suspend: false,
			},
			States: []*Port{
				{
					Number: "20001",
					State:  1,
				},
			},
			Timestamp: 0,
		},
		{
			Process: Process{
				Name:    "A",
				Path:    "",
				Ports:   nil,
				Group:   "BROKER",
				Host:    "",
				PIDFile: "",
				Flag:    "",
				Suspend: false,
			},
			States: []*Port{
				{
					Number: "20001",
					State:  1,
				},
			},
			Timestamp: 0,
		},
	}

	expect = []Group{
		{
			Name: "BROKER",
			Processes: []*CacheProcess{
				{
					Name: "A",
					Host: "1.1.1.1",
					Ports: []CachePort{
						{
							Number: "20001",
							State:  1,
						},
					},
					State:     1,
					StartTime: "2021-11-12 14:54:05",
					Suspend:   false,
				},
				{
					Name: "B",
					Host: "1.1.1.1",
					Ports: []CachePort{
						{
							Number: "20001",
							State:  1,
						},
					},
					State:     1,
					StartTime: "2021-11-12 14:54:05",
					Suspend:   false,
				},
			},
		},
		{
			Name: "wechat",
			Processes: []*CacheProcess{
				{
					Name: "F",
					Host: "1.1.1.1",
					Ports: []CachePort{
						{
							Number: "20001",
							State:  0,
						},
					},
					State:     0,
					StartTime: "",
					Suspend:   false,
				},
			},
		},
	}
)

func TestStateCache_MergeSort(t *testing.T) {
	SC.SyncPro(ps)
	SC.SyncPort(ports)
	data := SC.MergeSort(SC.FetchPro(), SC.FetchPort())
	t.Logf("output:%+v", data)
	t.Logf("output:%+v", expect)
	if !reflect.DeepEqual(data, expect) {
		t.Errorf("merge sort failed")
	}
}
