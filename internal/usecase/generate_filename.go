package usecase

import (
	"crypto/sha1"
	"fmt"
	"time"
)

func GenerateFilename(baseName, prefix, ext string) string {
	return fmt.Sprintf(
		"%s-%x.%s",
		prefix,
		sha1.Sum([]byte(baseName+time.Now().String()+"."+ext)),
		ext,
	)
}
