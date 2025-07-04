// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of weewoe

package xtime

import "time"

func GetUTC8Location() *time.Location {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.FixedZone("CST-8", 8*3600)
	}
	return loc
}
