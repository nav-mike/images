package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
	"github.com/joho/godotenv"
)

const IMAGES_DIR = "images"

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

	err = os.MkdirAll(IMAGES_DIR+"/", os.ModePerm)
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
		rszFilename, err := resizeImage(filename, size)
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
	return fmt.Sprintf("%s/%s", os.Getenv("SERVER_HOST_VALUE"), filename)
}

func saveToFile(input UploadImageDTO) (string, error) {
	// Decode base64 string to []byte
	decoded, err := base64.StdEncoding.DecodeString(input.File)
	if err != nil {
		return "", err
	}

	// Create file
	filename := generateFilename("image.png", "original") + ".png"
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

func resizeImage(filename string, size ImageSize) (string, error) {
	// read image from the file
	src, err := imaging.Open(filename)
	if err != nil {
		return "", err
	}

	result := imaging.Resize(src, 0, size.Height, imaging.Lanczos)

	// save the resulting image using png format.
	resizedFilename := generateFilename(filename, size.String()) + ".png"
	err = imaging.Save(result, resizedFilename)
	if err != nil {
		return "", err
	}

	return resizedFilename, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	fs := http.FileServer(http.Dir(IMAGES_DIR))

	http.HandleFunc("/upload", UploadHandler)
	http.Handle("/"+IMAGES_DIR+"/", http.StripPrefix("/"+IMAGES_DIR+"/", fs))
	http.ListenAndServe(":8080", nil)
}
