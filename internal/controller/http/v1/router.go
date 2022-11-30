package v1

import (
	"log"
	"net/http"

	"github.com/nav-mike/images/config"
	"github.com/nav-mike/images/internal/entity"
	"github.com/nav-mike/images/internal/usecase/repo/filesystem"
)

// ImageWriter represents service to save image
type ImageWriter interface {
	SaveImage(input entity.UploadImageDTO) (entity.UploadedImageResponse, error)
}

// ImageReader represents service to get image from some storage
type ImageReader interface {
	GetStaticImagePath(userId, requestPath string) (string, error)
}

// NewRouter defines routes for http server
func NewRouter(config *config.Config) {
	log.Println("Initialize routes")

	fs := filesystem.NewFileSystem(config.ImagesDir)

	http.HandleFunc("/upload", PostUploadImageHandler(config, fs))

	http.HandleFunc("/images/", GetImageHandler(fs))

	http.ListenAndServe(":8080", nil)
}
