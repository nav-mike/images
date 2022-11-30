package filesystem_test

import (
	"os"
	"testing"

	"github.com/nav-mike/images/internal/entity"
	"github.com/nav-mike/images/internal/usecase/repo/filesystem"
)

const TEST_IMAGES_PATH = "../../../../data/images"

func removeTestFiles(files map[string]string) {
	if files == nil {
		return
	}

	for _, filename := range files {
		os.Remove(TEST_IMAGES_PATH + "/test/" + filename)
	}
}

func TestGetStaticImagePath(t *testing.T) {
	fs := filesystem.NewFileSystem(TEST_IMAGES_PATH)

	t.Run("should return error if image does not exist", func(t *testing.T) {
		_, err := fs.GetStaticImagePath("anotherUser", "test-img.jpg")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("should return path to the image", func(t *testing.T) {
		actual, err := fs.GetStaticImagePath("test", "test-img.jpg")
		if err != nil {
			t.Errorf("expected no error, got %s", err.Error())
		} else if actual != TEST_IMAGES_PATH+"/test/test-img.jpg" {
			t.Errorf("expected path to be '%s', got %s", TEST_IMAGES_PATH+"/test/test-img.jpg", actual)
		}
	})
}

func TestSaveImage(t *testing.T) {
	fs := filesystem.NewFileSystem(TEST_IMAGES_PATH)

	t.Run("should return map of different sizes", func(t *testing.T) {
		file, err := os.ReadFile(TEST_IMAGES_PATH + "/test/correct_base64.txt")
		if err != nil {
			t.Errorf("expected no error, got %s", err.Error())
		}

		input := entity.UploadImageDTO{
			Filename: "test-img.jpg",
			File:     string(file),
			UserId:   "test",
		}
		actual, err := fs.SaveImage(input)
		defer removeTestFiles(actual)

		if err != nil {
			t.Errorf("expected no error, got %s", err.Error())
		} else if len(actual) != 4 {
			t.Errorf("expected 4 sizes, got %d", len(actual))
		}

		for key := range actual {
			if key != "micro" && key != "small" && key != "medium" && key != "original" {
				t.Errorf("expected key to be 'micro', 'small', 'medium' or 'original', got %s", key)
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
