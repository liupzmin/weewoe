package log_test

import (
	"testing"

	"github.com/liupzmin/weewoe/pkg/log"
)

func Test_Info(t *testing.T) {
	log.Info("hello", log.Any("a", "b"))
}
