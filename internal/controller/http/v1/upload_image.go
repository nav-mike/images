package v1

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/disintegration/imaging"
	"github.com/nav-mike/images/config"
)

type UploadImageDTO struct {
	File   string
	UserId string
}

type ImageSize struct {
	Height int // avoid using width here because proportional scaling is required
	Label  string
}

func (is ImageSize) String() string {
	return is.Label
}

// const array of image sizes
var imageSizes = [...]ImageSize{
	{Label: "micro", Height: 100},
	{Label: "small", Height: 200},
	{Label: "medium", Height: 300},
}

func PostUploadImageHandler(config *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Upload image")
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var file UploadImageDTO

		err := json.NewDecoder(r.Body).Decode(&file)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		err = os.MkdirAll(dirPath(config, file.UserId), os.ModePerm)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		originalFilename, err := saveToFile(config, file)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		resizedFilenames := make([]string, len(imageSizes))

		for i, size := range imageSizes {
			resizedFilenames[i], err = resizeImage(config, originalFilename, file.UserId, size)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		result := make(map[string]string)
		for index, sizeString := range imageSizes {
			result[sizeString.String()] = imageUrl(config, resizedFilenames[index])
		}
		result["original"] = imageUrl(config, originalFilename)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func dirPath(config *config.Config, userId string) string {
	return config.ImagesDir + "/" + userId
}

func saveToFile(config *config.Config, input UploadImageDTO) (string, error) {
	// Decode base64 string to []byte
	decoded, err := base64.StdEncoding.DecodeString(input.File)
	if err != nil {
		return "", err
	}

	// Create file
	filename := generateFilename(config, input.UserId, "original")
	file, err := os.Create(filename)
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

func generateFilename(config *config.Config, userId, prefix string) string {
	if prefix == "" {
		prefix = "image"
	}

	return fmt.Sprintf("%s/%s/%s-%x.png",
		config.ImagesDir,
		userId,
		prefix,
		sha1.Sum([]byte("image.png"+time.Now().String())),
	)
}

func resizeImage(config *config.Config, filename, userId string, size ImageSize) (string, error) {
	// read image from the file
	source, err := imaging.Open(filename)
	if err != nil {
		return "", err
	}

	result := imaging.Resize(source, 0, size.Height, imaging.Lanczos)

	// save the resulting image using png format.
	resizedFilename := generateFilename(config, userId, size.String())
	err = imaging.Save(result, resizedFilename)
	if err != nil {
		return "", err
	}

	return resizedFilename, nil
}

func imageUrl(config *config.Config, filename string) string {
	return config.ServerHostUrl + "/" + filename
}
