package main

import (
	"net"
	"strconv"
	"time"
)

func RawConnect(host string, port int64) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(int(port))), timeout)
	if err != nil {
		return false
	}

	_ = conn.Close()
	return true
}
