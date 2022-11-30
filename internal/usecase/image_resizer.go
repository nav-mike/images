package usecase

import (
	"image"

	"github.com/disintegration/imaging"
)

// ImageSize represents size of image
type ImageSize struct {
	Height int // avoid using width here because proportional scaling is required
	Label  string
}

func (is ImageSize) String() string {
	return is.Label
}

// imageSizes represents list of sizes for resizing
var imageSizes = [...]ImageSize{
	{Label: "micro", Height: 100},
	{Label: "small", Height: 200},
	{Label: "medium", Height: 300},
}

// ImageStorage defines a service to save image
type ImageStorage interface {
	ImageFileFullPath(userId, filename string) string
}

// ResizeImage resizes image to different (defined) sizes
func ResizeImage(storage ImageStorage, originalFilename, userId string) (map[string]string, error) {
	resizedFilenames := make(map[string]string)
	baseImage, err := imaging.Open(storage.ImageFileFullPath(userId, originalFilename))
	if err != nil {
		return nil, err
	}

	for _, size := range imageSizes {
		resizedFilename, err := createResizedImage(storage, baseImage, originalFilename, userId, size)
		if err != nil {
			return nil, err
		}

		resizedFilenames[size.Label] = resizedFilename
	}

	return resizedFilenames, nil
}

func createResizedImage(storage ImageStorage, baseImage image.Image, baseFilename, userId string, size ImageSize) (string, error) {
	result := imaging.Resize(baseImage, 0, size.Height, imaging.Lanczos)

	// save the resized image
	resizedFilename, err := GenerateFilename(baseFilename, size.Label)
	if err != nil {
		return "", err
	}

	err = imaging.Save(result, storage.ImageFileFullPath(userId, resizedFilename))
	if err != nil {
		return "", err
	}

	return resizedFilename, nil
}
