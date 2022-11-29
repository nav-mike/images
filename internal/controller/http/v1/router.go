package v1

import (
	"log"
	"net/http"

	"github.com/nav-mike/images/config"
)

func NewRouter(config *config.Config) {
	log.Println("Initialize routes")

	http.HandleFunc("/upload", PostUploadImageHandler(config))

	http.HandleFunc("/images", GetImageHandler(config))

	http.ListenAndServe(":8080", nil)
}
