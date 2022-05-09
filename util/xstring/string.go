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
