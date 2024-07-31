package scrape

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/liupzmin/weewoe/log"
	"github.com/liupzmin/weewoe/ssh"
	"github.com/liupzmin/weewoe/util/xmap"
	"github.com/liupzmin/weewoe/util/xstring"

	ssh2 "golang.org/x/crypto/ssh"
)

const MaxRetries = 5

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
		return "", fmt.Errorf("you have to specify either pidfile or flag for %s", c.p.Name)
	}
	if c.p.PIDFile != "" {
		return c.getPIDByFile(c.p.PIDFile)
	}

	return c.getPIDByFlag(c.p.Flag)
}

func (c *Command) GetProcessStat() (*ProcessState, error) {
	bad := &ProcessState{
		Process:       c.p,
		State:         Bad,
		StateDescribe: "",
		StartTime:     0,
		Timestamp:     time.Now().Unix(),
	}

	btime, err := GetBootTime(c.target.Conn)
	if err != nil {
		log.Warnf("GetBootTime failed for %s, error: %s", c.target.Conn.Addr(), err)
		return bad, nil
	}

	pid, err := c.GetPID()
	if err != nil {
		log.Warnf("GetPID failed for %s, error: %s", c.p.Name, err)
		return bad, nil
	}
	log.Debugf("%s's pid is %s", c.p.Name, pid)
	cmd := fmt.Sprintf("/bin/cat /proc/%s/stat", strings.TrimSpace(pid))
	output, err := c.target.Conn.SingleRun(cmd)
	if err != nil {
		var exit *ssh2.ExitError
		if errors.As(err, &exit) {
			log.Warnf("GetProcessStat error exit: %s", err)
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
		StartTime:     ps.StartTime(btime),
		Timestamp:     time.Now().Unix(),
	}, nil
}

func (c *Command) getPIDByFile(file string) (out string, err error) {
	cmd := fmt.Sprintf("cat %s", file)

	for i := 0; i < MaxRetries; i++ {
		out, err = c.target.Conn.SingleRun(cmd)
		if err == nil {
			return
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return "", fmt.Errorf("%s : %s", err, out)
	}

	if out == "" {
		return "", fmt.Errorf("pid file is empty")
	}
	return out, nil
}

func (c *Command) getPIDByFlag(flag string) (output string, err error) {
	flags := xmap.Map(strings.Split(flag, ","), func(s string) string {
		return "grep " + s
	})

	//cmd := fmt.Sprintf("ps -ef|%s|grep -v grep|awk '{print $2}'", strings.Join(flags, "|"))
	cmd := fmt.Sprintf("ps -ef|%s|grep -v grep", strings.Join(flags, "|"))
	log.Infof("Run CMD: %s", cmd)
	for i := 0; i < MaxRetries; i++ {
		output, err = c.target.Conn.SingleRun(cmd)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		log.Errorf("Get Failed for CMD [%s], error: %s", cmd, err)
		return "", err
	}
	output = xstring.TrimEmptyLines([]byte(output))
	log.Infof("Output for [%s]: %s", flag, output)

	count, err := xstring.GetNoEmptyLineNumber(output)
	if err != nil {
		log.Errorf("get line number failed: %s", err.Error())
		return "", err
	}

	if count > 1 {
		return "", fmt.Errorf("more than 1 line in pid cmd output")
	}

	reader := bufio.NewReader(strings.NewReader(output))
	pidStr, _, err := reader.ReadLine()
	if err == io.EOF {
		return "", err
	}
	if err != nil {
		log.Errorf("read pid line failed: %s", err.Error())
		return "", err
	}

	pidSlice := strings.Fields(string(pidStr))

	return pidSlice[1], nil
}

func GetBootTime(conn *ssh.Connection) (int64, error) {
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
