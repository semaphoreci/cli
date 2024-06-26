package utils

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

func Ask(in io.Reader, name string) error {
	reader := bufio.NewReader(in)
	for {
		s, _ := reader.ReadString('\n')
		s = strings.TrimSuffix(s, "\n")
		s = strings.ToLower(s)
		if strings.Compare(s, name) == 0 {
			break
		} else {
			return errors.New("user confirmation failed")
		}
	}
	return nil
}
