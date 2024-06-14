package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func Ask() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		s, _ := reader.ReadString('\n')
		s = strings.TrimSuffix(s, "\n")
		s = strings.ToLower(s)
		if len(s) > 1 {
			fmt.Fprintln(os.Stderr, "Please enter Y or N")
			continue
		}
		if strings.Compare(s, "n") == 0 {
			return errors.New("canceled by user")
		} else if strings.Compare(s, "y") == 0 {
			break
		} else {
			continue
		}
	}
	return nil
}
