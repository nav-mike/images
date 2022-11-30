package filesystem

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/nav-mike/images/internal/entity"
	"github.com/nav-mike/images/internal/usecase"
)

// MAX_FILE_SIZE represents max file size in MB
const MAX_FILE_SIZE = 1024 * 1024 * 10 // 10MB

// SaveImage saves image to file system and create resized copies of the message. Returns map of images' urls
func (fs *FileSystem) SaveImage(input entity.UploadImageDTO) (entity.UploadedImageResponse, error) {
	err := fs.createDir(input.UserId)
	if err != nil {
		return nil, err
	}

	orignalFilename, err := fs.saveToFile(input)
	if err != nil {
		return nil, err
	}

	result, err := usecase.ResizeImage(fs, orignalFilename, input.UserId)
	if err != nil {
		return nil, err
	}

	result["original"] = orignalFilename

	return result, nil
}

// GetStaticImagePath returns path to image
func (fs *FileSystem) GetStaticImagePath(userId, requestPath string) (string, error) {
	path := fs.ImageFileFullPath(userId, strings.Replace(requestPath, "/images/", "", 1))
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	return path, nil
}

func (fs *FileSystem) dirPath(userId string) string {
	return fs.Path + "/" + userId
}

func (fs *FileSystem) createDir(userId string) error {
	return os.MkdirAll(fs.dirPath(userId), os.ModePerm)
}

func (fs *FileSystem) ImageFileFullPath(userId, filename string) string {
	return fmt.Sprintf("%s/%s/%s", fs.Path, userId, filename)
}

func (fs *FileSystem) saveToFile(input entity.UploadImageDTO) (string, error) {
	filename, err := usecase.GenerateFilename(input.Filename, "original")
	if err != nil {
		return "", err
	}

	// Decode base64 string to []byte
	decoded, err := base64.StdEncoding.DecodeString(input.File)
	if err != nil {
		return "", err
	}

	err = validateFileSize(decoded)
	if err != nil {
		return "", err
	}

	// Create file
	file, err := os.Create(fs.ImageFileFullPath(input.UserId, filename))
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write data to file
	writedSize, err := file.Write(decoded)
	if err != nil {
		return "", err
	}
	if writedSize != len(decoded) {
		return "", errors.New("file size is not equal to decoded size")
	}

	return filename, nil
}

func validateFileSize(base64Value []byte) error {
	if len(base64Value) > MAX_FILE_SIZE {
		return NewValidationError("file size is too big")
	}

	return nil
}
