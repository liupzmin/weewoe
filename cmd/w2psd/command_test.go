package main

import (
	"testing"

	"github.com/liupzmin/weewoe/ssh"
)

func TestGetPID(t *testing.T) {
	conn, err := ssh.NewConnection("192.168.0.127:22", "root", false)
	if err != nil {
		t.Errorf("new connnection error: %v", err)
	}

	bt, err := getBootTime(conn)
	if err != nil {
		t.Errorf("get boot time failed: %s", err.Error())
	}

	cmd := NewCommand(Target{
		Conn:     conn,
		BootTime: bt,
	}, Process{Flag: "rcu,gp,par"})
	pid, err := cmd.GetPID()
	if err != nil {
		t.Errorf("som error happend:%s", err.Error())
	}

	t.Logf("pid is :%s", pid)
}

func TestCommand_GetProcessStat(t *testing.T) {
	conn, err := ssh.NewConnection("192.168.0.127:22", "root", false)
	if err != nil {
		t.Errorf("new connnection error: %v", err)
	}

	bt, err := getBootTime(conn)
	if err != nil {
		t.Errorf("get boot time failed: %s", err.Error())
	}

	cmd := NewCommand(Target{
		Conn:     conn,
		BootTime: bt,
	}, Process{Flag: "ssssssssssssssssss"})

	ps, err := cmd.GetProcessStat()
	if err != nil {
		t.Errorf("get boot time failed: %s", err.Error())
	} else {
		t.Logf("process state: %+v", *ps)
	}
}
