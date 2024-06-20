package utils

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func Ask(name string) error {
	reader := bufio.NewReader(os.Stdin)
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
