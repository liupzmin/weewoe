package ssh

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	"github.com/liupzmin/weewoe/log"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type Connection struct {
	client      *ssh.Client
	done        chan struct{}
	isKeepAlive bool
}

func NewConnection(addr, user string, keepAlive bool) (*Connection, error) {
	hostCallBack, err := knownhosts.New(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		return nil, err
	}

	key, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"))
	if err != nil {
		log.Errorf("unable to read private key: %v", err)
		return nil, err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Errorf("unable to parse private key: %v", err)
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostCallBack,
	}
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	conn := &Connection{
		client: client,
	}

	if keepAlive {
		conn.isKeepAlive = true
		conn.done = make(chan struct{})
		go conn.keepAlive()
	}

	return conn, nil
}

func (c *Connection) SingleRun(cmd string) (string, error) {
	sess, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	defer sess.Close()

	var buf, errBuf bytes.Buffer
	sess.Stdout = &buf
	sess.Stderr = &errBuf
	if err := sess.Run(cmd); err != nil {
		return errBuf.String(), err
	}

	return buf.String(), nil
}

func (c *Connection) MultipleRun(commands ...string) (string, error) {
	sess, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	defer sess.Close()

	stdinBuf, err := sess.StdinPipe()
	if err != nil {
		return "", err
	}

	var outBuf, errBuf bytes.Buffer
	sess.Stdout = &outBuf
	sess.Stderr = &errBuf

	err = sess.Shell()
	if err != nil {
		return "", err
	}
	for _, cmd := range commands {
		cmd = cmd + "\n"
		_, _ = stdinBuf.Write([]byte(cmd))
	}
	err = sess.Wait()
	if err != nil {
		return errBuf.String(), err
	}
	return outBuf.String(), nil
}

func (c *Connection) Close() error {
	if c.isKeepAlive {
		c.done <- struct{}{}
	}
	return c.client.Close()
}

func (c *Connection) keepAlive() {
	const keepAliveInterval = time.Minute
	t := time.NewTicker(keepAliveInterval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			_, _, err := c.client.SendRequest("keepalive@golang.org", true, nil)
			if err != nil {
				log.Errorf("failed to send keep alive")
			}
		case <-c.done:
			return
		}
	}
}
