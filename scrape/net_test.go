package scrape

import (
	"testing"
)

func TestRawConnect(t *testing.T) {
	open := RawConnect("127.0.0.1", 22)
	if !open {
		t.Errorf("wrong wrong")
	}
}

func TestRawConnect2(t *testing.T) {
	open := RawConnect("127.0.0.1", 2222)
	if open {
		t.Errorf("wrong wrong")
	}
}
