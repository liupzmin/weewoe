package scrape

import (
	"net"
	"strconv"
	"time"

	"github.com/liupzmin/weewoe/log"
)

func RawConnect(host string, port int64) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(int(port))), timeout)
	if err != nil {
		log.Errorf("Connect %s:%s Failed: %s", host, port, err)
		return false
	}

	_ = conn.Close()
	return true
}
