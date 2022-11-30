package usecase

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"regexp"
	"time"
)

func GenerateFilename(baseName, prefix string) (string, error) {
	re, err := regexp.Compile(`(png)|(jpeg)|(jpg)$`)
	if err != nil {
		return "", err
	}

	ext := re.FindString(baseName)
	if ext == "" {
		return "", errors.New("invalid file extension")
	}

	return fmt.Sprintf(
		"%s-%x.%s",
		prefix,
		sha1.Sum([]byte(baseName+time.Now().String()+"."+ext)),
		ext,
	), nil
}
