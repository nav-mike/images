package filesystem

import (
	"encoding/base64"
	"fmt"
	"os"

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

	result := make(UploadedImageResponse)
	result["original"] = orignalFilename

	return result, nil
}

func (fs *FileSystem) dirPath(userId string) string {
	return fs.Path + "/" + userId
}

func (fs *FileSystem) createDir(userId string) error {
	return os.MkdirAll(fs.dirPath(userId), os.ModePerm)
}

func (fs *FileSystem) filePath(filename, userId string) string {
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
	file, err := os.Create(fs.filePath(filename, input.UserId))
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
