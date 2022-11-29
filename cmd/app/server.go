package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/nav-mike/images/config"
	v1 "github.com/nav-mike/images/internal/controller/http/v1"
)

const IMAGES_DIR = "data/images"

type UploadImageDTO struct {
	File   string
	UserId string
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

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId := r.Header.Get("X-Custom-Auth-Token") // in real life it'd be JWT token for example
	if userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/"+IMAGES_DIR+"/"+userId) {
		http.ServeFile(w, r, r.URL.Path[1:])
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var file UploadImageDTO

	err := json.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	err = os.MkdirAll(IMAGES_DIR+"/"+file.UserId, os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		rszFilename, err := resizeImage(filename, file.UserId, size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		filenames = append(filenames, rszFilename)
	}

	result := map[string]string{
		"result":               "success",
		"filename":             makeImageUrl(filename),
		imageSizes[0].String(): makeImageUrl(filenames[0]),
		imageSizes[1].String(): makeImageUrl(filenames[1]),
		imageSizes[2].String(): makeImageUrl(filenames[2]),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func makeImageUrl(filename string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("SERVER_HOST_URL"), filename)
}

func saveToFile(input UploadImageDTO) (string, error) {
	// Decode base64 string to []byte
	decoded, err := base64.StdEncoding.DecodeString(input.File)
	if err != nil {
		return "", err
	}

	// Create file
	filename := generateFilename("image.png", input.UserId+"/original") + ".png"
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

func generateFilename(original, prefix string) string {
	if prefix == "" {
		prefix = "original"
	}

	return fmt.Sprintf("%s/%s-%x", IMAGES_DIR, prefix, sha1.Sum([]byte(original)))
}

func resizeImage(filename, userId string, size ImageSize) (string, error) {
	// read image from the file
	src, err := imaging.Open(filename)
	if err != nil {
		return "", err
	}

	result := imaging.Resize(src, 0, size.Height, imaging.Lanczos)

	// save the resulting image using png format.
	resizedFilename := generateFilename(filename, userId+"/"+size.String()) + ".png"
	err = imaging.Save(result, resizedFilename)
	if err != nil {
		return "", err
	}

	return resizedFilename, nil
}

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		return
	}

	log.Println("Starting server on port 8080")

	v1.NewRouter(conf)
}
