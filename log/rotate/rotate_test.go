// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of weewoe

//go:build linux
// +build linux

package rotate_test

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/liupzmin/weewoe/log/rotate"
)

// Example of how to rotate in response to SIGHUP.
func ExampleLogger_Rotate() {
	l := &rotate.Logger{}
	log.SetOutput(l)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for {
			<-c
			l.Rotate()
		}
	}()
}
