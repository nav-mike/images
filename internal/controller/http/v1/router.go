package v1

import (
	"log"
	"net/http"

	"github.com/nav-mike/images/config"
	"github.com/nav-mike/images/internal/entity"
	"github.com/nav-mike/images/internal/usecase/repo/filesystem"
)

type ImageWriter interface {
	SaveImage(input entity.UploadImageDTO) (entity.UploadedImageResponse, error)
}

type ImageReader interface {
	GetStaticImagePath(userId, requestPath string) (string, error)
}

func NewRouter(config *config.Config) {
	log.Println("Initialize routes")

	fs := filesystem.NewFileSystem(config.ImagesDir)

	http.HandleFunc("/upload", PostUploadImageHandler(config, fs))

	http.HandleFunc("/images/", GetImageHandler(fs))

	http.ListenAndServe(":8080", nil)
}
