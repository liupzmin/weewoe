package scrape

import pb "github.com/liupzmin/weewoe/proto"

func encodeProcessData(data []*ProcessState) *pb.ProcessCollection {
	list := make([]*pb.ProcessState, 0)
	for _, v := range data {
		ps := &pb.ProcessState{
			Base: &pb.Process{
				Name:    v.Name,
				Host:    v.Host,
				Path:    v.Path,
				Ports:   v.Ports,
				PIDFile: v.PIDFile,
				Group:   v.Group,
				Suspend: v.Suspend,
			},
			State:     v.State,
			StartTime: v.StartTime,
			Timestamp: v.Timestamp,
		}
		list = append(list, ps)
	}
	return &pb.ProcessCollection{
		List: list,
	}
}

func decodeProcessData(data *pb.ProcessCollection) []*ProcessState {
	pss := make([]*ProcessState, 0)
	for _, v := range data.List {
		ps := &ProcessState{
			Process: Process{
				Name:    v.Base.Name,
				Path:    v.Base.Path,
				Ports:   v.Base.Ports,
				Group:   v.Base.Group,
				Host:    v.Base.Host,
				PIDFile: v.Base.PIDFile,
				Suspend: v.Base.Suspend,
			},
			State:         v.State,
			StateDescribe: v.StateDescribe,
			StartTime:     v.StartTime,
			Timestamp:     v.Timestamp,
		}
		pss = append(pss, ps)
	}
	return pss
}

func encodePortData(data []*PortState) *pb.PortCollection {
	list := make([]*pb.PortState, 0)
	for _, v := range data {
		ss := make([]*pb.Port, 0)
		for _, p := range v.States {
			port := &pb.Port{
				Number: p.Number,
				State:  p.State,
			}
			ss = append(ss, port)
		}
		ps := &pb.PortState{
			Base: &pb.Process{
				Name:    v.Name,
				Host:    v.Host,
				Path:    v.Path,
				Ports:   v.Ports,
				PIDFile: v.PIDFile,
				Group:   v.Group,
				Suspend: v.Suspend,
			},
			States:    ss,
			Timestamp: v.Timestamp,
		}
		list = append(list, ps)
	}
	return &pb.PortCollection{
		List: list,
	}
}

func decodePortData(data *pb.PortCollection) []*PortState {
	pss := make([]*PortState, 0)
	for _, v := range data.List {
		states := make([]*Port, 0)
		for _, p := range v.States {
			state := &Port{
				Number: p.Number,
				State:  p.State,
			}
			states = append(states, state)
		}
		ps := &PortState{
			Process: Process{
				Name:    v.Base.Name,
				Path:    v.Base.Path,
				Ports:   v.Base.Ports,
				Group:   v.Base.Group,
				Host:    v.Base.Host,
				PIDFile: v.Base.PIDFile,
				Suspend: v.Base.Suspend,
			},
			States:    states,
			Timestamp: v.Timestamp,
		}
		pss = append(pss, ps)
	}
	return pss
}
