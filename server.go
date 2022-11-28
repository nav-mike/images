package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
)

type UploadImageDTO struct {
	File string
}

type ImageSize struct {
	Width  int
	Height int
}

func (is ImageSize) String() string {
	return fmt.Sprintf("%dx%d", is.Width, is.Height)
}

// const array of image sizes
var imageSizes = [...]ImageSize{
	{Width: 100, Height: 100},
	{Width: 200, Height: 200},
	{Width: 300, Height: 300},
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	var file UploadImageDTO

	err := json.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	filename, err := saveToFile(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	filenames := []string{}

	for _, size := range imageSizes {
		rszFilename, err := resizeImage(filename, size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		filenames = append(filenames, rszFilename)
	}

	result := map[string]string{"result": "success", "filename": filename, imageSizes[0].String(): filenames[0], imageSizes[1].String(): filenames[1], imageSizes[2].String(): filenames[2]}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func saveToFile(input UploadImageDTO) (string, error) {
	// Decode base64 string to []byte
	decoded, err := base64.StdEncoding.DecodeString(input.File)
	if err != nil {
		return "", err
	}

	// Create file
	filename := generateFilename("image.png") + ".png"
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

func generateFilename(original string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(original)))
}

func resizeImage(filename string, size ImageSize) (string, error) {
	// read image from the file
	src, err := imaging.Open(filename)
	if err != nil {
		return "", err
	}

	result := imaging.Resize(src, 0, size.Height, imaging.Lanczos)

	// save the resulting image using png format.
	resizedFilename := size.String() + "-" + generateFilename(filename) + ".png"
	err = imaging.Save(result, resizedFilename)
	if err != nil {
		return "", err
	}

	return resizedFilename, nil
}

func main() {
	http.HandleFunc("/upload", UploadHandler)
	http.ListenAndServe(":8080", nil)
}
