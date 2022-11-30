package usecase

import (
	"image"
	"log"

	"github.com/disintegration/imaging"
)

type ImageSize struct {
	Height int // avoid using width here because proportional scaling is required
	Label  string
}

func (is ImageSize) String() string {
	return is.Label
}

var imageSizes = [...]ImageSize{
	{Label: "micro", Height: 100},
	{Label: "small", Height: 200},
	{Label: "medium", Height: 300},
}

type ImageStorage interface {
	ImageFileFullPath(userId, filename string) string
}

func ResizeImage(storage ImageStorage, originalFilename, userId string) (map[string]string, error) {
	resizedFilenames := make(map[string]string)
	baseImage, err := imaging.Open(storage.ImageFileFullPath(userId, originalFilename))
	if err != nil {
		return nil, err
	}

	for _, size := range imageSizes {
		resizedFilename, err := createResizedImage(storage, baseImage, userId, size)
		if err != nil {
			return nil, err
		}

		resizedFilenames[size.Label] = resizedFilename
	}

	return resizedFilenames, nil
}

func createResizedImage(storage ImageStorage, baseImage image.Image, userId string, size ImageSize) (string, error) {
	log.Printf("Resizing image to %s", size)

	result := imaging.Resize(baseImage, 0, size.Height, imaging.Lanczos)

	// save the resized image
	resizedFilename := GenerateFilename("image.png", size.Label, "png")
	err := imaging.Save(result, storage.ImageFileFullPath(userId, resizedFilename))
	if err != nil {
		return "", err
	}

	return resizedFilename, nil
}
