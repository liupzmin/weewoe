// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of weewoe

package rotate_test

import (
	"log"

	"github.com/liupzmin/weewoe/log/rotate"
)

// To use rotate with the standard library's log package, just pass it into
// the SetOutput function when your application starts.
func Example() {
	log.SetOutput(&rotate.Logger{
		Filename:   "/var/log/myapp/foo.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // disabled by default
	})
}
