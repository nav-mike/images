package usecase

import (
	"crypto/sha1"
	"fmt"
	"regexp"
	"time"
)

type FileValidationError struct {
	Message string
}

func NewFileValidationError(message string) *FileValidationError {
	return &FileValidationError{Message: message}
}

func (e *FileValidationError) Error() string {
	return e.Message
}

func GenerateFilename(baseName, prefix string) (string, error) {
	re, err := regexp.Compile(`(png)|(jpeg)|(jpg)$`)
	if err != nil {
		return "", err
	}

	ext := re.FindString(baseName)
	if ext == "" {
		return "", NewFileValidationError("invalid file extension")
	}

	return fmt.Sprintf(
		"%s-%x.%s",
		prefix,
		sha1.Sum([]byte(baseName+time.Now().String()+"."+ext)),
		ext,
	), nil
}
