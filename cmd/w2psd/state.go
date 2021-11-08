package main

import (
	"io"

	"github.com/liupzmin/weewoe/log"

	pb "github.com/liupzmin/weewoe/proto"
)

type State struct {
}

func (s *State) DrainProcessState(stream pb.State_DrainProcessStateServer) error {
	ch := make(chan []*ProcessState)
	coll := GetCollector()
	defer func() {
		coll.UnRegisterChan(ch)
		close(ch)
		log.Debugf("bye bye")
	}()

	//	监听数据变化
	go func() {
		coll.RegisterChan(ch)
		for data := range ch {
			if err := stream.Send(s.buildProcessData(data)); err != nil {
				log.Errorf("stream send error: %s", err.Error())
			}
		}
	}()
	// 接受指令
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		switch in.ID {
		case UseCache:
			data := s.buildProcessData(coll.FetchFromCache())
			if err := stream.Send(data); err != nil {
				return err
			}
		case Reload:
			coll.ReCollect()
		}
	}
}

func (s *State) DrainPortState(stream pb.State_DrainPortStateServer) error {
	panic("implement me")
}

func (s *State) buildProcessData(data []*ProcessState) *pb.ProcessCollection {
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
