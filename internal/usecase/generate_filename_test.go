package usecase_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/nav-mike/images/internal/usecase"
)

func TestGenerateFilename(t *testing.T) {
	t.Run("should return error if file extension is invalid", func(t *testing.T) {
		actual, err := usecase.GenerateFilename("test.txt", "prefix")
		if err == nil {
			t.Errorf("expected error, got %s", actual)
		} else if err.Error() != "invalid file extension" {
			t.Errorf("expected error to be 'invalid file extension', got %s", err.Error())
		}
	})

	t.Run("should return filename with prefix and sha1 hash", func(t *testing.T) {
		re := regexp.MustCompile(`^prefix-[a-f0-9]{40}\.png$`)
		actual, err := usecase.GenerateFilename("test.png", "prefix")
		if err != nil {
			t.Errorf("expected no error, got %s", err.Error())
		} else if strings.HasPrefix(actual, "prefix-") == false {
			t.Errorf("expected filename to has 'prefix'', got %s", actual)
		} else if strings.HasSuffix(actual, ".png") == false {
			t.Errorf("expected filename to has '.png', got %s", actual)
		} else if re.MatchString(actual) == false {
			t.Errorf("expected filename to match regex, got %s", actual)
		}
	})
}
