package ssh

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/liupzmin/weewoe/log"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type Connection struct {
	sync.RWMutex
	user, addr          string
	client              *ssh.Client
	cc                  *ssh.ClientConfig
	done, stopKeepAlive chan struct{}
	reconnect           chan error
	valid               bool
}

func NewConnection(addr, user string) (*Connection, error) {
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
		Timeout:         30 * time.Second,
	}

	conn := &Connection{
		user:          user,
		addr:          addr,
		cc:            config,
		done:          make(chan struct{}),
		stopKeepAlive: make(chan struct{}),
		reconnect:     make(chan error),
	}

	go conn.connect()

	return conn, nil
}

func (c *Connection) IsValid() bool {
	return c.valid
}

func (c *Connection) SingleRun(cmd string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("connection not complete")
	}
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
	if c.client == nil {
		return "", fmt.Errorf("connection not complete")
	}
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

func (c *Connection) Close() {
	c.RLock()
	defer c.RUnlock()
	c.done <- struct{}{}
}

func (c *Connection) connect() {
	err := c.dial()
	if err != nil {
		log.Warnf("connect to ssh server %s failed: %s, enter into trying loop", c.addr, err)
		for {
			<-time.After(10 * time.Second)
			err = c.dial()
			if err != nil {
				log.Warnf("connect to ssh server %s failed: %s, try again after 10s", c.addr, err)
				continue
			}
			break
		}
	}
	log.Infof("connect to %s successful!", c.addr)
	go func() {

		c.Lock()
		c.valid = true
		c.Unlock()

		select {
		case err := <-c.reconnect:
			if err != nil {

				c.Lock()
				c.valid = false
				c.Unlock()

				log.Warnf("reconnecting to ssh server %s", c.addr)

				c.stopKeepAlive <- struct{}{}
				_ = c.client.Close()

				go c.connect()
				return
			}
		case <-c.done:
			if c.client != nil {
				c.stopKeepAlive <- struct{}{}
				_ = c.client.Close()
			}
			return
		}
	}()
}

func (c *Connection) dial() error {
	//client, err := ssh.Dial("tcp", c.addr, c.cc)
	//if err != nil {
	//	return err
	//}

	timeout := 60 * time.Second

	conn, err := net.DialTimeout("tcp", c.addr, timeout)
	if err != nil {
		return err
	}

	timeoutConn := &Conn{conn, timeout, timeout}
	co, chans, reqs, err := ssh.NewClientConn(timeoutConn, c.addr, c.cc)
	if err != nil {
		return err
	}
	client := ssh.NewClient(co, chans, reqs)

	c.client = client
	go c.keepAlive()
	return nil
}

func (c *Connection) keepAlive() {
	const keepAliveInterval = 30 * time.Second
	t := time.NewTicker(keepAliveInterval)
	defer t.Stop()
	retry := 3
	for {
		select {
		case <-t.C:
			_, _, err := c.client.SendRequest("keepalive@golang.org", true, nil)
			if err != nil {
				if retry != 0 {
					log.Errorf("failed to send keep alive: %s", err)
					retry--
					continue
				}
				c.reconnect <- err
			} else {
				log.Debugf("ssh keep alive: %s", c.addr)
			}

		case <-c.stopKeepAlive:
			return
		}
	}
}

// Conn wraps a net.Conn, and sets a deadline for every read
// and write operation.
type Conn struct {
	net.Conn
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (c *Conn) Read(b []byte) (int, error) {
	err := c.Conn.SetReadDeadline(time.Now().Add(c.ReadTimeout))
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = c.Conn.SetReadDeadline(time.Time{})
	}()
	return c.Conn.Read(b)
}

func (c *Conn) Write(b []byte) (int, error) {
	err := c.Conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = c.Conn.SetWriteDeadline(time.Time{})
	}()
	return c.Conn.Write(b)
}
