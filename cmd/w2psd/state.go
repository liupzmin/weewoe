package main

import pb "github.com/liupzmin/weewoe/proto"

type State struct {
}

func (s State) DrainProcessState(server pb.State_DrainProcessStateServer) error {
	panic("implement me")
}

func (s State) DrainPortState(server pb.State_DrainPortStateServer) error {
	panic("implement me")
}
