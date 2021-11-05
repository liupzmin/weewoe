package main

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/liupzmin/weewoe/log"
	"github.com/liupzmin/weewoe/ssh"
	"github.com/liupzmin/weewoe/util/xmap"
	"github.com/liupzmin/weewoe/util/xstring"

	ssh2 "golang.org/x/crypto/ssh"
)

type Command struct {
	target Target
	p      Process
}

func NewCommand(t Target, p Process) *Command {
	return &Command{
		target: t,
		p:      p,
	}
}

func (c *Command) GetPID() (string, error) {
	if c.p.Flag == "" && c.p.PIDFile == "" {
		return "", fmt.Errorf("you have to specify either pidfile or flag")
	}
	if c.p.PIDFile != "" {
		return c.getPIDByFile(c.p.PIDFile)
	}

	return c.getPIDByFlag(c.p.Flag)
}

func (c *Command) GetProcessStat() (*ProcessState, error) {
	bad := &ProcessState{
		Process:       c.p,
		State:         0,
		StateDescribe: "",
		StartTime:     0,
		Timestamp:     time.Now().Unix(),
	}

	pid, err := c.GetPID()
	if err != nil {
		return bad, nil
	}
	cmd := fmt.Sprintf("cat /proc/%s/stat", pid)
	output, err := c.target.Conn.SingleRun(cmd)
	if err != nil {
		var exit *ssh2.ExitError
		if errors.As(err, &exit) {
			return bad, nil
		} else {
			return nil, err
		}
	}
	reader := bufio.NewReader(strings.NewReader(output))
	stat, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}

	ps, err := NewProcStat(stat)
	if err != nil {
		return nil, err
	}

	return &ProcessState{
		Process:       c.p,
		State:         Good,
		StateDescribe: ps.State,
		StartTime:     ps.StartTime(c.target.BootTime),
		Timestamp:     time.Now().Unix(),
	}, nil
}

func (c *Command) getPIDByFile(file string) (string, error) {
	cmd := fmt.Sprintf("cat %s", file)
	pid, err := c.target.Conn.SingleRun(cmd)
	if err != nil {
		return "", err
	}

	if pid == "" {
		return "", fmt.Errorf("pid file is empty")
	}
	return pid, nil
}

func (c *Command) getPIDByFlag(flag string) (string, error) {
	flags := xmap.Map(strings.Split(flag, ","), func(s string) string {
		return "grep " + s
	})

	cmd := fmt.Sprintf("ps -ef|%s|grep -v grep|awk '{print $2}'", strings.Join(flags, "|"))
	output, err := c.target.Conn.SingleRun(cmd)
	if err != nil {
		return "", err
	}

	count, err := xstring.GetNoEmptyLineNumber(output)
	if err != nil {
		log.Errorf("get line number failed: %s", err.Error())
		return "", err
	}

	if count > 1 {
		return "", fmt.Errorf("more than 1 line in pid cmd output")
	}

	reader := bufio.NewReader(strings.NewReader(output))
	pid, _, err := reader.ReadLine()
	if err != nil {
		log.Errorf("read pid line failed: %s", err.Error())
		return "", err
	}
	return string(pid), nil
}

func getBootTime(conn *ssh.Connection) (int64, error) {
	cmd := "cat /proc/stat|grep btime|awk '{print $2}'"
	output, err := conn.SingleRun(cmd)
	if err != nil {
		return 0, err
	}
	reader := bufio.NewReader(strings.NewReader(output))
	btime, _, err := reader.ReadLine()
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(btime), 10, 64)
}
