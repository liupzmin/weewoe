// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of weewoe

package xmap

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
