package utils

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

func ParseFileFlag(raw string) (string, string, error) {
	// expected format <local-path>:<mount-path>
	matchFormat, err := regexp.MatchString(`^[^: ]+:[^: ]+$`, raw)

	if err != nil {
		return "", "", err
	}

	if matchFormat == false {
		msg := "The format of --file flag must be: <local-path>:<semaphore-path>"
		return "", "", fmt.Errorf(msg)
	}

	flagPaths := strings.Split(raw, ":")
	localPath := flagPaths[0]
	remotePath := flagPaths[1]

	// #nosec
	content, err := ioutil.ReadFile(localPath)

	if err != nil {
		return "", "", err
	}

	base64Content := base64.StdEncoding.EncodeToString(content)

	return remotePath, base64Content, nil
}
