package scrape

import (
	"context"
	"io"

	"github.com/liupzmin/weewoe/log"

	pb "github.com/liupzmin/weewoe/proto"
)

var SendMailCH = make(chan struct{})

type State struct {
}

func (s *State) SendMail(ctx context.Context, command *pb.Command) (*pb.Empty, error) {
	if command.ID == 2 {
		SendMailCH <- struct{}{}
	}
	return &pb.Empty{}, nil
}

func (s *State) DrainProcessState(stream pb.State_DrainProcessStateServer) error {
	ch := make(chan []*ProcessState, 1)
	coll := GetCollector()
	defer func() {
		coll.UnRegisterProChan(ch)
		log.Debugf("bye bye")
	}()

	//	监听数据变化
	go func() {
		coll.RegisterProChan(ch)
		for data := range ch {
			if err := stream.Send(encodeProcessData(data)); err != nil {
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
			data := encodeProcessData(coll.FetchProFromCache())
			if err := stream.Send(data); err != nil {
				return err
			}
		case Reload:
			coll.ReCollect()
		}
	}
}

func (s *State) DrainPortState(stream pb.State_DrainPortStateServer) error {
	ch := make(chan []*PortState, 1)
	coll := GetCollector()
	defer func() {
		coll.UnRegisterPortChan(ch)
		log.Debugf("bye bye")
	}()

	//	监听数据变化
	go func() {
		coll.RegisterPortChan(ch)
		for data := range ch {
			if err := stream.Send(encodePortData(data)); err != nil {
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
			data := encodePortData(coll.FetchPortFromCache())
			if err := stream.Send(data); err != nil {
				return err
			}
		case Reload:
			coll.ReCollect()
		}
	}
}
