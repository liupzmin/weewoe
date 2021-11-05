package log_test

import (
	"testing"

	log2 "github.com/liupzmin/weewoe/log"
)

func Test_Info(t *testing.T) {
	log2.Info("hello", log2.Any("a", "b"))
}

func Test_Debug(t *testing.T) {
	log2.Debug("hello", log2.Any("a", "b"))
}

func Test_Error(t *testing.T) {
	log2.Errorf("hello", log2.Any("a", "b"))
}
