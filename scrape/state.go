package scrape

import (
	"context"
	"fmt"

	"github.com/spf13/viper"

	"github.com/liupzmin/weewoe/serialize"

	"github.com/liupzmin/weewoe/log"

	pb "github.com/liupzmin/weewoe/proto"
)

var (
	globalCtx    context.Context
	globalCancel context.CancelFunc
)

func init() {
	globalCtx, globalCancel = context.WithCancel(context.Background())
}

func Stop() {
	globalCancel()
	for _, v := range CollectorMap {
		if v.Running() {
			v.Stop()
		}
	}
}

type State struct{}

func (s *State) GetDomain(ctx context.Context, empty *pb.Empty) (*pb.Domain, error) {
	name := viper.GetString("domain.name")
	return &pb.Domain{
		Name: name,
	}, nil
}

func (s *State) Drain(kind *pb.Kind, stream pb.State_DrainServer) error {
	var (
		c  Collector
		ok bool
	)
	if c, ok = CollectorMap[kind.Name]; !ok {
		return fmt.Errorf("not surpported collector")
	}

	if err := c.Start(); err != nil {
		return err
	}

	w := WrapperFuncMap[kind.Name](c, serialize.ProcessGob{})
	c.AddListener(w)
	defer c.RemoveListener(w)

	for {
		select {
		case data := <-w.Chan():
			if err := stream.Send(&pb.Data{Content: data}); err != nil {
				log.Errorf("stream send error: %s", err.Error())
				return err
			}
		case <-stream.Context().Done():
			log.Info("stream done!")
			return nil
		case <-globalCtx.Done():
			log.Info("server done!")
			return nil
		}
	}
}

func (s *State) SendCommand(ctx context.Context, command *pb.Command) (*pb.Empty, error) {
	var (
		c  Collector
		ok bool
	)
	if c, ok = CollectorMap[command.Kind]; !ok {
		return &pb.Empty{}, fmt.Errorf("not surpported collector")
	}
	c.SendCommand(command.ID)
	return &pb.Empty{}, nil
}
