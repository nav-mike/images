package usecase_test

import (
	"os"
	"testing"

	"github.com/nav-mike/images/internal/usecase"
	"github.com/nav-mike/images/internal/usecase/repo/filesystem"
)

const TEST_IMAGES_PATH = "../../data/images"

func removeTestFiles(files map[string]string) {
	if files == nil {
		return
	}

	for _, filename := range files {
		os.Remove(TEST_IMAGES_PATH + "/test/" + filename)
	}
}

func TestResizeImage(t *testing.T) {
	fs := filesystem.NewFileSystem(TEST_IMAGES_PATH)

	t.Run("should return error if file does not exists", func(t *testing.T) {
		_, err := usecase.ResizeImage(fs, "not-found.png", "test")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("should return map of the image sizes", func(t *testing.T) {
		actual, err := usecase.ResizeImage(fs, "test-img.jpg", "test")
		defer removeTestFiles(actual)

		if err != nil {
			t.Errorf("expected no error, got %s", err.Error())
		} else if len(actual) != 3 {
			t.Errorf("expected 3 sizes, got %d", len(actual))
		}

		for key := range actual {
			if key != "micro" && key != "small" && key != "medium" {
				t.Errorf("expected key to be 'micro', 'small' or 'medium', got %s", key)
			}
		}

		for _, filename := range actual {
			if filename == "" {
				t.Error("expected filename to be not empty")
			}
		}

		for _, filename := range actual {
			if _, err := os.Stat(TEST_IMAGES_PATH + "/test/" + filename); os.IsNotExist(err) {
				t.Errorf("expected file %s to exists", filename)
			}
		}
	})
}
