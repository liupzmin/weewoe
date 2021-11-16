package main

import (
	"os"
	"os/signal"
	"syscall"
)

func waitSignals() <-chan os.Signal {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	return signals
}
