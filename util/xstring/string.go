// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of weewoe

package xstring

import (
	"bufio"
	"strings"
)

func GetNoEmptyLineNumber(input string) (int, error) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanLines)

	var count int
	for scanner.Scan() {
		if scanner.Text() != "" {
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return count, nil
}

func TrimEmptyLines(b []byte) string {
	strs := strings.Split(string(b), "\n")
	str := ""
	for _, s := range strs {
		if len(strings.TrimSpace(s)) == 0 {
			continue
		}
		str += s + "\n"
	}
	str = strings.TrimSuffix(str, "\n")

	return str
}
