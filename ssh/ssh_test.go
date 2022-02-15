package ssh

import (
	"errors"
	"testing"

	ssh2 "golang.org/x/crypto/ssh"
)

func TestNewConnection(t *testing.T) {
	conn, err := NewConnection("192.168.0.127:22", "root")
	if err != nil {
		t.Errorf("new connnection error: %v", err)
	}
	conn.Close()
}

func TestConnection_SingleRun(t *testing.T) {
	conn, err := NewConnection("192.168.0.127:22", "root")
	if err != nil {
		t.Errorf("new connnection error: %v", err)
	}
	out, err := conn.SingleRun("cat /proc/9527/stat")
	if err != nil {
		var exit *ssh2.ExitError
		if errors.As(err, &exit) {
			t.Logf("exit err: %s", err.Error())
		} else {
			t.Errorf("single run failed: %s", err)
		}
	}
	t.Logf("single run output: %s", out)
	conn.Close()
}

func TestConnection_MultipleRun(t *testing.T) {
	conn, err := NewConnection("192.168.0.127:22", "root")
	if err != nil {
		t.Errorf("new connnection error: %v", err)
	}

	out, err := conn.MultipleRun("ls", "pwd", "date", "whoami", "exit")
	if err != nil {
		t.Errorf("single run failed: %s", err)
	}
	t.Logf("multiple run output: %s", out)
	conn.Close()
}
