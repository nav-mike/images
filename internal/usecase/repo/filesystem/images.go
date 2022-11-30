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

func (fs *FileSystem) SaveImage(input entity.UploadImageDTO) (UploadedImageResponse, error) {
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
	// Decode base64 string to []byte
	decoded, err := base64.StdEncoding.DecodeString(input.File)
	if err != nil {
		return "", err
	}

	// Create file
	filename := usecase.GenerateFilename("image.png", "original", "png")
	file, err := os.Create(fs.ImageFileFullPath(input.UserId, filename))
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write data to file
	_, err = file.Write(decoded)
	if err != nil {
		return "", err
	}

	return filename, nil
}
