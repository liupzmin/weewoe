package ssh

import (
	"testing"
)

func TestNewConnection(t *testing.T) {
	conn, err := NewConnection("192.168.0.127:22", "root", true)
	if err != nil {
		t.Errorf("new connnection error: %v", err)
	}
	conn.Close()
}

func TestConnection_SingleRun(t *testing.T) {
	conn, err := NewConnection("192.168.0.127:22", "root", false)
	if err != nil {
		t.Errorf("new connnection error: %v", err)
	}
	out, err := conn.SingleRun("date")
	if err != nil {
		t.Errorf("single run failed: %s", err)
	}
	t.Logf("single run output: %s", out)
	_ = conn.Close()
}

func TestConnection_MultipleRun(t *testing.T) {
	conn, err := NewConnection("192.168.0.127:22", "root", false)
	if err != nil {
		t.Errorf("new connnection error: %v", err)
	}

	out, err := conn.MultipleRun("ls", "pwd", "date", "whoami", "exit")
	if err != nil {
		t.Errorf("single run failed: %s", err)
	}
	t.Logf("multiple run output: %s", out)
	_ = conn.Close()
}
